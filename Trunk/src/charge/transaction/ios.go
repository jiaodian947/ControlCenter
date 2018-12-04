package transaction

import (
	"bytes"
	"charge/models"
	"charge/protocol"
	"charge/server"
	"charge/setting"
	"charge/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
)

// 苹果凭证验证
type AppleVerifyReceipt struct {
	*log.Logger
	BaseTransaction
	sandboxurl  string
	transaction *models.AppleTransaction
}

func NewAppleVerifyReceipt(l *log.Logger) server.Transaction {
	v := &AppleVerifyReceipt{}
	v.Logger = l
	v.platform = "ios"
	v.url = setting.Platforms["ios"].Path
	v.sandboxurl = setting.Platforms["ios"].TestPath
	return v
}

// 消息协议
// "verify" string	// 固定字符串(index:2)
// Platform string	// 平台标识
// GameId int32		// 游戏ID
// ServerId int32	// 服务器ID
// Account string 	// 玩家帐号
// RoleGuid uint64 	// 玩家唯一ID
// OrderId string	// 交易流水号
// Receipt string 	// 凭证
// ProductId string	// 物品ID
// Price int32		// 价格
func (i *AppleVerifyReceipt) ParseArgs(args *protocol.VarMessage) error {
	index := 4
	i.gameId = args.IntVal(index)
	index++
	i.serverId = args.IntVal(index)
	index++
	account := args.StringVal(index)
	index++
	roleguid := args.Int64Val(index)
	index++
	orderid := args.StringVal(index)
	index++
	receipt := args.StringVal(index)
	index++
	productid := args.StringVal(index)
	index++
	price := args.IntVal(index)

	if account == "" || roleguid == 0 || orderid == "" || receipt == "" {
		return fmt.Errorf("args error")
	}
	md5code := util.MD5([]byte(receipt))
	i.transaction = &models.AppleTransaction{}
	i.transaction.GameId = i.gameId
	i.transaction.ServerId = i.serverId
	i.transaction.Account = account
	i.transaction.RoleGuid = uint64(roleguid)
	i.transaction.OrderId = orderid
	i.transaction.Receipt = receipt
	i.transaction.ReceiptMd5 = md5code
	i.transaction.ProductId = productid
	i.transaction.Price = price
	return nil
}

func (i *AppleVerifyReceipt) Check() bool {
	transaction := models.FindByReceipt(i.transaction.Receipt, i.transaction.ReceiptMd5)
	if transaction != nil { // 已经存在
		if transaction.VerifyState == models.VS_DONE { //已经完成
			i.errcode = protocol.CAS_ERR_ALREADY_DONE
			return false
		}

		if transaction.VerifyState == models.VS_FAILED { //已经验证为失败的凭证
			i.errcode = protocol.CAS_ERR_VERIFY_FAILED
			return false
		}

		if i.transaction.GameId != transaction.GameId ||
			i.transaction.ServerId != transaction.ServerId ||
			i.transaction.RoleGuid != transaction.RoleGuid ||
			i.transaction.ProductId != transaction.ProductId {
			i.errcode = protocol.CAS_ERR_BILL_NOT_MATCH
			return false
		}

		if transaction.VerifyState < models.VS_DONE && time.Now().Sub(transaction.Timestamp).Seconds() < setting.VerifyInterval { //短时间内重复提交
			i.errcode = protocol.CAS_ERR_VERIFY_TOO_OFTEN
			return false
		}

		i.transaction = transaction
		return true
	}
	// 不存在
	if err := i.transaction.Insert(); err != nil {
		i.errcode = protocol.CAS_ERR_BILL_ERROR
		return false
	}
	if err := i.transaction.Read("OrderId"); err != nil {
		i.errcode = protocol.CAS_ERR_BILL_ERROR
		return false
	}
	i.Println("verify order", i.transaction.OrderId)
	return true
}

func (i *AppleVerifyReceipt) Process(done chan server.Trader) error {
	i.verifytime = time.Now()
	i.transaction.Timestamp = time.Now()
	i.transaction.VerifyState = models.VS_VERIFYING
	if err := i.transaction.Update("Timestamp", "VerifyState"); err != nil {
		return err
	}
	go i.doVerify(done)
	return nil
}

type ReceiptData struct {
	Receipt string `json:"receipt-data"`
}

type Receipt struct {
	TransactionId string `json:"transaction_id"`
}

type ReceiptResp struct {
	Status  int     `json:"status"`
	Receipt Receipt `json:"receipt"`
}

func (i *AppleVerifyReceipt) doVerify(done chan server.Trader) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	rd := &ReceiptData{}
	/*
		data, err := Decompress([]byte(i.transaction.Receipt))
		if err != nil {
			i.err = err
			done <- i
		}
		fmt.Println("dec", string(data))*/
	rd.Receipt = i.transaction.Receipt
	postdata, _ := json.Marshal(rd)
	sandbox := false
	url := ""
	retry := 0
	for {
		if !sandbox {
			url = i.url
		} else {
			url = i.sandboxurl
		}
		i.Println("post data to", url)
		resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(postdata))
		if err != nil {
			retry++
			if retry <= setting.RetryMax { //重试几次
				continue
			}
			i.Println("post data error", err)
			i.err = err
			done <- i
			return
		}

		buf, err := ioutil.ReadAll(resp.Body)
		receipt := &ReceiptResp{}
		if err := json.Unmarshal(buf, receipt); err != nil {
			resp.Body.Close()
			i.err = err
			done <- i
			return
		}
		resp.Body.Close()
		i.Println("recv resp status", receipt.Status)
		if receipt.Status == 21007 { //receipt是Sandbox receipt，但却发送至生产系统的验证服务
			sandbox = true
			retry = 0
			continue
		}

		if receipt.Status != 0 { //失败
			i.transaction.VerifyState = models.VS_FAILED
			i.errcode = protocol.CAS_ERR_VERIFY_FAILED
			break
		}

		i.transaction.VerifyState = models.VS_VERIFYED
		break
	}

	done <- i
}

func (i *AppleVerifyReceipt) Complete(ctx *server.Context) {
	client := ctx.Server.FindClient(i.connId)
	if client == nil {
		return
	}
	msg := protocol.NewVarMsg(10)
	msg.AddString(i.identity)
	msg.AddString("custom")
	if i.errcode != 0 { //如果有错误码，直接发送错误码
		msg.AddString("error")
		msg.AddInt(i.errcode)
		msg.AddString(i.transaction.OrderId)
		msg.AddInt64(int64(i.transaction.RoleGuid))
		client.SendMessage(msg)
		return
	}

	if i.err != nil { //验证过程出错
		i.transaction.VerifyState = models.VS_NONE //置为初始状态
		for retry := 0; retry < 3; retry++ {
			if err := i.transaction.Update("VerifyState"); err != nil {
				continue
			}
			break
		}
		msg.AddString("error")
		msg.AddInt(protocol.CAS_ERR_BILL_ERROR)
		msg.AddString(i.transaction.OrderId)
		msg.AddInt64(int64(i.transaction.RoleGuid))
		client.SendMessage(msg)
		return
	}

	i.transaction.VerifyState = models.VS_WAITCONFIRM //改变状态，等待游戏服务器确认订单
	for retry := 0; retry < 3; retry++ {
		if err := i.transaction.Update("VerifyState"); err != nil {
			continue
		}
		//发送帐单验证通过
		// "verify"
		// Platform string	// 平台标识
		// OrderId  string	// 交易流水号
		// RoleGuid int64 	// 玩家唯一ID
		msg.AddString("verify")
		msg.AddString(i.platform)
		msg.AddString(i.transaction.OrderId)
		msg.AddInt64(int64(i.transaction.RoleGuid))
		client.SendMessage(msg)
		return
	}
	// 写库失败
	msg.AddString("error")
	msg.AddInt(protocol.CAS_ERR_BILL_ERROR)
	msg.AddString(i.transaction.OrderId)
	msg.AddInt64(int64(i.transaction.RoleGuid))
	client.SendMessage(msg)
}

// 苹果交易服务器确认
type AppleTransactionConfirm struct {
	*log.Logger
	BaseTransaction
	transaction *models.AppleTransaction
}

func NewAppleTransactionConfirm(l *log.Logger) server.Transaction {
	c := &AppleTransactionConfirm{}
	c.platform = "ios"
	c.Logger = l
	return c
}

// 消息协议
// "confirm" string	// 固定字符串(index:2)
// Platform  string	// 平台标识
// OrderId   string	// 交易流水号
func (c *AppleTransactionConfirm) ParseArgs(args *protocol.VarMessage) error {
	orderid := args.StringVal(4)
	if orderid == "" {
		return fmt.Errorf("order id is empty")
	}
	c.transaction = &models.AppleTransaction{}
	c.transaction.OrderId = orderid
	return nil
}

func (c *AppleTransactionConfirm) Check() bool {
	err := c.transaction.Read("OrderId")

	if err == orm.ErrNoRows {
		c.errcode = protocol.CAS_ERR_TRANSACTION_NOT_FOUND
		return false
	}

	if err != nil { //读取交易失败
		c.errcode = protocol.CAS_ERR_BILL_ERROR
		return false
	}

	if c.transaction.VerifyState != models.VS_WAITCONFIRM {
		c.errcode = protocol.CAS_ERR_BILL_NOT_MATCH
		return false
	}

	return true
}

func (c *AppleTransactionConfirm) Process(done chan server.Trader) error {
	c.Println("confirm", c.transaction.OrderId)
	done <- c //直接完成
	return nil
}

func (c *AppleTransactionConfirm) Complete(ctx *server.Context) {
	client := ctx.Server.FindClient(c.connId)
	if client == nil {
		return
	}
	msg := protocol.NewVarMsg(10)
	msg.AddString(c.identity)
	msg.AddString("custom")
	if c.errcode != 0 { //如果有错误码，直接发送错误码
		msg.AddString("error")
		msg.AddInt(c.errcode)
		msg.AddString(c.transaction.OrderId)
		msg.AddInt64(int64(c.transaction.RoleGuid))
		client.SendMessage(msg)
		return
	}

	c.transaction.VerifyState = models.VS_DONE
	if err := c.transaction.Update("VerifyState"); err != nil {
		msg.AddString("error")
		msg.AddInt(protocol.CAS_ERR_BILL_ERROR)
		msg.AddString(c.transaction.OrderId)
		msg.AddInt64(int64(c.transaction.RoleGuid))
		client.SendMessage(msg)
		return
	}

	msg.AddString("confirm")
	msg.AddString("ios")
	msg.AddString(c.transaction.OrderId)
	msg.AddInt64(int64(c.transaction.RoleGuid))
	client.SendMessage(msg)
	c.Println("confirm", c.transaction.OrderId, "complete")
}

func NewAppleTransaction(typ string, l *log.Logger) server.Transaction {
	switch typ {
	case "verify":
		return NewAppleVerifyReceipt(l)
	case "confirm":
		return NewAppleTransactionConfirm(l)
	default:
		return nil
	}
}
