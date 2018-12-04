package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

type CheckRight struct {
	BaseRouter
}

type Reply struct {
	Status int
	Data   interface{}
}

func (c *CheckRight) CheckRight() {

}

func (c *CheckRight) CheckLogin() {
	if !c.IsLogin {
		c.Abort("403")
	}
}

func (c *CheckRight) NestPrepare() {
	c.CheckLogin()
	c.CheckRight()
}

func init() {
	beego.ErrorHandler("403", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(403)
		rw.Write([]byte("Forbidden"))
	})
}
