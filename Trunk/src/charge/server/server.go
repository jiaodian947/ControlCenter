package server

import (
	"charge/setting"
	"charge/util"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// 主服务器
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
	servers          map[int64]*ServerInfo
	useracchash      map[int]int64
	userloginhash    map[string]int64
	access           *access
	tradingHall      *TradingHall
}

func New() *Server {
	s := &Server{}
	s.clients = make(map[int64]*Client, 128)
	s.servers = make(map[int64]*ServerInfo, 128)
	s.useracchash = make(map[int]int64)
	s.userloginhash = make(map[string]int64)
	return s
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

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", setting.AppHost, setting.AppPort))
	if err != nil {
		s.log.Fatalf("tcp listen port: %d, err :%s", setting.AppPort, err.Error())
	}
	s.tcpListener = listener
	ctx := &Context{s}
	tcp := &tcpServer{ctx}
	s.waitGroup.Wrap(func() {
		util.TCPServer(s.tcpListener, tcp, s.log)
	})

	s.access = NewAccess(ctx, setting.WorkThreads, setting.QueueLen)
	s.access.Start()
	s.tradingHall = NewTradingHall(ctx, setting.PerChannelLen)
	for k := range transactions {
		s.tradingHall.CreatePlatform(k)
	}

	s.tradingHall.StartAll()
	s.log.Printf("server started")
}

// 增加一个客户端连接
func (s *Server) AddClient(client *Client) bool {
	s.Lock()
	if _, exist := s.clients[client.Id]; exist {
		s.Unlock()
		return false
	}

	s.clients[client.Id] = client
	s.Unlock()
	return true
}

// 查找客户端连接
func (s *Server) FindClient(clientid int64) *Client {
	s.RLock()
	c, _ := s.clients[clientid]
	s.RUnlock()
	return c
}

// 移除客户端连接
func (s *Server) RemoveClient(clientid int64) {
	s.Lock()
	delete(s.clients, clientid)
	s.Unlock()
}

// 增加一个游戏服务器
func (s *Server) AddServer(server *ServerInfo) bool {
	s.Lock()
	if _, exist := s.servers[server.Id]; exist {
		s.Unlock()
		return false
	}

	s.servers[server.Id] = server
	s.Unlock()
	return true
}

// 移除一个游戏服务器
func (s *Server) RemoveServer(serverid int64) {
	s.Lock()
	delete(s.servers, serverid)
	s.Unlock()
}

// 查找游戏服务器
func (s *Server) FindServer(serverid int64) *ServerInfo {
	s.RLock()
	if srv, exist := s.servers[serverid]; exist {
		s.RUnlock()
		return srv
	}
	s.RUnlock()
	return nil
}

// 增加一个交易
func (s *Server) AddTrader(w Trader) {
	s.tradingHall.AddTrader(w)
}

// 退出
func (s *Server) Exit() {
	s.log.Printf("server quit")
	if s.httpListener != nil {
		s.httpListener.Close()
	}
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}
	s.access.Close()
	s.access.Wait()
	s.tradingHall.Shutdown()
	s.waitGroup.Wait()
	if s.logFile != nil {
		s.logFile.Close()
	}
}
