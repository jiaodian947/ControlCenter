package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"manage/models"
	"manage/protocol"
	"manage/util"
	"net"
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
)

var (
	ServerApp *Server
)

type Server struct {
	sync.RWMutex
	log              *log.Logger
	logFile          *os.File
	tcpListener      net.Listener
	httpListener     net.Listener
	waitGroup        util.WaitGroupWrapper
	clientIDSequence int64
	msgIdSequence    int64
	userIdSequence   int64
	clients          map[int64]*Client
	servers          map[int]*GameServer
}

func init() {
	ServerApp = &Server{}
	ServerApp.clients = make(map[int64]*Client, 128)
	ServerApp.servers = make(map[int]*GameServer, 128)
}

func Run() {
	ServerApp.Main()
}

func Exit() {
	ServerApp.Exit()
}

func (s *Server) Main() {

	fileName := "log.log"
	logfile, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("open file error !")
	}
	w := io.MultiWriter(logfile, os.Stdout)
	s.logFile = logfile
	// 创建一个日志对象
	s.log = log.New(w, "", log.LstdFlags)

	s.LoadAll()
	s.log.Printf("server started")
}

func (s *Server) LoadAll() {
	//s.LoadAllServers()
}

func (s *Server) LoadAllServers() {
	var servers []*models.Server
	_, err := models.Servers().All(&servers)
	if err != nil && err != orm.ErrNoRows {
		s.log.Fatal(err)
		return
	}

	for k := range servers {
		gs := NewGameServer(servers[k])
		s.AddServer(gs)
	}
}

func (s *Server) AddClient(client *Client) bool {
	s.Lock()
	defer s.Unlock()
	if _, exist := s.clients[client.Id]; exist {
		return false
	}

	s.clients[client.Id] = client
	return true
}

func (s *Server) FindClient(clientid int64) *Client {
	s.RLock()
	defer s.RUnlock()
	c, _ := s.clients[clientid]
	return c
}

func (s *Server) RemoveClient(clientid int64) {
	s.Lock()
	defer s.Unlock()
	delete(s.clients, clientid)
}

func (s *Server) AddServer(server *GameServer) bool {
	s.Lock()
	defer s.Unlock()
	if _, exist := s.servers[server.Info.Id]; exist {
		return false
	}

	s.servers[server.Info.Id] = server
	return true
}

func (s *Server) RemoveServer(id int) {
	s.Lock()
	defer s.Unlock()
	delete(s.servers, id)
}

func (s *Server) FindServer(id int) *GameServer {
	s.RLock()
	defer s.RUnlock()
	if srv, exist := s.servers[id]; exist {
		return srv
	}
	return nil
}

func FindServerByServerId(gameid, serverid int) *GameServer {
	return ServerApp.FindServerByServerId(gameid, serverid)
}

func (s *Server) FindServerByServerId(gameid, serverid int) *GameServer {
	s.RLock()
	for _, v := range s.servers {
		if v.Info.GameId == gameid && v.Info.ServerId == serverid {
			s.RUnlock()
			return v
		}
	}
	s.RUnlock()

	server := &models.Server{}
	server.GameId = gameid
	server.ServerId = serverid
	if err := server.Read("GameId", "ServerId"); err != nil {
		return nil
	}
	gs := NewGameServer(server)
	if !s.AddServer(gs) { // 已经存在了
		return s.servers[gs.Info.Id]
	}
	return gs
}

func SendMessage(gameid, serverid int, msg *protocol.Message, need_response bool) (*protocol.Message, error) {
	return ServerApp.SendMessage(gameid, serverid, msg, need_response)
}

func (s *Server) SendMessage(gameid, serverid int, msg *protocol.Message, need_response bool) (*protocol.Message, error) {
	gs := s.FindServerByServerId(gameid, serverid)
	if gs == nil {
		return nil, errors.New("server not found")
	}

	if !gs.Connected {
		if err := gs.Connect(); err != nil {
			return nil, err
		}
	}

	c := s.FindClient(gs.Client)
	if c == nil {
		return nil, errors.New("client not found")
	}

	if need_response {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		ch := c.Call(msg).Done
		select {
		case call := <-ch:
			return call.msg, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("time out")
		}
	} else {
		if !c.SendMessage(msg) {
			return nil, errors.New("client is quit")
		}
	}

	return nil, nil
}

func (s *Server) Exit() {
	s.log.Printf("server quit")
	if s.httpListener != nil {
		s.httpListener.Close()
	}
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}
	s.waitGroup.Wait()
	if s.logFile != nil {
		s.logFile.Close()
	}
}
