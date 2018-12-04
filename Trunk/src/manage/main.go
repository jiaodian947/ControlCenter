package main

import (
	"manage/routers"
	"manage/server"
	"manage/setting"
	"os"
	"os/signal"
	"syscall"

	"github.com/astaxie/beego"

	_ "manage/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initialize() {
	setting.LoadConfig()
	orm.RunCommand()
	//orm.RunSyncdb("default", false, true)
	routers.InitRouter()
}

func main() {
	initialize()

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	server.Run()
	go beego.Run()
	<-exitChan
	server.Exit()
}
