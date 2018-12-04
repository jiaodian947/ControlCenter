package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"login/protocol"
	"login/util"
	"net/http"
	"time"
)

func (u *User) name_password_base64() string {
	str := fmt.Sprintf("%s %s", u.Account, u.Password)
	return "token=" + base64.StdEncoding.EncodeToString([]byte(str))
}

func (u *User) name_login_string() string {
	str := fmt.Sprintf("account=%s&token=%s", u.Account, u.LoginString)
	return str
}

func (u *User) EncodeRequest(url string, enc_func string) string {
	var enc string
	switch enc_func {
	case "name_password_base64":
		enc = u.name_password_base64()
	case "name_login_string":
		enc = u.name_login_string()
	}
	if enc == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", url, enc)
}

func (u *User) AuthRequest(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		u.ctx.server.log.Println("(AuthRequest) http failed,", err, url)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		u.ctx.server.log.Println("(AuthRequest) http read failed,", err, url)
		return false
	}
	var result AuthResult
	if err := json.Unmarshal(body, &result); err != nil {
		u.ctx.server.log.Println("(AuthRequest) decode result failed,", err, url)
		return false
	}

	u.ctx.server.log.Println("(AuthRequest) result,", result.Result, result.Error, url)
	return result.Result == "ok"
}

func (u *User) RequestSelf(authurl string, enc_func string) {
	url := u.EncodeRequest(authurl, enc_func)

	msg := protocol.NewVarMsg(14)
	msg.AddString(u.UserId)
	msg.AddString("login")
	msg.AddInt(1)
	errcode := 0
	if url != "" {
		if u.AuthRequest(url) {
			acc, ret := u.ReadChargeUser()
			if ret {
				if acc.ServerId == u.ServerId { //已经登录了(可能计费服务器挂了，先登出原来的帐号)
					acc.ServerId = 0
					acc.LogonId = ""
					acc.Update("LogonId", "ServerId", "LastlogoutTime", "LastExec")
				}

				validtime := acc.ValidTime.Sub(time.Now()).Seconds()
				if acc.ServerId == 0 &&
					acc.Status == 1 &&
					validtime <= 1 {
					u.dbuser = acc
					acc.LogonId = util.UUID()
					acc.ServerId = u.ServerId
					acc.LastlogAddr = u.IpAddr
					acc.LastlogTime = time.Now()
					acc.Update("LogonId", "LastlogTime", "ServerId", "LastlogAddr", "LastExec")
					msg.AddString(acc.Account)
					msg.AddInt(1) //result
					msg.AddInt(acc.Id)
					msg.AddString(acc.LogonId)
					msg.AddString("") //user name
					msg.AddInt(acc.GmLevel)
					msg.AddString("") //password
					msg.AddInt(acc.Points)
					msg.AddDouble(0)  //limit
					msg.AddInt(1)     //is free
					msg.AddString("") // acc_info
					msg.AddInt(0)     //login type

					c := u.ctx.server.FindClient(u.ConnId)
					if c != nil {
						c.SendMessage(msg)
						u.ctx.server.log.Println("user login success,", acc.Account)
						u.ctx.server.UpdateUserHash(u.Index)
						return
					}
					u.ctx.server.log.Println("user not found", acc.Account)
					return
				} else {
					if validtime > 1 {
						errcode = 20105 // 在一段时间内禁止登陆
					} else if acc.ServerId != 0 {
						errcode = 51011 // 此帐号已在分区内其他服务器登录
					} else if acc.Status == 0 {
						errcode = 51003 // 此帐号已被冻结
					}
					u.ctx.server.log.Printf("user login failed, account:%s, serverid: %d, need:0, status:%d, need:1, validtime:%d, need <=0", acc.Account, acc.ServerId, acc.Status, int(validtime))
				}
			}
		}
	}

	msg.AddString(u.Account)
	msg.AddInt(errcode) //result

	c := u.ctx.server.FindClient(u.ConnId)
	if c != nil {
		c.SendMessage(msg)
		u.ctx.server.log.Println("user login failed", u.Account)
	}

	u.ctx.server.RemoveUserByIndex(u.Index, u.ConnId)
}
