package controllers

type BaseUserRouter struct {
	BaseRouter
}

func (this *BaseUserRouter) NestPrepare() {
	if this.CheckActiveRedirect() {
		return
	}
}
