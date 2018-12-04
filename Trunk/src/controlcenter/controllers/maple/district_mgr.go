package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/maple"
	"controlcenter/modules/models"
	"fmt"
	"strconv"
)

type DistrictController struct {
	controllers.BaseAdmin
}

func (c *DistrictController) Get() {
	c.Data["districts"] = nil
	c.Data["gamename"] = ""
	c.Data["active"] = "server"
	c.Data["title"] = ""
	c.TplName = "district_manager.html"

	id := c.Ctx.Input.Param(":id")
	c.Data["gameid"] = id
	gameid, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	gameinfo := models.ServerGame{Id: gameid}

	if gameinfo.Read("Id") != nil {
		return
	}

	c.Data["gamename"] = gameinfo.Name
	c.Data["title"] = fmt.Sprintf("%s分区管理", gameinfo.Name)

	var districts []*models.ServerDistrict
	_, err = models.Districts().Filter("game_id", gameid).All(&districts)
	if err != nil {
		return
	}
	c.Data["districts"] = districts

	curpage, err := c.GetInt("p")
	if err != nil {
		curpage = 1
	}

	curpage--
	pagelimit := 20

	count, err := models.Servers().Filter("game_id", gameid).Count()
	if err != nil {
		return
	}

	c.SetPaginator(pagelimit, count)
	startpos := curpage * pagelimit

	var servers []*models.ServerInfo

	models.Servers().Filter("game_id", gameid).Limit(pagelimit, startpos).All(&servers)
	c.Data["servers"] = servers
}

func (c *DistrictController) DeleteDistrict() {
	var district models.ServerDistrict
	var err error
	district.Id, err = c.GetInt(":id2")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}
	if district.Read("Id") != nil { //没有找到
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var server models.ServerInfo
	server.DistrictId = district.Id
	server.Delete("DistrictId")

	district.Delete("Id")

	c.Redirect(c.Ctx.Request.Referer(), 302)
}

func (c *DistrictController) EditDistrict() {
	c.TplName = "district_edit.html"
	c.Data["title"] = ""
	c.Data["active"] = "server"

	id, err := c.GetInt(":id1")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}
	gameinfo := models.ServerGame{Id: id}
	if gameinfo.Read("Id") != nil {
		return
	}

	c.Data["gameid"] = id
	c.Data["gamename"] = gameinfo.Name

	var district models.ServerDistrict
	district.Id, err = c.GetInt(":id2")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}
	if district.Read("Id") != nil { //没有找到
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	c.Data["district"] = district
}

func (c *DistrictController) UpdateDistrict() {
	var district models.ServerDistrict
	var err error
	district.Id, err = c.GetInt(":id2")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}
	if district.Read("Id") != nil { //没有找到
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var form maple.DistrictEditForm
	if c.ValidFormSets(&form) == false {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	if form.Id != form.DistrictId { // 不能修改id
		c.Data["err"] = "不能修改Id"
		return
	}

	dis := models.ServerDistrict{}
	dis.Id = form.Id
	dis.GameId = form.GameId
	dis.DistrictName = form.DistrictName
	dis.Group = form.Group
	dis.Comment = form.Comment

	dis.Update()

	redurl := fmt.Sprintf("/game/%d", form.GameId)
	c.Redirect(redurl, 302)
}
