package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/models"
)

type GameController struct {
	controllers.BaseAdmin
}

func (c *GameController) Get() {
	c.TplName = "game_manager.html"
	c.Data["title"] = "服务器管理"
	c.Data["active"] = "server"

	var games []*models.ServerGame
	_, err := models.Games().All(&games)
	if err != nil {
		c.Data["games"] = nil
		return
	}
	c.Data["games"] = games
}

func (c *GameController) DeleteGame() {

	var game models.ServerGame
	var err error
	game.Id, err = c.GetInt(":id")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}
	if game.Read("Id") != nil { //没有找到
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var server models.ServerInfo
	server.GameId = game.Id
	server.Delete("GameId")

	var district models.ServerDistrict
	district.GameId = game.Id
	district.Delete("GameId")

	game.Delete("Id")

	c.Redirect(c.Ctx.Request.Referer(), 302)
}
