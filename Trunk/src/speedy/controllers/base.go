package controllers

import (
	"fmt"
	"html/template"
	"net/url"
	"speedy/models"
	"speedy/modules/auth"
	"time"

	"github.com/astaxie/beego"
)

type NestPreparer interface {
	NestPrepare()
}

type BaseRouter struct {
	beego.Controller
	User    models.User
	IsLogin bool
}

func (this *BaseRouter) Prepare() {
	// page start time
	this.Data["PageStartTime"] = time.Now()

	// start session
	this.StartSession()

	switch {
	// save logined user if exist in session
	case auth.GetUserFromSession(&this.User, this.CruSession):
		this.IsLogin = true
	// save logined user if exist in remember cookie
	case auth.LoginUserFromRememberCookie(&this.User, this.Ctx):
		this.IsLogin = true
	}

	// read flash message
	beego.ReadFromRequest(&this.Controller)

	// pass xsrf helper to template context
	xsrfToken := this.Controller.XSRFToken()
	this.Data["xsrf_token"] = xsrfToken
	this.Data["xsrf_html"] = template.HTML(this.Controller.XSRFFormHTML())

	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

func (this *BaseRouter) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !this.IsLogin {
		if len(redirect_to) == 0 {
			req := this.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/#/login?to=" + url.QueryEscape(redirect_to)
		this.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && this.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		this.Redirect(redirect_to, code)
		return true
	}
	return false
}

// read beego flash message
func (this *BaseRouter) FlashRead(key string) (string, bool) {
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (this *BaseRouter) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

// check xsrf and show a friendly page
func (this *BaseRouter) CheckXsrfCookie() bool {
	return this.Controller.CheckXSRFCookie()
}
