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
	result := models.SearchAllDepartments()

	list := make([]map[string]interface{}, len(result))
	for k, d := range result {
		row := make(map[string]interface{})
		row["id"] = d.ID
		row["name"] = d.Name
		u, err := models.SearchUserByID(d.Leader)
		if err == nil {
			row["leaderName"] = u.Name
		}
		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *DepartmentInfoController) Add() {
	c.Data["pageTitle"] = "添加部门"
	leaders := models.SearchAllDepartmentLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *DepartmentInfoController) Edit() {
	c.Data["pageTitle"] = "编辑部门"

	id, _ := c.GetInt("id", 0)

	d, err := models.SearchDepartmentByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	leaders := models.SearchAllDepartmentLeaders()

	row := make(map[string]interface{})
	row["id"] = d.ID
	row["name"] = d.Name
	row["leader"] = d.Leader
	c.Data["department"] = row
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *DepartmentInfoController) AjaxSave() {
	id, _ := c.GetInt("id")
	name := strings.TrimSpace(c.GetString("name"))
	leader, _ := c.GetInt("leader", 0)

	if id == 0 {
		department := new(models.Department)
		department.Name = name
		department.Leader = leader
		department.Status = 1

		if err := models.AddDepartment(department); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}

		c.ajaxMsg(MSG_OK, "success")
	}

	// 修改
	department, _ := models.SearchDepartmentByID(id)
	department.Name = name
	department.Leader = leader

	if err := department.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}

	c.ajaxMsg(MSG_OK, "success")
}
