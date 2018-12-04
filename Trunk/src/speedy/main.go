package main

import (
	_ "speedy/models"
	_ "speedy/routers"
	"speedy/server"
	"speedy/setting"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	setting.LoadConfig()
	orm.RunCommand()
	server.Run()
	beego.Run()
}
