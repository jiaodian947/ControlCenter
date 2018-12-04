package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	VS_NONE        = iota
	VS_VERIFYING   //验证中
	VS_VERIFYED    //验证服务器已经返回
	VS_FAILED      //验证失败，无效的凭证
	VS_WAITCONFIRM //等待游戏服务器确认
	VS_DONE        //完成验证
)

// 苹果交易数据库
type AppleTransaction struct {
	Id          uint64    //流水号
	OrderId     string    `orm:"size(128);unique"` //订单号
	ProductId   string    `orm:"size(128)"`        //物品编号
	Price       int       //价格
	Timestamp   time.Time `orm:"type(datetime);auto_now_add"` //订单时间
	ServerId    int       //服务器ID
	GameId      int       //游戏ID
	Account     string    //帐号
	RoleGuid    uint64    //角色唯一ID
	Receipt     string    `orm:"type(text)"`     //凭证
	ReceiptMd5  string    `orm:"size(32);index"` //凭证md5,用来快速查找
	VerifyState int8      //验证状态
}

func (t *AppleTransaction) Insert() error {
	if _, err := orm.NewOrm().Insert(t); err != nil {
		return err
	}
	return nil
}

func (t *AppleTransaction) Read(fields ...string) error {
	if err := orm.NewOrm().Read(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *AppleTransaction) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *AppleTransaction) Delete() error {
	if _, err := orm.NewOrm().Delete(t); err != nil {
		return err
	}
	return nil
}

func FindByReceipt(receipt, md5 string) *AppleTransaction {
	var ret []*AppleTransaction
	_, err := orm.NewOrm().QueryTable("cc_apple_transaction").Filter("ReceiptMd5", md5).All(&ret)
	if err != nil {
		return nil
	}

	for _, v := range ret {
		if v.Receipt == receipt {
			return v
		}
	}

	return nil
}
