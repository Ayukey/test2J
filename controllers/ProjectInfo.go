package controllers

import (
	"jg2j_server/models"
	"strings"
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
	result := models.SearchAllProjects()

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		u, err := models.SearchUserByID(v.Leader)
		if err == nil {
			row["leaderName"] = u.Name
		}
		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ProjectInfoController) Add() {
	c.Data["pageTitle"] = "添加项目"
	leaders := models.SearchAllProjectLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *ProjectInfoController) Edit() {
	c.Data["pageTitle"] = "编辑项目"

	id, _ := c.GetInt("id", 0)
	p, err := models.SearchProjectByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	leaders := models.SearchAllProjectLeaders()

	row := make(map[string]interface{})
	row["id"] = p.ID
	row["name"] = p.Name
	row["leader"] = p.Leader

	c.Data["project"] = row
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *ProjectInfoController) AjaxSave() {
	id, _ := c.GetInt("id")
	name := strings.TrimSpace(c.GetString("name"))
	leader, _ := c.GetInt("leader", 0)

	if id == 0 {
		project := new(models.Project)
		project.Name = name
		project.Active = 1
		project.Leader = leader
		project.LeaderActive = 1
		project.Status = 1

		if err := models.AddProject(project); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
		c.ajaxMsg(MSG_OK, "success")
	}

	// 修改
	project, _ := models.SearchProjectByID(id)
	project.Name = name
	project.Leader = leader

	if err := project.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "success")
}
