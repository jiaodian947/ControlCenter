package setting

import (
	"fmt"
	"os"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/orm"
)

type Platform struct {
	Name     string
	Path     string
	TestPath string
}

var (
	AppHost        string
	AppPort        int
	HeartTimeout   int
	GameId         int
	WorkThreads    int
	QueueLen       int
	OrderRedirect  string
	Platforms      map[string]Platform
	PerChannelLen  int
	RetryMax       int
	VerifyInterval float64
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
	RetryMax = Cfg.MustInt("", "retry_max", 0)
	VerifyInterval = Cfg.MustFloat64("", "verify_interval", 60)
	HeartTimeout = Cfg.MustInt("tcp", "heart_timeout", 300)
	PerChannelLen = Cfg.MustInt("", "channel_len", 32)
	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	count := Cfg.MustInt("platform", "count", 0)
	Platforms = make(map[string]Platform)
	for i := 0; i < count; i++ {
		sec := fmt.Sprintf("platform%d", i)
		pt := Platform{}
		pt.Name = Cfg.MustValue(sec, "name", "")
		pt.Path = Cfg.MustValue(sec, "path", "")
		pt.TestPath = Cfg.MustValue(sec, "test", "")
		if pt.Name == "" {
			panic("platform name is nil")
		}
		Platforms[pt.Name] = pt
	}

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	maxLifeTime := Cfg.MustInt("orm", "max_life_time", 3600)
	db, _ := orm.GetDB("default")
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTime))
}
