package server

import (
	"charge/protocol"
	"fmt"
	"log"
	"time"
)

type newfn func(string, *log.Logger) Transaction
type Transaction interface {
	Platform() string
	SetUrl(url string)
	Url() string
	SetIdentity(identity string)
	Identity() string
	SetConnId(connid int64)
	ConnId() int64
	SetServerId(id int)
	ServerId() int
	SetGameId(id int)
	GameId() int
	ParseArgs(msg *protocol.VarMessage) error
	Check() bool
	Process(ch chan Trader) error
	VerifyTime() time.Time
	Err() error
	ErrCode() int
	Complete(ctx *Context)
}

var (
	transactions = make(map[string]newfn)
)

// 注册交易，platform为平台名称，一个平台对应一个
func RegisterTransaction(platform string, fn newfn) {
	if _, dup := transactions[platform]; dup {
		panic(fmt.Sprintf("transaction(%s) is dup", platform))
	}

	transactions[platform] = fn
}

// 获取交易实例，通过平台和交易类型获取
func RetrieveTransaction(platform string, typ string, l *log.Logger) Transaction {
	if fn, find := transactions[platform]; find {
		return fn(typ, l)
	}

	return nil
}
