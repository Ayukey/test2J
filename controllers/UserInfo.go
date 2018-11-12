package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/models"
	"strings"
)

// 用户信息
type UserInfoController struct {
	BaseController
}

// 跳转用户信息模块
func (c *UserInfoController) List() {
	c.Data["pageTitle"] = "用户列表"
	postions := models.SearchAllPositionRoleInfo()
	c.Data["postions"] = postions
	c.display()
}

// 查询用户信息数据
func (c *UserInfoController) Table() {
	postion, _ := c.GetInt("postion", 0)

	page, err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 30
	}
	c.pageSize = limit

	positionRoleInfos := models.SearchAllPositionRoleInfo()
	departmentInfos := models.SearchAllDepartmentInfo()
	projectInfos := models.SearchAllProjectInfo()

	//查询条件

	var result []*models.UserInfo
	var count int64 = 0

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if postion != 0 {
		filters = append(filters, "role_id", postion)
		result = models.SearchAllUserInfoList(filters...)
		count = int64(len(result))
	} else {
		result, count = models.SearchUserInfoList(page, c.pageSize, filters...)
	}

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["account"] = v.Account
		row["name"] = v.Name
		row["phone"] = v.Phone
		if getPostionRoleInfo(positionRoleInfos, v.RoleID) != nil {
			row["role"] = getPostionRoleInfo(positionRoleInfos, v.RoleID).Name
		} else {
			row["role"] = ""
		}

		if getDepartmentInfo(departmentInfos, v.DepartmentID) != nil {
			row["department"] = getDepartmentInfo(departmentInfos, v.DepartmentID).Name
		} else {
			row["department"] = ""
		}

		if getProjectInfo(projectInfos, v.ProjectID) != nil {
			row["project"] = getProjectInfo(projectInfos, v.ProjectID).Name
		} else {
			row["project"] = ""
		}

		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

// 跳转用户信息新增模块
func (c *UserInfoController) Add() {
	c.Data["pageTitle"] = "新增用户"

	positionRoleInfos := models.SearchAllPositionRoleInfo()
	departmentInfos := models.SearchAllDepartmentInfo()
	projectInfos := models.SearchAllProjectInfo()

	c.Data["postionList"] = positionRoleInfos
	c.Data["departmentList"] = departmentInfos
	c.Data["projectList"] = projectInfos

	c.display()
}

// 跳转用户信息编辑模块
func (c *UserInfoController) Edit() {
	c.Data["pageTitle"] = "编辑用户"

	id, _ := c.GetInt("id", 0)
	u, err := models.SearchUserInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	positionRoleInfos := models.SearchAllPositionRoleInfo()
	departmentInfos := models.SearchAllDepartmentInfo()
	projectInfos := models.SearchAllProjectInfo()

	row := make(map[string]interface{})
	row["id"] = u.ID
	row["account"] = u.Account
	row["name"] = u.Name
	row["phone"] = u.Phone
	row["role_id"] = u.RoleID
	row["department_id"] = u.DepartmentID
	row["project_id"] = u.ProjectID
	c.Data["Source"] = row

	c.Data["postionList"] = positionRoleInfos
	c.Data["departmentList"] = departmentInfos
	c.Data["projectList"] = projectInfos
	c.display()
}

// 跳转用户信息修改密码模块
func (c *UserInfoController) Change() {
	c.Data["pageTitle"] = "修改密码"

	id, _ := c.GetInt("id", 0)
	u, err := models.SearchUserInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["id"] = u.ID
	row["account"] = u.Account
	row["name"] = u.Name
	c.Data["Source"] = row
	c.display()
}

// 保存用户信息
func (c *UserInfoController) AjaxSave() {
	id, _ := c.GetInt("id")

	if id == 0 {
		// 修改用户信息逻辑
		password := strings.TrimSpace(c.GetString("password"))
		rpassword := strings.TrimSpace(c.GetString("r_password"))

		if password != rpassword {
			c.ajaxMsg("两遍密码不一致", MSG_ERR)
		}

		user := new(models.UserInfo)
		user.Account = strings.TrimSpace(c.GetString("account"))
		user.Name = strings.TrimSpace(c.GetString("name"))
		user.Password = libs.Md5([]byte(password))
		user.Phone = strings.TrimSpace(c.GetString("phone"))
		user.RoleID, _ = c.GetInt("position_id", 0)
		user.DepartmentID, _ = c.GetInt("department_id", 0)
		user.ProjectID, _ = c.GetInt("project_id", 0)
		user.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchUserInfoByAccount(user.Account)
		if err == nil {
			c.ajaxMsg("该员工已经存在", MSG_ERR)
		}

		if _, err := models.AddUserInfo(user); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}

		if user.RoleID == 5 {
			projectorInfo := new(models.ProjectorInfo)
			projectorInfo.UserID = user.ID
			projectorInfo.BeginDate = 1528330000000
			projectorInfo.EndDate = 1559350000000
			models.AddPProjectorInfo(projectorInfo)
		}

		c.ajaxMsg("", MSG_OK)
	}

	// 修改用户信息逻辑
	user, _ := models.SearchUserInfoByID(id)
	user.Name = strings.TrimSpace(c.GetString("name"))
	user.Phone = strings.TrimSpace(c.GetString("phone"))
	user.DepartmentID, _ = c.GetInt("department_id", 0)
	user.RoleID, _ = c.GetInt("position_id", 0)
	user.ProjectID, _ = c.GetInt("project_id", 0)

	if err := user.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}

// 保存用户密码信息
func (c *UserInfoController) AjaxSavePassword() {
	password := strings.TrimSpace(c.GetString("password"))
	rpassword := strings.TrimSpace(c.GetString("r_password"))

	if password != rpassword {
		c.ajaxMsg("两遍密码不一致", MSG_ERR)
	}

	id, _ := c.GetInt("id")
	user, _ := models.SearchUserInfoByID(id)
	user.Password = libs.Md5([]byte(password))

	if err := user.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}

// 删除用户
func (c *UserInfoController) DeleteUser() {
	id, err := c.GetInt("id")
	if err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}

	user, err := models.SearchUserInfoByID(id)
	if err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	user.Status = 0

	if err = user.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
