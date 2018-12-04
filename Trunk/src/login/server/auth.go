package server

import (
	"fmt"
	"login/models"
	"login/protocol"
	"login/setting"
	"time"

	"github.com/astaxie/beego/orm"
)

type AuthResult struct {
	Result string
	Error  string
}

type User struct {
	ctx         *Context
	Index       int64
	LoginType   int
	UserId      string //原样返回
	ConnId      int64
	ServerId    int
	GameId      int
	Account     string
	Password    string
	IpAddr      string
	Port        int
	Prefix      string
	LogonTime   time.Time
	LoginString string
	dbuser      *models.Account
}

func (u *User) AuthPostRequest(url string) bool {
	return true
}

func (u *User) RequestWithChannel(channel setting.ChannelInfo) {
	switch channel.ChannelType {
	case 0: //自营
		u.RequestSelf(channel.AuthUrl, channel.EncodeFunc)
	}
}

func (u *User) ReadChargeUser() (*models.Account, bool) {
	realacc := fmt.Sprintf("%s%s", u.Prefix, u.Account)
	account := &models.Account{}
	account.Account = realacc

	if account.Read("Account") == orm.ErrNoRows {
		account.From = u.LoginType
		account.Status = 1
		if err := account.Insert(); err != nil {
			u.ctx.server.log.Println("create user failed,", err)
			return nil, false
		}
		if account.Read("Account") == orm.ErrNoRows {
			u.ctx.server.log.Println("database user not found", realacc)
			return nil, false
		}
	}

	return account, true
}

func (u *User) Logout() {
	if u.dbuser != nil {
		onlinetime := int(time.Now().Sub(u.dbuser.LastlogTime).Minutes())
		u.dbuser.ServerId = 0
		u.dbuser.LogonId = ""
		u.dbuser.TotalTime += onlinetime
		u.dbuser.OnlineTime += onlinetime
		u.dbuser.Update("LogonId", "ServerId", "LastlogoutTime", "TotalTime", "OnlineTime", "LastExec")
	}

	u.ctx.server.RemoveUserByIndex(u.Index, u.ConnId)

	msg := protocol.NewVarMsg(4)
	c := u.ctx.server.FindClient(u.ConnId)
	if c != nil {
		msg.AddString(u.UserId)
		msg.AddString("logout")
		msg.AddInt(u.dbuser.Id)
		msg.AddInt(1)
		c.SendMessage(msg)
	}

	u.ctx.server.log.Println("user logout, ", u.Account)
}

func UserLogin(ctx *Context, msg *protocol.VarMessage) {
	srv := ctx.server.FindServer(msg.ConnId)
	if srv == nil {
		ctx.server.log.Println("(UserLogin) server not found")
		return
	}
	user_id := msg.StringVal(0)
	k := 2
	login_type := msg.IntVal(k)
	k++
	account := msg.StringVal(k)
	k++
	password := msg.StringVal(k)
	k++
	ip_addr := msg.StringVal(k)
	k++
	port := msg.IntVal(k)
	k++
	login_string := msg.StringVal(k)

	if login_type < 0 || login_type >= len(setting.Channels) {
		ctx.server.log.Println("(UserLogin) login type error,", login_type)
		return
	}

	channel := setting.Channels[login_type]

	user := ctx.server.AddUser(msg.ConnId, account, password, ip_addr, port, srv.ServerId)
	user.ctx = ctx
	user.UserId = user_id
	user.Prefix = channel.AccountPrefix
	user.LoginType = login_type
	user.LoginString = login_string
	//异步请求
	go user.RequestWithChannel(channel)
}

func UserLogout(ctx *Context, msg *protocol.VarMessage) {
	user_id := msg.StringVal(0)
	k := 2
	accid := msg.IntVal(k)
	k++
	logonid := msg.StringVal(k)
	user := ctx.server.GetUserByAccId(accid)
	if user == nil {
		user = ctx.server.GetUserByLogonId(logonid)
		if user == nil {
			ctx.server.log.Printf("user not found, accid:%d, logonid:%s", accid, logonid)
			return
		}
	}
	user.UserId = user_id
	go user.Logout()
}
