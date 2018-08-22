package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Data["pageTitle"] = "系统首页"
	c.TplName = "public/main.html"
}

func (c *HomeController) Start() {
	c.Data["pageTitle"] = "控制面板"
	c.display()
}
