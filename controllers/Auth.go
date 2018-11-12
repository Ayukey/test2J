package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/models"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type AuthController struct {
	BaseController
}

// 登录
func (c *AuthController) Login() {
	session := c.StartSession()
	defer session.SessionRelease(c.Ctx.ResponseWriter)

	if c.userID > 0 {
		c.redirect(beego.URLFor("HomeController.Index"))
	}

	beego.ReadFromRequest(&c.Controller)
	if c.isPost() {
		username := strings.TrimSpace(c.GetString("username"))
		password := strings.TrimSpace(c.GetString("password"))
		if username != "" && password != "" {
			user, err := models.AdminGetByName(username)
			flash := beego.NewFlash()
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "账号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该账号已禁用"
			} else {
				user.LastIP = c.getClientIP()
				user.LastLogin = time.Now().Unix()
				user.Update()
				auth := libs.Md5([]byte(user.Password + user.Salt))

				// 保存session
				session.Set("id", user.ID)
				session.Set("auth", auth)

				c.redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&c.Controller)
			c.redirect(beego.URLFor("AuthController.Login"))
		}
	}

	c.TplName = "login/login.html"
}

// 登出
func (c *AuthController) Logout() {
	session := c.StartSession()
	defer session.SessionRelease(c.Ctx.ResponseWriter)
	session.Delete("id")
	session.Delete("auth")
	c.redirect(beego.URLFor("AuthController.Login"))
}
