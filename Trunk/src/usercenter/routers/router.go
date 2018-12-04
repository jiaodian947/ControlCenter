package routers

import (
	"usercenter/controllers"
	"usercenter/controllers/api"
	"usercenter/controllers/auth"
	"usercenter/setting"

	"github.com/astaxie/beego"
)

func InitRouter() {
	beego.ErrorController(&controllers.ErrorController{})
	beego.Router("/", &controllers.MainController{})
	beego.InsertFilter("/captcha/*", beego.BeforeRouter, setting.Captcha.Handler) // 验证码接口

	login := new(auth.LoginController) // 网站登录登出接口
	beego.Router("/login", login, "get:Get;post:Post")
	beego.Router("/logout", login, "get:Logout")

	register := new(auth.RegisterController) // 网站激活接口
	beego.Router("/register", register, "get:Get;post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	game := new(api.GameActiveController) // 游戏内激活接口
	beego.Router("/game/active", game, "post:Active")
	beego.Router("/game/generate/:count", game, "get:Generate") // 生成激活码

	ga := new(api.GameAuthController) // 游戏内的帐号验证
	beego.Router("/auth", ga, "get:Get;post:Post")

	gr := new(api.GameReg) // 游戏内帐号注册
	beego.Router("/game/register", gr, "post:Register")

	s := new(auth.SettingsRouter)
	beego.Router("/settings/profile", s, "get:Profile")
	beego.Router("/language", s, "get:ChangeLang")

	p := new(api.ProfileRouter)
	beego.Router("/profile/verify", p, "get:Verify")
}
