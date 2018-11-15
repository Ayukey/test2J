package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 项目负责人评分模版
type ProjectorScoreTypeInfoController struct {
	BaseController
}

func (c *ProjectorScoreTypeInfoController) List() {
	c.Data["pageTitle"] = "项目负责人评分模版"
	c.display()
}

func (c *ProjectorScoreTypeInfoController) Table() {

	result := models.SearchAllProjectLeaderTemplates()
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["score_limit"] = v.ScoreLimit
		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ProjectorScoreTypeInfoController) Edit() {
	c.Data["pageTitle"] = "编辑项目负责人评分模版"

	id, _ := c.GetInt("id", 0)
	t, err := models.SearchProjectLeaderTemplateByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = t.ID
	row["name"] = t.Name
	row["score_limit"] = t.ScoreLimit
	c.Data["template"] = row
	c.display()
}

//存储资源
func (c *ProjectorScoreTypeInfoController) AjaxSave() {
	tid, _ := c.GetInt("id")
	template, _ := models.SearchProjectLeaderTemplateByID(tid)
	// 修改
	template.Name = strings.TrimSpace(c.GetString("name"))
	template.ScoreLimit, _ = c.GetFloat("score_limit")

	if err := template.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "success")
}
