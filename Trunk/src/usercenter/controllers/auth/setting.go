package auth

import (
	"usercenter/controllers"

	"github.com/beego/i18n"
)

type SettingsRouter struct {
	controllers.BaseRouter
}

// Profile implemented user profile settings page.
func (this *SettingsRouter) Profile() {
	this.TplName = "setting/profile.html"

	// need login
	if this.CheckLoginRedirect() {
		return
	}
}

func (this *SettingsRouter) ChangeLang() {
	if this.IsLogin {
		lanidx := i18n.IndexLang(this.Lang)
		this.User.Lang = lanidx
		this.User.Update("lang")
	}

	this.Redirect(this.Ctx.Request.Referer(), 302)
}
