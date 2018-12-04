package auth

import (
	"usercenter/controllers"
	"usercenter/modules/auth"
	"usercenter/modules/models"
	"usercenter/setting"
)

type LoginController struct {
	controllers.BaseRouter
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
	c.Data["title"] = "Login"
	c.Data["errname"] = ""
	c.Data["errpass"] = ""
	captcha := setting.Captcha.CreateCaptchaHTML()

	c.Data["captcha"] = captcha

	if c.CheckLoginRedirect(false) {
		c.Redirect("/", 302)
		return
	}
}

func (c *LoginController) Logout() {
	lang := c.Lang
	auth.LogoutUser(c.Ctx)
	c.Redirect("/login?lang="+lang, 302)
}

func (c *LoginController) Post() {

	c.Data["title"] = "Login"
	c.Data["errname"] = c.Tr("login.name_or_pass_err")
	c.Data["errpass"] = c.Tr("login.name_or_pass_err")
	captcha := setting.Captcha.CreateCaptchaHTML()
	c.Data["captcha"] = captcha
	c.TplName = "login.html"

	if c.CheckLoginRedirect(false) {
		c.Redirect("/", 302)
		return
	}

	var form auth.LoginForm
	// valid form and put errors to template context
	if c.ValidFormSets(&form) == false {
		return
	}

	if !setting.Captcha.VerifyReq(c.Ctx.Request) {
		c.Data["errname"] = ""
		c.Data["errpass"] = ""
		c.Data["errcaptcha"] = c.Tr("login.captcha_err")
		return
	}
	//setting.Captcha.VerifyReq(c.Ctx.Request)
	var user models.User
	if auth.VerifyUser(&user, form.UserName, form.Password) {
		c.LoginUser(&user, form.Remember)
		c.Redirect("/", 302)
		return
	}

}
