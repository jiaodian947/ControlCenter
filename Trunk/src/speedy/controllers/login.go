package controllers

import (
	"encoding/json"
	"log"
	"speedy/models"
	"speedy/modules/auth"
)

type Login struct {
	BaseRouter
}

func (l *Login) Post() {
	var login map[string]string

	var reply Reply
	reply.Status = 500
	log.Println(l.Ctx.Input.RequestBody)
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &login); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}

	user, err := models.GetUserByName(login["username"])
	if err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}

	auth.LoginUser(l.Ctx, user, true)
	in_white_list := auth.LoginClientWhiteList(l.Ctx)
	if !in_white_list {
		reply.Data = "clent ip is not in whiteList"
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"name": user.Name,
	}
	l.Data["json"] = &reply
	l.ServeJSON()

}
