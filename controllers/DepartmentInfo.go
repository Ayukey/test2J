package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 部门信息
type DepartmentInfoController struct {
	BaseController
}

func (c *DepartmentInfoController) List() {
	c.Data["pageTitle"] = "部门列表"
	c.display()
}

func (c *DepartmentInfoController) Table() {
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

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, count := models.SearchDepartmentInfoList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		u, err := models.SearchUserInfoByID(v.UserID)
		if err == nil {
			row["user_name"] = u.Name
		}
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *DepartmentInfoController) Add() {
	c.Data["pageTitle"] = "新增部门"
	leaders := models.SearchAllDepartmentLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *DepartmentInfoController) Edit() {
	c.Data["pageTitle"] = "编辑部门"

	id, _ := c.GetInt("id", 0)
	d, err := models.SearchDepartmentInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	leaders := models.SearchAllDepartmentLeaders()

	row := make(map[string]interface{})
	row["id"] = d.ID
	row["name"] = d.Name
	row["user_id"] = d.UserID
	c.Data["Source"] = row
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *DepartmentInfoController) AjaxSave() {
	id, _ := c.GetInt("id")

	if id == 0 {
		department := new(models.DepartmentInfo)
		department.Name = strings.TrimSpace(c.GetString("name"))
		department.UserID, _ = c.GetInt("user_id", 0)
		department.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchDepartmentInfoByName(department.Name)
		if err == nil {
			c.ajaxMsg("该部门已经存在", MSG_ERR)
		}

		id, err := models.AddDepartmentInfo(department)
		if err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}

		if department.UserID > 0 {
			user, err := models.SearchUserInfoByID(department.UserID)
			if err == nil {
				user.DepartmentID = int(id)
				_ = user.Update()
			}
		}
		c.ajaxMsg("", MSG_OK)
	}

	department, _ := models.SearchDepartmentInfoByID(id)
	// 修改
	department.Name = strings.TrimSpace(c.GetString("name"))
	department.UserID, _ = c.GetInt("user_id", 0)

	if err := department.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}

	if department.UserID > 0 {
		user, err := models.SearchUserInfoByID(id)
		if err == nil {
			user.DepartmentID = int(id)
			_ = user.Update()
		}
	}
	c.ajaxMsg("", MSG_OK)
}

func (c *DepartmentInfoController) Members() {
	c.Data["pageTitle"] = "部门成员列表"
	id, _ := c.GetInt("id", 0)
	row := make(map[string]interface{})
	row["id"] = id
	c.Data["Source"] = row
	c.display()
}

func (c *DepartmentInfoController) MembersTable() {
	id, _ := c.GetInt("id", 0)

	//列表
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 30
	}

	c.pageSize = limit

	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result := models.SearchAllDepartmentMembers(id)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["account"] = v.Account
		row["name"] = v.Name
		row["phone"] = v.Phone
		role, err := models.SearchPositionRoleInfoByID(v.RoleID)
		if err == nil {
			row["role"] = role.Name
		}
		list[k] = row
	}
	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}
