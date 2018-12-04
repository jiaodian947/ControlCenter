package auth

import (
	"manage/controllers"
	"manage/models"
	"manage/modules/auth"
	"manage/setting"
)

type LoginController struct {
	controllers.BaseRouter
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
	captcha := setting.Captcha.CreateCaptchaHTML()

	c.Data["captcha"] = captcha
	if c.GetString("quit") == "true" {
		auth.LogoutUser(c.Ctx)
		c.Redirect("/login", 302)
		return
	}

	if c.CheckLoginRedirect(false) {
		c.Redirect("/", 302)
		return
	}
}

func (c *LoginController) Post() {
	captcha := setting.Captcha.CreateCaptchaHTML()
	c.Data["captcha"] = captcha
	c.TplName = "login.html"

	if c.CheckLoginRedirect(false) {
		c.Redirect("/", 302)
		return
	}

	var form auth.LoginForm
	// valid form and put errors to template context

	if !setting.Captcha.VerifyReq(c.Ctx.Request) {
		return
	}

	var user models.User
	if auth.VerifyUser(&user, form.UserName, form.Password) {
		c.LoginUser(&user, form.Remember)
		c.Redirect("/", 302)
		return
	}

}
