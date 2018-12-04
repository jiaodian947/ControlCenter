package routers

import (
	"controlcenter/controllers"
	"controlcenter/controllers/api"
	"controlcenter/controllers/auth"
	"controlcenter/controllers/log"
	"controlcenter/controllers/maple"
	"controlcenter/setting"

	"github.com/astaxie/beego"
)

func InitRouter() {
	beego.InsertFilter("/captcha/*", beego.BeforeRouter, setting.Captcha.Handler)
	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.MainController{})
	beego.Router("/login", &auth.LoginController{})
	beego.Router("/user/add", &auth.UserAddController{})
	maple.AddMapleRouter()

	//api
	mapapi := new(api.MapleAPI)
	beego.Router("/api/maple/:gameid:int/:seckey/:account/:filter:int", mapapi, "get:PrePull")
	beego.Router("/api/maple/:gameid:int/:seckey/:account", mapapi, "get:PrePull")
	beego.Router("/api/maple/:gameid:int/:seckey", mapapi, "get:PrePull")
	beego.Router("/api/maple", mapapi, "post:PrePull")

	//clientlog
	l := new(log.ClientLog)
	beego.Router("/log/upload", l, "get:GetUpload;post:Upload")
	beego.Router("/log", l, "get:GetLogs")
}
