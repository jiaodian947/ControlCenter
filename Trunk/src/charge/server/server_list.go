package server

import (
	"charge/models"
	"charge/protocol"
	"charge/setting"
)

// 管理所有连接的服务器信息
type ServerInfo struct {
	ctx        *Context
	Id         int64
	Serial     int
	GameId     int
	ServerId   int
	ServerName string
	Address    string
	Md5        string
	Port       int
}

// 注册服务器
func RegisterServer(client *Client, msg *protocol.VarMessage) {
	if client.ctx.Server.FindServer(client.Id) != nil {
		client.ctx.Server.log.Println("server already register")
		return
	}
	k := 2
	user_id := msg.StringVal(0)
	game_type := msg.IntVal(k)
	k++
	server_id := msg.IntVal(k)
	k++
	server_name := msg.StringVal(k)
	k++
	md5_pswd := msg.StringVal(k)
	k++
	is_new_reg := msg.IntVal(k)

	s := &ServerInfo{
		ctx:        client.ctx,
		Id:         client.Id,
		GameId:     game_type,
		ServerId:   server_id,
		ServerName: server_name,
		Address:    client.Addr,
		Port:       client.Port,
		Md5:        md5_pswd,
	}

	if is_new_reg == 1 {
		m := models.Server{}
		m.ServerId = server_id
		if m.Read("ServerId") == nil {
			if m.Logged == 1 {
				m.Logged = 0
				m.Update("Logged")
			}
		}
	}

	m := models.Server{}
	m.ServerId = server_id
	if m.Read("ServerId") != nil {
		//调试时不存在则直接写入
		m.GameId = setting.GameId
		m.Insert()
		//client.ctx.server.log.Printf("server not found")
	}

	if m.Logged == 1 {
		client.ctx.Server.log.Printf("server has not normal closed")
	}

	m.ServerIp = s.Address
	m.Logged = 1

	if err := m.Update(); err != nil {
		client.ctx.Server.log.Println("update to server error", err)
		return
	}

	out := protocol.NewVarMsg(3)
	out.AddString(user_id)
	out.AddString("register")
	s.GameId = m.GameId
	if client.ctx.Server.AddServer(s) {
		client.ctx.Server.log.Printf("register server: %d %s:%d", s.ServerId, s.Address, s.Port)
		out.AddInt(1)
	} else {
		out.AddInt(0)
	}

	client.SendMessage(out)
}

// 注销服务器
func UnregisterServer(client *Client, msg *protocol.VarMessage) {

	user_id := msg.StringVal(0)

	s := client.ctx.Server.FindServer(client.Id)
	if s != nil {
		m := models.Server{}
		m.ServerId = s.ServerId
		if m.Read() == nil {
			m.ServerIp = ""
			m.Logged = 0
			m.Update("ServerIp", "Logged")
		}
		RemoveServer(client, s)
	}

	out := protocol.NewVarMsg(3)
	out.AddString(user_id)
	out.AddString("unregister")
	out.AddInt(1)
	client.SendMessage(out)
}

// 移除服务器
func RemoveServer(client *Client, s *ServerInfo) {
	if s == nil {
		s = client.ctx.Server.FindServer(client.Id)
		if s == nil {
			return
		}
	}

	s.ctx.Server.RemoveServer(s.Id)
	s.ctx.Server.log.Printf("remove server: %d %s:%d", s.ServerId, s.Address, s.Port)
}
