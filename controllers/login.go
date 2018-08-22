package controllers

import (
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type LoginController struct {
	BaseController
}

//登录 TODO:XSRF过滤
func (c *LoginController) LoginIn() {
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
				fmt.Println("当前用户:")
				fmt.Println(user)
				authkey := libs.Md5([]byte(c.getClientIP() + "|" + user.Password + user.Salt))
				// 保存cookie
				c.Ctx.SetCookie("auth", strconv.Itoa(user.ID)+"|"+authkey, 7*86400)
				c.redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&c.Controller)
			c.redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}
	c.TplName = "login/login.html"
}
