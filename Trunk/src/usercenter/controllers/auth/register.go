package auth

import (
	"usercenter/controllers"
	"usercenter/modules/auth"
	"usercenter/modules/models"
	"usercenter/setting"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type RegisterController struct {
	controllers.BaseRouter
}

func (c *RegisterController) Get() {
	c.TplName = "register.html"
	captcha := setting.Captcha.CreateCaptchaHTML()
	c.Data["captcha"] = captcha
}

func (c *RegisterController) Register() {
	c.TplName = "register.html"
	captcha := setting.Captcha.CreateCaptchaHTML()
	c.Data["captcha"] = captcha

	var form auth.CreateForm
	if c.ValidFormSets(&form) == false {
		c.Data["errormsg"] = "验证表单失败"
		return
	}

	if !setting.Captcha.VerifyReq(c.Ctx.Request) {
		c.Data["errormsg"] = "验证码错误"
		return
	}

	var user models.User
	user.Lang = i18n.IndexLang(c.Lang)
	if err := auth.RegisterUser(&user, form.UserName, form.UserName, form.Email, form.PassWord); err != nil {
		c.Data["errormsg"] = err.Error()
		return
	}

	auth.SendRegisterMail(c.Locale, &user)

	loginRedirect := c.LoginUser(&user, false)
	if loginRedirect == "/" {
		c.FlashRedirect("/settings/profile", 302, "RegSuccess")
	} else {
		c.Redirect(loginRedirect, 302)
	}

	return

	c.Redirect("/", 302)
}

func (this *RegisterController) Active() {
	this.TplName = "active.html"

	// no need active
	if this.CheckActiveRedirect(false) {
		return
	}

	code := this.GetString(":code")

	var user models.User

	if auth.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = models.GetUserSalt()
		if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
			beego.Error("Active: user Update ", err)
		}
		if this.IsLogin {
			this.User = user
		}

		this.Redirect("/active/success", 302)

	} else {
		this.Data["Success"] = false
	}
}

// ActiveSuccess implemented success page when email active code verified.
func (this *RegisterController) ActiveSuccess() {
	this.TplName = "active.html"

	this.Data["Success"] = true
}
