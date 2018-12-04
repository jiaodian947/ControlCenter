package api

import (
	"usercenter/controllers"
	"usercenter/modules/auth"
	"usercenter/modules/models"
	"usercenter/modules/utils"

	"github.com/astaxie/beego"
)

type GameAuthController struct {
	controllers.BaseRouter
}

const (
	ERR_SUCCESS                        = 200
	MGS_ERR_VALIDATE_FAILED            = 21129
	CAS_ERR_NO_ACCOUNT                 = 51001 // 帐号不存在
	CAS_ERR_ACCOUNT_PSWD               = 51002
	CAS_ERR_ACCOUNT_FORBID             = 51013
	CAS_ERR_CDKEY_INVALID              = 54001 // 激活码无效
	CAS_ERR_ACCOUNT_ACTIVED            = 54002 // 帐号已经激活
	CAS_ERR_ACTIVE_FAILED              = 54003 // 账号激活失败
	CAS_ERR_ACCOUNT_EXIST              = 54005 // 用户名已经存在
	CAS_ERR_ACCOUNT_REG_FAILED         = 54006 // 注册账号失败
	CAS_ERR_ACCOUNT_NAME_OR_PASS_EMPTY = 54007 // 用户名或者密码为空
	CAS_ERR_ACCOUNT_TOO_LONG           = 54008 // 用户名过长
)

type Data struct {
	Account string
	Code    string
}

type LoginResult struct {
	ErrCode int
	Data    Data
}

func (c *GameAuthController) Get() {
	token := c.GetString("token")
	account := c.GetString("account")
	res := &VerifyResult{}
	res.Result = "error"
	res.Error = "validate failed"
	beego.Info("verify:", account, token)
	if token == "" || account == "" {
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if !utils.VerifyToken(account, token) {
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Result = "ok"
	res.Error = ""
	c.Data["json"] = res
	c.ServeJSON()
	return
}

func (c *GameAuthController) Post() {
	username := c.GetString("user")
	pass := c.GetString("pass")
	res := &LoginResult{}
	if username == "" || pass == "" {
		res.ErrCode = CAS_ERR_ACCOUNT_PSWD
		c.Data["json"] = res
		c.ServeJSON()
		beego.Info("auth failed,", username)
		return
	}

	var user models.User
	if !auth.VerifyUser(&user, username, pass) {
		res.ErrCode = CAS_ERR_ACCOUNT_PSWD
		c.Data["json"] = res
		c.ServeJSON()
		beego.Info("auth failed,", username)
		return
	}

	if !user.IsActive {
		res.ErrCode = CAS_ERR_ACCOUNT_FORBID
		c.Data["json"] = res
		c.ServeJSON()
		beego.Info("account forbid,", username)
		return
	}

	res.ErrCode = ERR_SUCCESS
	res.Data.Account = user.UserName
	res.Data.Code = utils.CreateToken(user.UserName, 5, 0)
	beego.Info(username, "auth succeed,", "token:", res.Data.Code)
	c.Data["json"] = res
	c.ServeJSON()
}
