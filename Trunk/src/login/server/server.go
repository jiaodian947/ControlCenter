package server

import (
	"fmt"
	"io"
	"log"
	"login/setting"
	"login/util"
	"net"
	"os"
	"sync"
	"sync/atomic"
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
	servers          map[int64]*ServerInfo
	users            map[int64]*User
	useracchash      map[int]int64
	userloginhash    map[string]int64
	access           *access
}

func New() *Server {
	s := &Server{}
	s.clients = make(map[int64]*Client, 128)
	s.servers = make(map[int64]*ServerInfo, 128)
	s.users = make(map[int64]*User, 2048)
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
	s.log.Printf("server started")
}

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

func (s *Server) FindClient(clientid int64) *Client {
	s.RLock()
	c, _ := s.clients[clientid]
	s.RUnlock()
	return c
}

func (s *Server) RemoveClient(clientid int64) {
	s.Lock()
	delete(s.clients, clientid)
	s.Unlock()
}

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

func (s *Server) RemoveServer(serverid int64) {
	s.Lock()
	delete(s.servers, serverid)
	s.Unlock()
}

func (s *Server) FindServer(serverid int64) *ServerInfo {
	s.RLock()
	if srv, exist := s.servers[serverid]; exist {
		s.RUnlock()
		return srv
	}
	s.RUnlock()
	return nil
}

func (s *Server) AddUser(connid int64, account string, password string, address string, port int, serverid int) *User {
	user := &User{}
	user.ConnId = connid
	user.ServerId = serverid
	user.Account = account
	user.Password = password
	user.IpAddr = address
	user.Port = port
	user.Index = atomic.AddInt64(&s.userIdSequence, 1)
	s.Lock()
	s.users[user.Index] = user
	s.Unlock()
	return user
}

func (s *Server) UpdateUserHash(index int64) {
	s.Lock()
	if u, has := s.users[index]; has {
		if u.dbuser != nil {
			if u.dbuser.Id != 0 {
				s.useracchash[u.dbuser.Id] = index
			}
			if u.dbuser.LogonId != "" {
				s.userloginhash[u.dbuser.LogonId] = index
			}
		}
	}
	s.Unlock()
}

func (s *Server) RemoveUserByIndex(index int64, connid int64) {
	s.Lock()
	if u, has := s.users[index]; has {
		if u.ConnId == connid {
			delete(s.users, index)
			if u.dbuser != nil {
				delete(s.useracchash, u.dbuser.Id)
				delete(s.userloginhash, u.dbuser.LogonId)
			}

		}
	}
	s.Unlock()
}

func (s *Server) RemoveAllUserByConnid(connid int64) {
	s.RLock()
	del := make([]int64, 0, 128)
	for _, u := range s.users {
		if u.ConnId == connid {
			del = append(del, u.Index)
		}
	}

	s.RUnlock()

	for _, index := range del {
		s.RemoveUserByIndex(index, connid)
	}
}

func (s *Server) GetUserByIndex(index int64) *User {
	s.RLock()
	if u, has := s.users[index]; has {
		s.RUnlock()
		return u
	}
	s.RUnlock()
	return nil
}

func (s *Server) GetUserByLogonId(id string) *User {
	var user *User
	s.RLock()
	if index, exist := s.userloginhash[id]; exist {
		user = s.users[index]
	}
	s.RUnlock()
	return user
}

func (s *Server) GetUserByAccId(acc_id int) *User {
	var user *User
	s.RLock()
	if index, exist := s.useracchash[acc_id]; exist {
		user = s.users[index]
	}
	s.RUnlock()
	return user
}

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
	s.waitGroup.Wait()
	if s.logFile != nil {
		s.logFile.Close()
	}
}
