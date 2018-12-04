package server

import "charge/protocol"

// 自定义消息处理类
type CustomHandler struct {
	ctx *Context
}

func NewCustomHandler(ctx *Context) *CustomHandler {
	c := &CustomHandler{ctx}
	return c
}

func (c *CustomHandler) Handler(msg *protocol.VarMessage) {
	identity := msg.StringVal(0)
	index := 2
	var transaction Transaction
	transaction = nil
	optype := msg.StringVal(index)
	c.ctx.Server.log.Println("recv custom", optype)
	switch optype {
	case "verify", "confirm":
		index++
		platform := msg.StringVal(index)
		transaction = RetrieveTransaction(platform, optype, c.ctx.Server.log)
		if transaction == nil {
			c.ctx.Server.log.Println("transaction not found", platform, optype)
			return
		}
		transaction.SetIdentity(identity)
		transaction.SetConnId(msg.ConnId)
		if err := transaction.ParseArgs(msg); err != nil {
			c.ctx.Server.log.Println("parse error:", err)
		}
		c.ctx.Server.AddTrader(transaction)
	default:
		c.ctx.Server.log.Println("unknown custom msg", optype)
	}
}
