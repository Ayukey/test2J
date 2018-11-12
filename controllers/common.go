package controllers

import (
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/models"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	MSG_OK  = 0
	MSG_ERR = -1
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *models.Admin
	userID         int
	userName       string
	loginName      string
	pageSize       int
	allowURL       string
}

func (c *BaseController) Prepare() {
	c.pageSize = 20
	controllerName, actionName := c.GetControllerAndAction()
	c.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	c.actionName = strings.ToLower(actionName)
	c.Data["version"] = beego.AppConfig.String("version")
	c.Data["siteName"] = beego.AppConfig.String("site.name")
	c.Data["curRoute"] = c.controllerName + "." + c.actionName
	c.Data["curController"] = c.controllerName
	c.Data["curAction"] = c.actionName

	c.auth()

	c.Data["loginUserId"] = c.userID
	c.Data["loginUserName"] = c.userName
}

//Auth验证
func (c *BaseController) auth() {
	c.userID = 0
	session := c.StartSession()
	defer session.SessionRelease(c.Ctx.ResponseWriter)

	uid := session.Get("id")
	auth := session.Get("auth")

	fmt.Println("当前用户 id:", uid, " auth:", auth)

	if uid != nil && auth != nil {
		user, err := models.AdminGetByID(uid.(int))
		if err == nil && auth == libs.Md5([]byte(user.Password+user.Salt)) {
			c.userID = user.ID
			c.loginName = user.LoginName
			c.userName = user.RealName
			c.user = user
			//管理员验证
			c.AdminAuth()
		}
	}

	if c.userID == 0 && (c.controllerName != "login" && c.actionName != "login") {
		c.redirect(beego.URLFor("AuthController.Login"))
	}
}

// 获取侧边栏权限菜单
func (c *BaseController) AdminAuth() {
	// 左侧导航栏
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if c.userID != 1 {
		//普通管理员
		adminAuthIds, _ := models.RoleAuthGetByIds(c.user.RoleIds)
		adminAuthIdArr := strings.Split(adminAuthIds, ",")
		filters = append(filters, "id__in", adminAuthIdArr)
	}
	result, _ := models.AuthGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	list2 := make([]map[string]interface{}, len(result))
	allow_url := ""
	i, j := 0, 0
	for _, v := range result {
		if v.AuthURL != " " || v.AuthURL != "/" {
			allow_url += v.AuthURL
		}
		row := make(map[string]interface{})
		if v.Pid == 1 && v.IsShow == 1 {
			row["Id"] = int(v.ID)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = v.AuthURL
			row["Icon"] = v.Icon
			row["Pid"] = int(v.Pid)
			list[i] = row
			i++
		}
		if v.Pid != 1 && v.IsShow == 1 {
			row["Id"] = int(v.ID)
			row["Sort"] = v.Sort
			row["AuthName"] = v.AuthName
			row["AuthUrl"] = v.AuthURL
			row["Icon"] = v.Icon
			row["Pid"] = int(v.Pid)
			list2[j] = row
			j++
		}
	}

	c.Data["SideMenu1"] = list[:i]  //一级菜单
	c.Data["SideMenu2"] = list2[:j] //二级菜单

	c.allowURL = allow_url + "/home/index"
}

//是否POST提交
func (c *BaseController) isPost() bool {
	return c.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (c *BaseController) getClientIP() string {
	s := c.Ctx.Request.RemoteAddr
	l := strings.LastIndex(s, ":")
	return s[0:l]
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

//加载模板
func (c *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		tplname = c.controllerName + "/" + c.actionName + ".html"
	}
	fmt.Println("====加载模板====")
	fmt.Println(tplname)
	c.Layout = "public/layout.html"
	c.TplName = tplname
}

//相关ajax

//ajax返回
func (c *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["message"] = msg
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

//ajax返回 列表
func (self *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	self.Data["json"] = out
	self.ServeJSON()
	self.StopRun()
}

// 业务相关

//获取单个职位信息
func getPostionRoleInfo(infoList []*models.PositionRoleInfo, id int) (info *models.PositionRoleInfo) {
	for _, v := range infoList {
		if v.ID == id {
			return v
		}
	}
	return nil
}

//获取单个部门信息
func getDepartmentInfo(infoList []*models.DepartmentInfo, id int) (info *models.DepartmentInfo) {
	for _, v := range infoList {
		if v.ID == id {
			return v
		}
	}
	return nil
}

//获取单个项目信息
func getProjectInfo(infoList []*models.ProjectInfo, id int) (info *models.ProjectInfo) {
	for _, v := range infoList {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func getCurrentYearAndQuarter() string {
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()
	currentQuarter := 0

	if currentMonth >= 1 && currentMonth <= 3 {
		currentQuarter = 1
	} else if currentMonth >= 4 && currentMonth <= 6 {
		currentQuarter = 2
	} else if currentMonth >= 7 && currentMonth <= 9 {
		currentQuarter = 3
	} else if currentMonth >= 10 && currentMonth <= 12 {
		currentQuarter = 4
	}

	return fmt.Sprintf("%d-%d", currentYear, currentQuarter)
}
