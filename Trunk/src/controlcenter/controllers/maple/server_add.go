package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/maple"
	"controlcenter/modules/models"
	"fmt"
	"strconv"
)

type ServerAddController struct {
	controllers.BaseAdmin
}

func (c *ServerAddController) Get() {
	c.Data["active"] = "server"
	c.Data["servers"] = nil
	c.Data["gamename"] = ""
	c.Data["districtname"] = ""
	c.TplName = "server_add.html"

	id := c.Ctx.Input.Param(":id1")
	id2 := c.Ctx.Input.Param(":id2")
	c.Data["gameid"] = id
	c.Data["districtid"] = id2
	gameid, err := strconv.Atoi(id)
	districtid, err1 := strconv.Atoi(id2)
	if err != nil || err1 != nil {
		return
	}

	var game models.ServerGame
	game.Id = gameid
	if game.Read("Id") != nil {
		return
	}
	var district models.ServerDistrict
	district.Id = districtid
	if district.Read("Id") != nil {
		return
	}

	c.Data["gamename"] = game.Name
	c.Data["districtname"] = district.DistrictName
	c.Data["title"] = fmt.Sprintf("%s-%s服务器增加", game.Name, district.DistrictName)
}

func (c *ServerAddController) Post() {
	c.Data["title"] = "服务器管理"
	c.Data["active"] = "server"
	c.TplName = "server_add.html"

	var form maple.ServerForm
	if c.ValidFormSets(&form) == false {
		c.Data["err"] = "args error"
		return
	}

	var info models.ServerInfo
	info.Id = form.Id
	if form.Id != 0 && info.Read("Id") == nil {
		c.Data["err"] = "id conflict"
		return
	}
	info.Id = form.Id
	info.DistrictId = form.DistrictId
	info.GameId = form.GameId
	info.ServerName = form.ServerName
	info.ServerType = form.ServerType
	info.ServerStatus = form.ServerStatus
	info.PlayerMaxCount = form.MaxPlayer
	info.ServerIp = form.ServerIp
	info.ServerPort = form.ServerPort
	info.Comment = form.Comment

	if info.Insert() != nil {
		c.Data["err"] = "增加服务器失败"
		return
	}
	c.Ctx.Redirect(302, fmt.Sprintf("/game/%d/%d", form.GameId, form.DistrictId))
}
