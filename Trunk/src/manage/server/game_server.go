package server

import (
	"database/sql"
	"fmt"
	"log"
	"manage/models"
	"manage/protocol"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type GameServer struct {
	sync.Mutex
	Info      *models.Server
	Client    int64
	Connected bool
	DB        *sql.DB
}

func NewGameServer(info *models.Server) *GameServer {
	gs := &GameServer{}
	gs.Info = info
	gs.Connected = false
	dbinfos := strings.Split(info.GameDb, ":")
	if len(dbinfos) == 6 {
		//sa:abc@tcp(192.168.1.52:3306)/sininm_game?charset=utf8&loc=UTC
		ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", dbinfos[1], dbinfos[2], dbinfos[0], dbinfos[4], dbinfos[3], dbinfos[5])
		db, err := sql.Open("mysql", ds)
		if err == nil {
			gs.DB = db
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
