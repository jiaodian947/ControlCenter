package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"usercenter/modules/models"
	"usercenter/modules/utils"
	"usercenter/routers"
	"usercenter/setting"

	"github.com/astaxie/beego"

	_ "github.com/go-sql-driver/mysql"
)

func initialize() {
	setting.LoadConfig()
	routers.InitRouter()
}

func RegisterUser(id int, username, nickname, email, password string) error {
	var user models.User
	user.Id = id
	user.Lang = 1

	// use random salt encode password
	salt := models.GetUserSalt()
	pwd := utils.EncodePassword(password, salt)

	user.UserName = strings.ToLower(username)
	user.Email = strings.ToLower(email)

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	// Use username as default nickname.
	user.NickName = nickname

	user.IsActive = true
	user.IsAdmin = false

	return user.Insert()
}

func runCommand() {
	if len(os.Args) < 2 || os.Args[1] != "ct" {
		return
	}

	count := 100
	if len(os.Args) >= 3 {
		c, err := strconv.Atoi(os.Args[2])
		if err == nil {
			count = c
		}
	}

	for i := 1; i < count; i++ {
		name := fmt.Sprintf("test%d", i)
		err := RegisterUser(0, name, name, fmt.Sprintf("%s@sininm.com", name), "123456")
		fmt.Println(name, err)
	}

	os.Exit(1)
}

func main() {
	initialize()
	runCommand()
	if !setting.IsProMode {
		beego.SetStaticPath("/static_source", "static_source")
		beego.BConfig.WebConfig.DirectoryIndex = true
	}

	beego.Run()
}
