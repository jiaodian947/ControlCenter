package setting

import (
	"fmt"
	"os"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/orm"
)

type ChannelInfo struct {
	ChannelType   int
	ChannelName   string
	AuthUrl       string
	EncodeFunc    string
	AccountPrefix string
}

var (
	AppHost       string
	AppPort       int
	HeartTimeout  int
	GameId        int
	WorkThreads   int
	QueueLen      int
	Channels      []ChannelInfo
	OrderRedirect string
)

var (
	Cfg *goconfig.ConfigFile
)

var (
	AppConfPath = "conf/conf.ini"
)

func LoadConfig() {

	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.

	Cfg, err = goconfig.LoadConfigFile(AppConfPath)
	if err != nil {
		fmt.Println("Fail to load configuration file: " + err.Error())
		os.Exit(2)
	}

	Cfg.BlockMode = false

	AppHost = Cfg.MustValue("", "app_host", "127.0.0.1")
	AppPort = Cfg.MustInt("", "app_port", 0)
	GameId = Cfg.MustInt("", "game_id", 0)
	WorkThreads = Cfg.MustInt("", "work_threads", 1)
	QueueLen = Cfg.MustInt("", "queue_len", 16)
	HeartTimeout = Cfg.MustInt("tcp", "heart_timeout", 300)
	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "sa:abc@tcp(192.168.1.180:3306)/nx_cc?charset=utf8&loc=UTC")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)
	maxLifeTime := Cfg.MustInt("orm", "max_life_time", 3600)

	channels := Cfg.MustInt("channel", "channels", 1)
	OrderRedirect = Cfg.MustValue("channel", "order_redirect_url", "/order/notify")
	Channels = make([]ChannelInfo, 0, channels)
	for i := 0; i < channels; i++ {
		sec := fmt.Sprintf("channel_%d", i)
		ch := ChannelInfo{}
		ch.ChannelType = Cfg.MustInt(sec, "channel_type", 0)
		ch.AuthUrl = Cfg.MustValue(sec, "auth_url", "")
		ch.ChannelName = Cfg.MustValue(sec, "channel_name", "")
		ch.EncodeFunc = Cfg.MustValue(sec, "encode_func", "")
		ch.AccountPrefix = Cfg.MustValue(sec, "account_prefix", "")
		Channels = append(Channels, ch)
	}

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, _ := orm.GetDB("default")
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTime))
}
