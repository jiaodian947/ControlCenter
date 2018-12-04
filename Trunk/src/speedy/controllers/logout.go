package controllers

import "speedy/modules/auth"

type Logout struct {
	BaseRouter
}

func (l *Logout) Post() {
	var reply Reply
	reply.Status = 200
	if l.IsLogin {
		auth.LogoutUser(l.Ctx)
	}
	l.Data["json"] = &reply
	l.ServeJSON()
}
