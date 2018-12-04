package server

import (
	"login/protocol"
	"login/util"
)

type access struct {
	ctx      *Context
	Pools    int
	MsgQueue []chan *protocol.VarMessage
	ExitChan chan struct{}
	wg       util.WaitGroupWrapper
}

func NewAccess(ctx *Context, pools int, queuelen int) *access {
	t := &access{ctx: ctx}
	t.Pools = pools
	t.MsgQueue = make([]chan *protocol.VarMessage, pools)
	for i := 0; i < pools; i++ {
		t.MsgQueue[i] = make(chan *protocol.VarMessage, queuelen)
	}
	t.ExitChan = make(chan struct{})
	return t
}

func (t *access) Start() error {
	for i := 0; i < t.Pools; i++ {
		id := i
		t.wg.Wrap(func() { t.work(id) })
	}

	t.ctx.server.log.Println("start works:", t.Pools)
	return nil
}

func (t *access) Close() {
	close(t.ExitChan)
}

func (t *access) Wait() {
	t.wg.Wait()
}

func (t *access) PushMessage(msg *protocol.VarMessage) {
	msgtype := msg.StringVal(1)
	var workid uint32
	switch msgtype {
	case "login":
		acc := msg.StringVal(3)
		workid = util.Hash(acc)
	case "logout":
		loginid := msg.StringVal(3)
		workid = util.Hash(loginid)
	default:
		workid = uint32(msg.ConnId & 0xFFFF)
	}

	t.MsgQueue[int(workid)%t.Pools] <- msg
}

func (t *access) work(id int) {
	queue := t.MsgQueue[id]
	for {
		select {
		case m := <-queue:
			t.ProcessMsg(m)
		case <-t.ExitChan:
			return
		}
	}
}

func (t *access) ProcessMsg(msg *protocol.VarMessage) {
	msg_type := msg.StringVal(1)
	switch msg_type {
	case "login":
		UserLogin(t.ctx, msg)
	case "logout":
		UserLogout(t.ctx, msg)
	default:
		t.ctx.server.log.Println("unknown msg type:", msg_type)
	}

	t.ctx.server.log.Println("recv msg:", msg_type)
}
