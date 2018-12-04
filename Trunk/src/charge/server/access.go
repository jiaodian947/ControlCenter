package server

import (
	"charge/protocol"
	"charge/util"
)

// 消息接入通路，通过多路消息队列进行并行处理
type access struct {
	ctx      *Context
	Pools    int
	MsgQueue []chan *protocol.VarMessage
	ExitChan chan struct{}
	handler  *CustomHandler
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
	t.handler = NewCustomHandler(ctx)
	return t
}

// 启动多个工作线程
func (t *access) Start() error {
	for i := 0; i < t.Pools; i++ {
		id := i
		t.wg.Wrap(func() { t.work(id) })
	}

	t.ctx.Server.log.Println("start works:", t.Pools)
	return nil
}

// 关闭所有的工作线程
func (t *access) Close() {
	close(t.ExitChan)
}

func (t *access) Wait() {
	t.wg.Wait()
}

// 压入消息，并进行消息分发处理
func (t *access) PushMessage(msg *protocol.VarMessage) {
	msgtype := msg.StringVal(1)
	var workid uint32
	switch msgtype {
	case "custom":
		identity := msg.StringVal(0)
		workid = util.Hash(identity)
	default:
		workid = uint32(msg.ConnId & 0xFFFF)
	}

	t.MsgQueue[int(workid)%t.Pools] <- msg
}

// 实际的工作线程，负责从消息队列中读取消息，并进行处理
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

// 消息处理函数
func (t *access) ProcessMsg(msg *protocol.VarMessage) {
	msg_type := msg.StringVal(1)
	switch msg_type {
	case "custom":
		t.handler.Handler(msg)
	default:
		t.ctx.Server.log.Println("unknown msg type:", msg_type)
	}

	t.ctx.Server.log.Println("recv msg:", msg_type)
}
