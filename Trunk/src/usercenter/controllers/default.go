package controllers

type MainController struct {
	BaseRouter
}

func (c *MainController) Get() {
	c.TplName = "index.html"
}
