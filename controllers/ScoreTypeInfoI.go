package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 项目评分顶级模版
type ScoreTypeInfoIController struct {
	BaseController
}

func (c *ScoreTypeInfoIController) List() {
	c.Data["pageTitle"] = "项目一级评分模版"
	c.display()
}

func (c *ScoreTypeInfoIController) Table() {
	//查询条件
	result := models.SearchAllProjectTemplate1s()
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ScoreTypeInfoIController) Add() {
	c.Data["pageTitle"] = "新增一级评分模版"
	c.display()
}

func (c *ScoreTypeInfoIController) Edit() {
	c.Data["pageTitle"] = "编辑一级评分模版"

	typeID, _ := c.GetInt("id", 0)
	scoreType, err := models.SearchProjectTemplate1ByID(typeID)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = scoreType.ID
	row["name"] = scoreType.Name
	c.Data["template"] = row
	c.display()
}

//存储资源
func (c *ScoreTypeInfoIController) AjaxSave() {
	typeID, _ := c.GetInt("id")

	if typeID == 0 {
		scoreType := new(models.ProjectTemplate1)
		scoreType.Name = strings.TrimSpace(c.GetString("name"))
		scoreType.Status = 1

		if err := models.AddProjectTemplate1(scoreType); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
		c.ajaxMsg(MSG_OK, "success")
	}

	scoreType, _ := models.SearchProjectTemplate1ByID(typeID)
	scoreType.Name = strings.TrimSpace(c.GetString("name"))

	if err := scoreType.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "success")
}
