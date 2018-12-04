package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/maple"
	"controlcenter/modules/models"
	"fmt"
	"strconv"
)

type DistrictAddController struct {
	controllers.BaseAdmin
}

func (c *DistrictAddController) Get() {
	c.TplName = "district_add.html"
	c.Data["districts"] = nil
	c.Data["gamename"] = ""

	id := c.Ctx.Input.Param(":id")

	gameid, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	gameinfo := models.ServerGame{Id: gameid}

	if gameinfo.Read("Id") != nil {
		return
	}
	c.Data["gamename"] = gameinfo.Name
	c.Data["gameid"] = id
	c.Data["title"] = fmt.Sprintf("%s分区管理", gameinfo.Name)
	c.Data["active"] = "server"
}

func (c *DistrictAddController) Post() {
	c.TplName = "server_game_add.html"
	c.Data["title"] = "服务器管理"
	c.Data["active"] = "server"

	var form maple.DistrictForm
	if c.ValidFormSets(&form) == false {
		c.Data["err"] = "输入的信息不合法"
		return
	}

	dis := models.ServerDistrict{}
	dis.Id = form.Id
	if form.Id != 0 && dis.Read("Id") == nil {
		c.Data["err"] = "Id conflict"
		return
	}
	dis.Id = form.Id
	dis.GameId = form.GameId
	dis.DistrictName = form.DistrictName
	dis.Comment = form.Comment
	if dis.Insert() != nil {
		c.Data["err"] = "添加游戏失败"
		return
	}

	c.FlashRedirect(fmt.Sprintf("/game/%d", form.GameId), 302, "CreateSuccess")
}
