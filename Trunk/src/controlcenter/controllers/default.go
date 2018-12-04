package controllers

type MainController struct {
	BaseAdmin
}

func (c *MainController) Get() {
	c.TplName = "index.html"
	c.Data["active"] = "index"
	c.Data["title"] = "控制台"
}
