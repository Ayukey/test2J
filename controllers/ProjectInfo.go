package controllers

import (
	"jg2j_server/models"
	"strings"
	"time"
)

// 项目信息
type ProjectInfoController struct {
	BaseController
}

func (c *ProjectInfoController) List() {
	c.Data["pageTitle"] = "项目列表"
	c.display()
}

func (c *ProjectInfoController) Table() {
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
	result, count := models.SearchProjectInfoList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		startStr := time.Unix(int64(v.BeginDate/1000), 0).Format("2006-01-02")
		endStr := time.Unix(int64(v.EndDate/1000), 0).Format("2006-01-02")
		if v.BeginDate > 0 {
			row["start_time"] = startStr
		}

		if v.EndDate > 0 {
			row["end_time"] = endStr
		}
		u, err := models.SearchUserInfoByID(v.UserID)
		if err == nil {
			row["userName"] = u.Name
		}
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ProjectInfoController) Add() {
	c.Data["pageTitle"] = "新增项目"
	leaders := models.SearchAllProjectLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *ProjectInfoController) Edit() {
	c.Data["pageTitle"] = "编辑项目"

	id, _ := c.GetInt("id", 0)
	p, err := models.SearchProjectInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	leaders := models.SearchAllProjectLeaders()

	row := make(map[string]interface{})
	row["id"] = p.ID
	row["name"] = p.Name
	row["user_id"] = p.UserID

	startStr := time.Unix(int64(p.BeginDate/1000), 0).Format("2006-01-02")
	endStr := time.Unix(int64(p.EndDate/1000), 0).Format("2006-01-02")
	if p.BeginDate > 0 {
		row["start_time"] = startStr
	}

	if p.EndDate > 0 {
		row["end_time"] = endStr
	}

	c.Data["Source"] = row
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *ProjectInfoController) AjaxSave() {
	id, _ := c.GetInt("id")

	if id == 0 {
		project := new(models.ProjectInfo)
		project.Name = strings.TrimSpace(c.GetString("name"))
		project.UserID, _ = c.GetInt("user_id", 0)
		project.BeginDate, _ = c.GetFloat("start_time", 0)
		project.EndDate, _ = c.GetFloat("end_time", 0)
		project.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchProjectInfoByName(project.Name)

		if err == nil {
			c.ajaxMsg("该项目已经存在", MSG_ERR)
		}

		if _, err := models.AddProjectInfo(project); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
		c.ajaxMsg("", MSG_OK)
	}

	project, _ := models.SearchProjectInfoByID(id)
	// 修改
	project.Name = strings.TrimSpace(c.GetString("name"))
	project.UserID, _ = c.GetInt("user_id", 0)
	project.BeginDate, _ = c.GetFloat("start_time", 0)
	project.EndDate, _ = c.GetFloat("end_time", 0)

	if err := project.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}

func (c *ProjectInfoController) Members() {
	c.Data["pageTitle"] = "项目成员列表"
	id, _ := c.GetInt("id", 0)
	row := make(map[string]interface{})
	row["id"] = id
	c.Data["Source"] = row
	c.display()
}

func (c *ProjectInfoController) MembersTable() {
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
	result := models.SearchAllProjectMembers(id)
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
