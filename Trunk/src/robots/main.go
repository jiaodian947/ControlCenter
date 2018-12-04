package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"robots/controller"
	"robots/game"
	"robots/ui"
	"robots/utils"
	"runtime"
	"syscall"
	"time"
)

var (
	pprof = flag.Bool("pprof", false, "pprof")
)

func StartPprof() {
	//runtime.MemProfileRate = 1
	if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
		log.Println("statr pprof failed", err)
		return
	}
	log.Println("start pprof at :6060")
}

func RobotJsonConf() {
	JsonParse := utils.NewJsonStruct()
	JsonParse.Load("./conf.json", utils.JsonConf)
}

func main() {

	RobotJsonConf()
	flag.Parse()
	if *pprof {
		go StartPprof()
	}
	rand.Seed(time.Now().Unix())
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	controller.New(func(acc, pwd, name string, index int) controller.Robot {
		return game.NewGameClient(acc, pwd, name, index)
	})
	go ui.Serv(utils.JsonConf.StartPort)
	openbrowser(fmt.Sprintf("http://127.0.0.1:%d", utils.JsonConf.StartPort))
	<-exitChan
	controller.Shutdown()
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println(err)
	}

}
