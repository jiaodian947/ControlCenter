package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"speedy/models"
	"speedy/protocol"
	"sync"
	"sync/atomic"
)

type GameServer struct {
	sync.Mutex
	Info      *models.Server
	Client    int64
	Connected bool
	DB        *sql.DB
	LogDB     *sql.DB
}

func NewGameServer(info *models.Server) *GameServer {
	gs := &GameServer{}
	gs.Info = info
	gs.Connected = false
	if info.GameDb != "" {
		db, err := sql.Open("mysql", info.GameDb)
		if err == nil {
			gs.DB = db
		}
	}
	if info.LogDb != "" {
		db, err := sql.Open("mysql", info.LogDb)
		if err == nil {
			gs.LogDB = db
		}
	}
	return gs
}

func (g *GameServer) Handle(client *Client, log *log.Logger) {
	log.Printf("TCP: new client(%s)", client.Conn.RemoteAddr())
	g.Connected = true
	var prot protocol.Protocol
	prot = &BinaryProtocol{client: client}
	err := prot.IOLoop()
	if err != nil {
		log.Printf("ERROR: client(%s) - %s", client.Conn.RemoteAddr(), err)
	}

	client.Quit()
	ServerApp.RemoveClient(client.Id)
	g.Connected = false
}

func (g *GameServer) Connect() error {
	g.Lock()
	defer g.Unlock()
	if g.Connected {
		return nil
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", g.Info.ServerIp, g.Info.ToolPort))
	if err != nil {
		return err
	}

	clientId := atomic.AddInt64(&ServerApp.clientIDSequence, 1)
	client := newClient(clientId, conn)
	if !ServerApp.AddClient(client) {
		return fmt.Errorf("add client(%d) error", clientId)
	}

	g.Client = clientId
	go g.Handle(client, ServerApp.log)

	return nil
}
