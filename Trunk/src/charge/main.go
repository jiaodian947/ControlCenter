package main

import (
	"charge/server"
	"charge/setting"
	"os"
	"os/signal"
	"syscall"

	_ "charge/models"
	_ "charge/transaction"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initialize() {
	setting.LoadConfig()
	orm.RunCommand()
	orm.RunSyncdb("default", false, true)
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

	srv := server.New()

	srv.Main()
	<-exitChan
	srv.Exit()
}
