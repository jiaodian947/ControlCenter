package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/maple"
	"controlcenter/modules/models"
)

type GameAddController struct {
	controllers.BaseAdmin
}

func (c *GameAddController) Get() {
	c.Data["title"] = "游戏管理"
	c.Data["active"] = "server"
	c.TplName = "game_add.html"
}

func (c *GameAddController) Post() {
	c.Data["title"] = "游戏管理"
	c.Data["active"] = "server"
	c.TplName = "game_add.html"

	var form maple.GameForm
	if c.ValidFormSets(&form) == false {
		c.Data["err"] = "输入的信息不合法"
		return
	}

	game := models.ServerGame{}
	game.Id = form.Id
	if form.Id != 0 && game.Read("Id") == nil {
		c.Data["err"] = "Id conflict"
		return
	}
	game.Id = form.Id
	game.Name = form.GameName
	game.Comment = form.Comment
	if game.Insert() != nil {
		c.Data["err"] = "插入失败"
		return
	}

	c.FlashRedirect("/game", 302, "CreateSuccess")
}
