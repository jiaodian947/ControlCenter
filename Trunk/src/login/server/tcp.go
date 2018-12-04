package server

import (
	"login/protocol"
	"net"
)

type tcpServer struct {
	ctx *Context
}

func (p *tcpServer) Handle(clientConn net.Conn) {
	p.ctx.server.log.Printf("TCP: new client(%s)", clientConn.RemoteAddr())
	var prot protocol.Protocol
	prot = &TextProtocol{p.ctx}
	err := prot.IOLoop(clientConn)
	if err != nil {
		p.ctx.server.log.Printf("ERROR: client(%s) - %s", clientConn.RemoteAddr(), err)
		return
	}
}
