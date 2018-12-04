package api

import (
	"encoding/base64"
	"usercenter/modules/auth"
	"usercenter/modules/models"

	"strings"

	"github.com/astaxie/beego"
)

type ProfileRouter struct {
	beego.Controller
}

type VerifyResult struct {
	Result string
	Error  string
}

func (c *ProfileRouter) Verify() {
	s := c.GetString("token")
	token, err := base64.StdEncoding.DecodeString(s)

	ret := &VerifyResult{}
	if err != nil {
		ret.Result = "error"
		ret.Error = err.Error()
		c.Data["json"] = ret
		c.ServeJSON()
		return
	}

	info := strings.Split(string(token), " ")
	if len(info) != 2 {
		ret.Result = "error"
		ret.Error = "token parse error"
		c.Data["json"] = ret
		c.ServeJSON()
		return
	}

	name, pass := info[0], info[1]
	beego.Info(name, pass)
	var user models.User
	if auth.VerifyUser(&user, name, pass) {
		ret.Result = "ok"
		c.Data["json"] = ret
		c.ServeJSON()
		return
	}

	ret.Result = "error"
	ret.Error = "name or pass error"
	c.Data["json"] = ret
	c.ServeJSON()
}
