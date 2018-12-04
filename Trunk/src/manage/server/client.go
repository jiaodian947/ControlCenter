package server

import (
	"bufio"
	"encoding/binary"
	"manage/protocol"
	"manage/setting"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const defaultBufferSize = 16 * 1024

type Call struct {
	msg  *protocol.Message
	Done chan *Call
}

type Client struct {
	sync.Mutex
	Id int64
	net.Conn
	quit     bool
	exitchan chan struct{}
	// reading/writing interfaces
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	Timeout   time.Duration
	sendqueue chan *protocol.Message
	Addr      string
	Port      int
	lenBuf    [4]byte
	lenSlice  []byte
	pending   map[uint64]*Call
	seq       uint64
}

func newClient(id int64, conn net.Conn) *Client {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &Client{
		Id:        id,
		Conn:      conn,
		Reader:    bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:    bufio.NewWriterSize(conn, defaultBufferSize),
		Timeout:   time.Duration(setting.IdleTimeout) * time.Second,
		sendqueue: make(chan *protocol.Message, 32),
		exitchan:  make(chan struct{}),
		Addr:      addr,
		Port:      int(p),
		pending:   make(map[uint64]*Call),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *Client) SendMessage(msg *protocol.Message) bool {
	if c.quit {
		return false
	}

	seq := atomic.AddUint64(&c.seq, 1)
	msg.Header = msg.Header[:8]
	binary.LittleEndian.PutUint64(msg.Header, seq)
	c.sendqueue <- msg //消息太多的情况可能会阻塞
	return true
}

func (c *Client) Call(msg *protocol.Message) *Call {
	if c.quit {
		return nil
	}
	seq := atomic.AddUint64(&c.seq, 1)
	msg.Header = msg.Header[:8]
	binary.LittleEndian.PutUint64(msg.Header, seq)
	c.Lock()
	call := &Call{}
	call.Done = make(chan *Call, 1)
	c.pending[seq] = call
	c.Unlock()
	c.sendqueue <- msg
	return call
}

func (c *Client) Response(seq uint64, msg *protocol.Message) {
	c.Lock()
	if call, has := c.pending[seq]; has {
		call.msg = msg.Dup()
		call.Done <- call
		delete(c.pending, seq)
	}
	c.Unlock()
}

func (c *Client) Quit() {
	if !c.quit {
		c.quit = true
		close(c.exitchan)
		c.Close()
	}
}
