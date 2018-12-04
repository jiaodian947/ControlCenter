package api

import (
	"usercenter/controllers"
	"usercenter/modules/auth"
	"usercenter/modules/models"

	"github.com/beego/i18n"
)

type GameReg struct {
	controllers.BaseRouter
}
type RegResult struct {
	ErrCode int
	Err     string
}

func (r *GameReg) Register() {
	name := r.GetString("name")
	pass := r.GetString("pass")

	ret := new(RegResult)
	if name == "" || pass == "" {
		ret.ErrCode = CAS_ERR_ACCOUNT_NAME_OR_PASS_EMPTY
		ret.Err = "account or pass empty"
		r.Data["json"] = ret
		r.ServeJSON()
		return
	}

	if len(name) > 30 {
		ret.ErrCode = CAS_ERR_ACCOUNT_TOO_LONG
		ret.Err = "account too long"
		r.Data["json"] = ret
		r.ServeJSON()
		return
	}

	var user models.User
	user.UserName = name
	if user.Read("UserName") == nil {
		ret.ErrCode = CAS_ERR_ACCOUNT_EXIST
		ret.Err = "account exist"
		r.Data["json"] = ret
		r.ServeJSON()
		return
	}
	user.Lang = i18n.IndexLang(r.Lang)
	if err := auth.RegisterUser(&user, name, name, "", pass); err != nil {
		ret.ErrCode = CAS_ERR_ACCOUNT_REG_FAILED
		ret.Err = err.Error()
		r.Data["json"] = ret
		r.ServeJSON()
		return
	}

	ret.ErrCode = ERR_SUCCESS
	r.Data["json"] = ret
	r.ServeJSON()
}
