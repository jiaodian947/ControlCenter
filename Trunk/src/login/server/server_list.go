package server

import (
	"login/models"
	"login/protocol"
	"login/setting"

	"github.com/astaxie/beego/orm"
)

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

func RegisterServer(client *Client, msg *protocol.VarMessage) {
	if client.ctx.server.FindServer(client.Id) != nil {
		client.ctx.server.log.Println("server already register")
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
		models.Accounts().Filter("ServerId", server_id).Update(orm.Params{"ServerId": 0})
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
		client.ctx.server.log.Printf("server has not normal closed")
	}

	m.ServerIp = s.Address
	m.Logged = 1

	if m.Update() != nil {
		client.ctx.server.log.Printf("update to server error")
		return
	}

	out := protocol.NewVarMsg(3)
	out.AddString(user_id)
	out.AddString("register")
	s.GameId = m.GameId
	if client.ctx.server.AddServer(s) {
		client.ctx.server.log.Printf("register server: %d %s:%d", s.ServerId, s.Address, s.Port)
		out.AddInt(1)
	} else {
		out.AddInt(0)
	}

	client.SendMessage(out)
}

func UnregisterServer(client *Client, msg *protocol.VarMessage) {

	user_id := msg.StringVal(0)

	s := client.ctx.server.FindServer(client.Id)
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

func RemoveServer(client *Client, s *ServerInfo) {
	if s == nil {
		s = client.ctx.server.FindServer(client.Id)
		if s == nil {
			return
		}
	}

	s.ctx.server.RemoveServer(s.Id)
	s.ctx.server.log.Printf("remove server: %d %s:%d", s.ServerId, s.Address, s.Port)
}
