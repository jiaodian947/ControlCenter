package auth

import (
	"log"
	"speedy/models"
	"speedy/setting"
	"speedy/utils"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"

	"github.com/astaxie/beego"
)

func LoginClientWhiteList(ctx *context.Context) (success bool) {

	ip_list, err := models.GetAllWhitelist(nil, []string{"ClientIp"}, nil, nil, 0, 0)
	if err != nil {
		success = false
		return
	}
	log.Println(ip_list)
	client_host := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	for i := 0; i < len(ip_list); i++ {
		// log.Println(ip_list[i].(map[string]interface{}))
		// log.Println(client_host)
		for _, value := range ip_list[i].(map[string]interface{}) {
			// strValue := fmt.Sprintf("%v", value)
			if client_host == value {
				success = true
				return
			}
		}
	}
	success = false
	return
}

func LoginUser(ctx *context.Context, user *models.User, remember bool) {

	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
	ctx.Input.CruSession = beego.GlobalSessions.SessionRegenerateID(ctx.ResponseWriter, ctx.Request)
	ctx.Input.CruSession.Set("auth_user_id", user.Id)
	ctx.Input.CruSession.Set("auth_user", user)

	if remember {
		WriteRememberCookie(user, ctx)
	}
}

func LoginUserFromRememberCookie(user *models.User, ctx *context.Context) (success bool) {
	userName := ctx.GetCookie(setting.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	defer func() {
		if !success {
			DeleteRememberCookie(ctx)
		}
	}()

	var err error
	user, err = models.GetUserByName(userName)
	if err != nil {
		success = false
		return
	}

	secret := utils.EncodeMd5(user.Password)
	value, _ := ctx.GetSecureCookie(secret, setting.CookieRememberName)
	if value != userName {
		return false
	}

	LoginUser(ctx, user, true)

	return true
}

// logout user
func LogoutUser(ctx *context.Context) {
	DeleteRememberCookie(ctx)
	ctx.Input.CruSession.Delete("auth_user_id")
	ctx.Input.CruSession.Delete("auth_user")
	ctx.Input.CruSession.Flush()
	beego.GlobalSessions.SessionDestroy(ctx.ResponseWriter, ctx.Request)
}

func WriteRememberCookie(user *models.User, ctx *context.Context) {
	secret := utils.EncodeMd5(user.Password)
	days := 86400 * setting.LoginRememberDays
	ctx.SetCookie(setting.CookieUserName, user.Name, days)
	ctx.SetSecureCookie(secret, setting.CookieRememberName, user.Name, days)
}

func DeleteRememberCookie(ctx *context.Context) {
	ctx.SetCookie(setting.CookieUserName, "", -1)
	ctx.SetCookie(setting.CookieRememberName, "", -1)
}

func GetUserIdFromSession(sess session.Store) int {
	if id, ok := sess.Get("auth_user_id").(int); ok && id > 0 {
		return id
	}
	return 0
}

func GetUserFromSession(user *models.User, sess session.Store) bool {
	au := sess.Get("auth_user")

	if au == nil {
		return false
	}

	if u, ok := au.(*models.User); ok && u != nil {
		*user = *u
		return true
	}

	return false
}
