package controllers

import (
	"controlcenter/modules/auth"
	"controlcenter/modules/models"
)

type BaseAdmin struct {
	BaseRouter
}

func (this *BaseAdmin) NestPrepare() {
	if this.CheckLoginRedirect() {
		return
	}

	if !this.User.IsAdmin {
		auth.LogoutUser(this.Ctx)

		// write flash message
		this.FlashWrite("NotPermit", "true")

		this.Redirect("/login", 302)
		return
	}

	// current in admin page
	this.Data["IsAdmin"] = true

	var games []*models.ServerGame
	models.Games().All(&games)
	this.Data["GameData"] = games
}
