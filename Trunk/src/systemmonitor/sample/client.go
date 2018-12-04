package sample

import (
	"bufio"
	"net"
	"strconv"
	"systemmonitor/protocol"
	"time"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	net.Conn
	quit     bool
	exitchan chan struct{}
	// reading/writing interfaces
	Reader            *bufio.Reader
	Writer            *bufio.Writer
	HeartbeatInterval time.Duration
	sendqueue         chan *protocol.Message
	Addr              string
	Port              int
	lenBuf            [4]byte
	lenSlice          []byte
}

func newClient(conn net.Conn) *Client {

	addr, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	p, _ := strconv.ParseInt(port, 10, 32)

	c := &Client{
		Conn:      conn,
		Reader:    bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:    bufio.NewWriterSize(conn, defaultBufferSize),
		sendqueue: make(chan *protocol.Message, 32),
		exitchan:  make(chan struct{}),
		Addr:      addr,
		Port:      int(p),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *Client) Shutdown() {
	if !c.quit {
		c.quit = true
		c.Close()
		close(c.exitchan)
	}
}
