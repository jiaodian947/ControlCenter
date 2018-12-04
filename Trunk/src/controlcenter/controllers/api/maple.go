package api

import (
	"controlcenter/controllers/maple"
	"controlcenter/modules/models"

	"github.com/astaxie/beego"
)

type District struct {
	DistrictId   int
	DistrictName string
}

type Server struct {
	DistrictId       int
	ServerId         int
	ServerName       string
	ServerStatus     int
	ServerPlayers    int
	ServerMaxPlayers int
	ServerIp         string
	ServerPort       int
}

type ServerList struct {
	Ok        bool
	Err       string
	Districts []District
	Servers   []Server
	MyServers []models.UserServerInfo
}

type MapleAPI struct {
	beego.Controller
}

func (c *MapleAPI) PrePullByAccount() {
}

func (c *MapleAPI) PrePull() {
	var gameid int
	var seckey string
	var account string
	var filter int
	var err error
	sl := ServerList{}
	if gameid, err = c.GetInt(":gameid"); err != nil {
		if gameid, err = c.GetInt("gameid"); err != nil {
			sl.Ok = false
			sl.Err = err.Error()
			c.Data["json"] = &sl
			c.ServeJSON()
			return
		}
	}

	seckey = c.GetString(":seckey")
	if seckey == "" {
		seckey = c.GetString("seckey")
	}

	if seckey == "" {
		sl.Ok = false
		sl.Err = "seckey error"
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	account = c.GetString(":account")
	if account == "" {
		account = c.GetString("account")
	}

	if filter, err = c.GetInt(":filter"); err != nil {
		filter, err = c.GetInt("filter")
	}

	filter--
	if account == "" {
		c.Pull(gameid, seckey, filter)
	} else {
		c.PullByAccount(gameid, seckey, account, filter)
	}
}

func (c *MapleAPI) Pull(gameid int, seckey string, filter int) {
	var err error
	sl := ServerList{}
	var game models.ServerGame
	game.Id = gameid
	game.SecretKey = seckey
	if err = game.Read("Id", "SecretKey"); err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	var districts []*models.ServerDistrict
	_, err = models.Districts().Filter("GameId", gameid).All(&districts)
	if err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	var servers []*models.ServerInfo

	_, err = models.Servers().Filter("GameId", gameid).All(&servers)
	if err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	sl.Ok = true

	sl.Districts = make([]District, 0, len(districts))

	filters := make(map[int]struct{})
	for _, v := range districts {
		if filter >= 0 && v.Group != filter {
			continue
		}
		var d District
		d.DistrictId = v.Id
		d.DistrictName = v.DistrictName
		sl.Districts = append(sl.Districts, d)
		filters[v.Id] = struct{}{}
	}

	sl.Servers = make([]Server, 0, len(servers))
	for _, v := range servers {
		if _, has := filters[v.DistrictId]; !has {
			continue
		}
		var s Server
		s.DistrictId = v.DistrictId
		s.ServerId = v.Id
		s.ServerName = v.ServerName
		s.ServerIp = v.ServerIp
		s.ServerPort = v.ServerPort
		s.ServerStatus = v.ServerStatus
		s.ServerPlayers = v.PlayerCount
		s.ServerMaxPlayers = v.PlayerMaxCount
		sl.Servers = append(sl.Servers, s)
	}
	c.Data["json"] = &sl
	c.ServeJSON()
}

func (c *MapleAPI) PullByAccount(gameid int, seckey string, account string, filter int) {
	var err error
	sl := ServerList{}
	var game models.ServerGame
	game.Id = gameid
	game.SecretKey = seckey
	if err = game.Read("Id", "SecretKey"); err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	var districts []*models.ServerDistrict
	_, err = models.Districts().Filter("GameId", gameid).All(&districts)
	if err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	var servers []*models.ServerInfo

	_, err = models.Servers().Filter("GameId", gameid).All(&servers)
	if err != nil {
		sl.Ok = false
		sl.Err = err.Error()
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	sl.Ok = true

	sl.Districts = make([]District, 0, len(districts))
	filters := make(map[int]struct{})
	for _, v := range districts {
		if filter >= 0 && v.Group != filter {
			continue
		}
		var d District
		d.DistrictId = v.Id
		d.DistrictName = v.DistrictName
		sl.Districts = append(sl.Districts, d)
		filters[v.Id] = struct{}{}
	}

	sl.Servers = make([]Server, 0, len(servers))
	for _, v := range servers {
		if _, has := filters[v.DistrictId]; !has {
			continue
		}
		var s Server
		s.DistrictId = v.DistrictId
		s.ServerId = v.Id
		s.ServerName = v.ServerName
		s.ServerIp = v.ServerIp
		s.ServerPort = v.ServerPort
		s.ServerStatus = v.ServerStatus
		s.ServerPlayers = v.PlayerCount
		s.ServerMaxPlayers = v.PlayerMaxCount
		sl.Servers = append(sl.Servers, s)
	}

	dbsession := maple.GetDBSession()
	if dbsession == nil {
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	defer dbsession.Close()

	DB := dbsession.DB(maple.MAPLEDB)
	collection := DB.C(models.UserInfoCollection)
	user := &models.UserInfo{}
	user.Account = account
	if err := user.Read(collection); err != nil {
		c.Data["json"] = &sl
		c.ServeJSON()
		return
	}

	usi := make([]models.UserServerInfo, 0, len(user.ServerInfos))
	for _, v := range user.ServerInfos {
		if _, has := filters[v.DistrictId]; !has {
			continue
		}

		usi = append(usi, v)

	}
	sl.MyServers = usi
	c.Data["json"] = &sl
	c.ServeJSON()
}
