package controllers

type BaseAdmin struct {
	BaseRouter
}

func (this *BaseAdmin) NestPrepare() {
	if this.CheckLoginRedirect() {
		return
	}
}
