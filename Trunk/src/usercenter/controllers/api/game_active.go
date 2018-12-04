package api

import (
	"usercenter/controllers"
	"usercenter/modules/auth"
	"usercenter/modules/models"
	"usercenter/modules/utils"
)

type GameActiveController struct {
	controllers.BaseRouter
}

type ActiveResult struct {
	ErrCode int
}

type CDKey struct {
	Key []string
}

func (c *GameActiveController) Generate() {
	count, _ := c.GetInt(":count", 1)
	ret := new(CDKey)
	ret.Key = make([]string, 0, count)
	for true {
		uuid := utils.UUID()
		ac := new(models.ActivationCode)
		ac.Code = uuid
		if ac.Insert() != nil {
			continue
		}
		ret.Key = append(ret.Key, uuid)
		if len(ret.Key) >= count {
			break
		}
	}

	c.Data["json"] = ret
	c.ServeJSON()
}

func (c *GameActiveController) Active() {
	username := c.GetString("user")
	cdkey := c.GetString("cdkey")
	res := ActiveResult{}
	user := new(models.User)
	if !auth.HasUser(user, username) {
		res.ErrCode = CAS_ERR_NO_ACCOUNT
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if user.IsActive {
		res.ErrCode = CAS_ERR_ACCOUNT_ACTIVED
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if !auth.VerifyUserWithCDKey(user, cdkey) {
		res.ErrCode = CAS_ERR_CDKEY_INVALID
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	user.IsActive = true
	user.Rands = models.GetUserSalt()
	if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
		res.ErrCode = CAS_ERR_ACTIVE_FAILED
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.ErrCode = ERR_SUCCESS
	c.Data["json"] = res
	c.ServeJSON()
}
