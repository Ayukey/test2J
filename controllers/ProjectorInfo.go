package controllers

import (
	"jg2j_server/models"
	"time"
)

// 项目负责人信息
type ProjectorInfoController struct {
	BaseController
}

func (c *ProjectorInfoController) List() {
	c.Data["pageTitle"] = "项目负责人信息列表"
	c.display()
}

func (c *ProjectorInfoController) Table() {
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
	result, count := models.SearchProjectorInfoList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		user, _ := models.SearchUserInfoByID(v.UserID)
		project, _ := models.SearchProjectInfoByID(user.ProjectID)
		row := make(map[string]interface{})
		row["id"] = v.ID
		if user != nil {
			row["account"] = user.Account
			row["userName"] = user.Name
		}
		if project != nil {
			row["projectName"] = project.Name
		}
		startStr := time.Unix(int64(v.BeginDate/1000), 0).Format("2006-01-02")
		if v.BeginDate > 0 {
			row["start_time"] = startStr
		}
		endStr := time.Unix(int64(v.EndDate/1000), 0).Format("2006-01-02")
		if v.EndDate > 0 {
			row["end_time"] = endStr
		}
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ProjectorInfoController) Edit() {
	c.Data["pageTitle"] = "编辑项目负责人"

	id, _ := c.GetInt("id", 0)
	p, err := models.SearchProjectorInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	user, _ := models.SearchUserInfoByID(p.UserID)
	project, _ := models.SearchProjectInfoByID(user.ProjectID)

	row := make(map[string]interface{})
	row["id"] = p.ID
	if user != nil {
		row["account"] = user.Account
		row["userName"] = user.Name
	}
	if project != nil {
		row["projectName"] = project.Name
	}

	startStr := time.Unix(int64(p.BeginDate/1000), 0).Format("2006-01-02")
	endStr := time.Unix(int64(p.EndDate/1000), 0).Format("2006-01-02")
	if p.BeginDate > 0 {
		row["start_time"] = startStr
	}

	if p.EndDate > 0 {
		row["end_time"] = endStr
	}

	c.Data["Source"] = row
	c.display()
}

//存储资源
func (c *ProjectorInfoController) AjaxSave() {
	id, _ := c.GetInt("id")

	projector, _ := models.SearchProjectorInfoByID(id)
	projector.BeginDate, _ = c.GetFloat("start_time", 0)
	projector.EndDate, _ = c.GetFloat("end_time", 0)

	if err := projector.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
