package setting

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
	"github.com/beego/compress"
	"github.com/beego/i18n"
)

var (
	AppHost             string
	AppUrl              string
	ActiveCodeLives     int
	ResetPwdCodeLives   int
	SecretKey           string
	EnforceRedirect     bool
	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	IsProMode           bool

	LoginRememberDays int
	LoginMaxRetries   int
	LoginFailedBlocks int

	CookieRememberName string
	CookieUserName     string

	StatusAddr string
	StatusPort int
	StatusDB   string
	Langs      []string
)

var (
	Cfg     *goconfig.ConfigFile
	Cache   cache.Cache
	Captcha *captcha.Captcha
)

var (
	AppConfPath      = "conf/app.conf"
	CompressConfPath = "conf/compress.json"
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

	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	Captcha = captcha.NewCaptcha("/captcha/", Cache)
	Captcha.FieldIDName = "CaptchaId"
	Captcha.FieldCaptchaName = "Captcha"

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = Cfg.MustValue("session", "session_provider", "file")
	//beego.BConfig.WebConfig.Session.SessionSavePath = Cfg.MustValue("session", "session_path", "sessions")
	beego.BConfig.WebConfig.Session.SessionName = Cfg.MustValue("session", "session_name", "cc_sess")
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8080")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8080/")
	EnforceRedirect = Cfg.MustBool("app", "enforce_redirect")

	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_minutes", 180)
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_minutes", 180)

	SecretKey = Cfg.MustValue("app", "secret_key")
	if len(SecretKey) == 0 {
		fmt.Println("Please set your secret_key in app.ini file")
	}

	IsProMode = beego.BConfig.RunMode == "pro"
	if IsProMode {
		beego.SetLevel(beego.LevelInformational)
	}

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "cc_magic")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name", "cc_powerful")

	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "sa:abc@tcp(192.168.1.180:3306)/nx_cc?charset=utf8&loc=UTC")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	StatusAddr = Cfg.MustValue("maple", "addr", "127.0.0.1")
	StatusPort = Cfg.MustInt("maple", "port", 5040)
	StatusDB = Cfg.MustValue("maple", "data_source", "127.0.0.1:27017?maxPoolSize=10")
	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		beego.Error(err)
	}
	orm.RunCommand()

	maxLifeTime := Cfg.MustInt("orm", "max_life_time", 3600)
	db, _ := orm.GetDB("default")
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTime))

	settingLocales()
	settingCompress()
}

func settingLocales() {
	// load locales with locale_LANG.ini files
	langs := "en-US|zh-CN"
	for _, lang := range strings.Split(langs, "|") {
		lang = strings.TrimSpace(lang)
		files := []string{"conf/" + "locale_" + lang + ".ini"}
		if fh, err := os.Open(files[0]); err == nil {
			fh.Close()
		} else {
			files = nil
		}
		if err := i18n.SetMessage(lang, "conf/global/"+"locale_"+lang+".ini", files...); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			os.Exit(2)
		}
	}
	Langs = i18n.ListLangs()
}

func settingCompress() {

	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, AppUrl)
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}
