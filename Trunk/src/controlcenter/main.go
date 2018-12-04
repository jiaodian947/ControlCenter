package main

import (
	"controlcenter/controllers/maple"
	"controlcenter/modules/auth"
	"controlcenter/modules/models"
	"controlcenter/routers"
	"os"

	"controlcenter/setting"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func CreateDefaultAdmin() error {
	var user models.User
	user.NickName = "Admin"
	user.IsActive = true
	user.IsAdmin = true
	user.IsForbid = false
	if err := auth.RegisterUser(&user, "admin", "admin", "admin@sininm.com", "admin"); err != nil {
		return err
	}

	return nil
}

func initialize() {
	setting.LoadConfig()
	RunCommand()
	routers.InitRouter()
}

func RunCommand() {
	if len(os.Args) < 2 || os.Args[1] != "syncdb" {
		return
	}

	if err := orm.RunSyncdb("default", false, false); err != nil {
		beego.Error(err)
		os.Exit(1)
	}

	if err := CreateDefaultAdmin(); err != nil {
		beego.Error(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func main() {
	initialize()
	conn := maple.StartUdpServer(setting.StatusAddr, setting.StatusPort)
	defer conn.Close()
	beego.Run()
	beego.Info("quit server")
}
