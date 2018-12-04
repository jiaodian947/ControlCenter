package maple

import (
	"controlcenter/controllers"
	"controlcenter/modules/maple"
	"controlcenter/modules/models"
	"fmt"
	"strconv"
)

type ServerController struct {
	controllers.BaseAdmin
}

func (c *ServerController) Get() {

	c.Data["active"] = "server"
	c.Data["servers"] = nil
	c.Data["gamename"] = ""
	c.Data["districtname"] = ""
	c.TplName = "server_manager.html"

	id := c.Ctx.Input.Param(":id1")
	id2 := c.Ctx.Input.Param(":id2")
	c.Data["gameid"] = id
	c.Data["districtid"] = id2
	gameid, err := strconv.Atoi(id)
	districtid, err1 := strconv.Atoi(id2)
	if err != nil || err1 != nil {
		return
	}

	var gameinfo models.ServerGame
	gameinfo.Id = gameid
	if gameinfo.Read("Id") != nil {
		return
	}
	var districtinfo models.ServerDistrict
	districtinfo.Id = districtid
	if districtinfo.Read("Id") != nil {
		return
	}
	c.Data["gamename"] = gameinfo.Name
	c.Data["districtname"] = districtinfo.DistrictName
	c.Data["title"] = fmt.Sprintf("%s-%s服务器管理", gameinfo.Name, districtinfo.DistrictName)
	c.Data["active"] = "server"

	var servers []*models.ServerInfo
	curpage, err := c.GetInt("p")
	if err != nil {
		curpage = 1
	}

	curpage--
	pagelimit := 20

	count, err := models.Servers().Filter("district_id", districtid).Count()
	if err != nil {
		return
	}

	c.SetPaginator(pagelimit, count)
	startpos := curpage * pagelimit

	models.Servers().Filter("district_id", districtid).Limit(pagelimit, startpos).All(&servers)
	c.Data["servers"] = servers
}

func (c *ServerController) DeleteServer() {
	var server models.ServerInfo
	var err error
	server.Id, err = c.GetInt(":id3")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	if server.Read("Id") != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	server.Delete("Id")
	c.Redirect(c.Ctx.Request.Referer(), 302)
}

func (c *ServerController) ShowServer() {
	c.TplName = "server_edit.html"
	c.Data["title"] = ""
	c.Data["active"] = "server"

	var server models.ServerInfo
	var err error
	var gameid, districtid int
	if gameid, err = c.GetInt(":id1"); err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	if districtid, err = c.GetInt(":id2"); err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	server.Id, err = c.GetInt(":id3")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var game models.ServerGame
	game.Id = gameid
	if game.Read("Id") != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var district models.ServerDistrict
	district.Id = districtid
	if district.Read("Id") != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	if server.Read("Id") != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	c.Data["gamename"] = game.Name
	c.Data["districtname"] = district.DistrictName
	c.Data["server"] = server
	c.Data["title"] = fmt.Sprintf("编辑-%s", server.ServerName)

}

func (c *ServerController) UpdateServer() {
	c.TplName = "server_edit.html"
	c.Data["title"] = ""
	c.Data["active"] = "server"

	var form maple.ServerEditForm
	if c.ValidFormSets(&form) == false {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var server models.ServerInfo
	server.Id = form.Id
	server.DistrictId = form.DistrictId
	server.GameId = form.GameId
	server.ServerName = form.ServerName
	server.ServerType = form.ServerType
	server.ServerStatus = form.ServerStatus
	server.PlayerMaxCount = form.MaxPlayer
	server.ServerIp = form.ServerIp
	server.ServerPort = form.ServerPort
	server.Comment = form.Comment

	redurl := fmt.Sprintf("/game/%d/%d", form.GameId, form.DistrictId)
	if form.Id != form.ServerId { //id改了
		if server.Insert() != nil {
			c.Redirect(redurl, 302)
			return
		}
		var del models.ServerInfo
		del.Id = form.ServerId
		del.Delete("Id")
		c.Redirect(redurl, 302)
		return
	}

	server.Update()
	c.Redirect(redurl, 302)
}

func (c *ServerController) OpServer() {
	var server models.ServerInfo
	var err error

	server.Id, err = c.GetInt(":id3")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	if server.Read("Id") != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	op := c.GetString("op")
	switch op {
	case "open":
		server.ServerStatus = 1
	case "close":
		server.ServerStatus = 0
	case "maintain":
		server.ServerStatus = 2
	}

	server.Update("ServerStatus")

	c.Redirect(c.Ctx.Request.Referer(), 302)

}

func (c *ServerController) AllOpServer() {
	var gameid int
	var err error
	if gameid, err = c.GetInt(":id1"); err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
		return
	}

	var servers []*models.ServerInfo
	models.Servers().Filter("game_id", gameid).All(&servers)
	for _, s := range servers {
		if s.ServerStatus == 0 {
			continue
		}

		op := c.GetString("op")
		switch op {
		case "open":
			s.ServerStatus = 1
		case "close":
			s.ServerStatus = 0
		case "maintain":
			s.ServerStatus = 2
		}

		s.Update("ServerStatus")
	}

	c.Redirect(c.Ctx.Request.Referer(), 302)
}
