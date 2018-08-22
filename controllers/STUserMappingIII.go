package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/models"
	"strings"
)

// 用户信息
type STUserMappingIIIController struct {
	BaseController
}

func (c *STUserMappingIIIController) List() {
	c.Data["pageTitle"] = "用户列表"
	c.display()
}

func (c *STUserMappingIIIController) Table() {
	//列表
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
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, count := models.SearchUserInfoList(page, c.pageSize, filters...)
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

func (c *STUserMappingIIIController) Add() {
	c.Data["pageTitle"] = "新增用户"

	positionRoleInfos := models.SearchAllPositionRoleInfo()
	departmentInfos := models.SearchAllDepartmentInfo()
	projectInfos := models.SearchAllProjectInfo()

	c.Data["postionList"] = positionRoleInfos
	c.Data["departmentList"] = departmentInfos
	c.Data["projectList"] = projectInfos

	c.display()
}

func (c *STUserMappingIIIController) Edit() {
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

func (c *STUserMappingIIIController) Change() {
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

//存储资源
func (c *STUserMappingIIIController) AjaxSave() {
	id, _ := c.GetInt("id")

	if id == 0 {
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
		c.ajaxMsg("", MSG_OK)
	}

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

//存储资源
func (c *STUserMappingIIIController) AjaxSavePassword() {
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
