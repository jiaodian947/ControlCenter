package routers

import (
	"manage/controllers"
	"manage/controllers/auth"
	"manage/setting"

	"github.com/astaxie/beego"
)

func InitRouter() {
	beego.InsertFilter("/captcha/*", beego.BeforeRouter, setting.Captcha.Handler)
	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.MainController{})

	beego.Router("/login", &auth.LoginController{})
	beego.Router("/send", &controllers.Mt2SrvController{})
	beego.Router("/query", &controllers.DatabaseController{})
}
