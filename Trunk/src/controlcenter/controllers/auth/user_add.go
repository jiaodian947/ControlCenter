package auth

import (
	"controlcenter/controllers"
	"controlcenter/modules/auth"
	"controlcenter/modules/models"
)

type UserAddController struct {
	controllers.BaseAdmin
}

func (c *UserAddController) Get() {
	c.Data["title"] = "人员管理"
	c.Data["active"] = "index"
	c.TplName = "user_add.html"
}

func (c *UserAddController) Post() {
	c.Data["title"] = "人员管理"
	c.Data["active"] = "index"
	c.TplName = "user_add.html"

	var form auth.CreateForm
	if c.ValidFormSets(&form) == false {
		return
	}

	var user models.User
	if err := auth.RegisterUser(&user, form.UserName, form.NickName, form.Email, form.PassWord); err != nil {
		return
	}
	c.Ctx.Redirect(302, "/")
}
