package setting

import (
	"fmt"
	"os"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
)

var (
	AppConfPath = "conf/app.conf"
)

var (
	SecretKey string
	Cache     cache.Cache
	Captcha   *captcha.Captcha

	AppHost            string
	LoginRememberDays  int
	LoginMaxRetries    int
	LoginFailedBlocks  int
	CookieRememberName string
	CookieUserName     string
	IdleTimeout        int = 600
)

func LoadConfig() {

	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.

	Cfg, err := goconfig.LoadConfigFile(AppConfPath)
	if err != nil {
		fmt.Println("Fail to load configuration file: " + err.Error())
		os.Exit(2)
	}

	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	Captcha = captcha.NewCaptcha("/captcha/", Cache)
	Captcha.FieldIDName = "CaptchaId"
	Captcha.FieldCaptchaName = "Captcha"

	AppHost = Cfg.MustValue("app", "app_host", "localhost")
	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "cc_magic")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name", "cc_powerful")

	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "sa:abc@tcp(192.168.1.180:3306)/nx_cc?charset=utf8&loc=UTC")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)
	maxLifeTime := Cfg.MustInt("orm", "max_life_time", 3600)

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, _ := orm.GetDB("default")
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTime))
}
