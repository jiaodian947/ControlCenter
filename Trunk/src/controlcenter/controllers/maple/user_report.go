package maple

import (
	"controlcenter/modules/models"
	"encoding/json"

	"github.com/astaxie/beego"

	mgo "gopkg.in/mgo.v2"
)

type RoleInfo struct {
	Account    string
	RoleName   string
	GameId     int
	DistrictId int
	ServerId   int
	Level      int
}

func ReportRoleInfo(handler *Handler) {
	var roleinfo RoleInfo
	err := json.Unmarshal(handler.data, &roleinfo)
	if err != nil {
		beego.Error(err)
		return
	}

	if roleinfo.GameId == 0 &&
		roleinfo.DistrictId == 0 &&
		roleinfo.ServerId == 0 &&
		roleinfo.RoleName == "" {
		beego.Error("role info error")
		return
	}

	c := handler.DB.C(models.UserInfoCollection)
	user := &models.UserInfo{}
	user.Account = roleinfo.Account
	si := models.UserServerInfo{}
	si.GameId = roleinfo.GameId
	si.DistrictId = roleinfo.DistrictId
	si.ServerId = roleinfo.ServerId
	si.RoleName = roleinfo.RoleName
	si.RoleLevel = roleinfo.Level

	if err := user.Read(c); err != nil {
		if err == mgo.ErrNotFound {
			user.ServerInfos = make([]models.UserServerInfo, 0, 8)
			user.ServerInfos = append(user.ServerInfos, si)
			if err = user.Insert(c); err != nil {
				beego.Error(err)
				return
			}
		} else {
			beego.Error(err)
			return
		}
	}

	if err := user.UpsetServerInfo(c, si); err != nil {
		beego.Error(err)
		return
	}

	beego.Info("update player:", user.Account, "role:", si.RoleName, "level:", si.RoleLevel)
}
