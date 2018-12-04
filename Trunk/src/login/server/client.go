package server

import (
	"bufio"
	"login/protocol"
	"login/setting"
	"net"
	"strconv"
	"time"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	Id  int64
	ctx *Context
	net.Conn
	quit     bool
	exitchan chan struct{}
	// reading/writing interfaces
	Reader            *bufio.Reader
	Writer            *bufio.Writer
	HeartbeatInterval time.Duration
	sendqueue         chan *protocol.VarMessage
	Addr              string
	Port              int
	lenBuf            [4]byte
	lenSlice          []byte
}

func newClient(id int64, conn net.Conn, ctx *Context) *Client {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &Client{
		Id:                id,
		ctx:               ctx,
		Conn:              conn,
		Reader:            bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:            bufio.NewWriterSize(conn, defaultBufferSize),
		HeartbeatInterval: time.Duration(setting.HeartTimeout) * time.Second,
		sendqueue:         make(chan *protocol.VarMessage, 32),
		exitchan:          make(chan struct{}),
		Addr:              addr,
		Port:              int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *Client) SendMessage(msg *protocol.VarMessage) bool {
	if c.quit {
		return false
	}

	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

func (c *Client) Quit() {
	if !c.quit {
		c.quit = true
		close(c.exitchan)
		c.Close()
	}
}
