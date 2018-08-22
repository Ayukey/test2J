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
	result, count := models.SearchScoreTypeInfoIList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ScoreTypeInfoIController) Add() {
	c.Data["pageTitle"] = "新增一级评分模版"
	c.display()
}

func (c *ScoreTypeInfoIController) Edit() {
	c.Data["pageTitle"] = "编辑一级评分模版"

	typeID, _ := c.GetInt("id", 0)
	scoreType, err := models.SearchScoreTypeInfoIByID(typeID)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = scoreType.ID
	row["name"] = scoreType.Name
	c.Data["Source"] = row
	c.display()
}

//存储资源
func (c *ScoreTypeInfoIController) AjaxSave() {
	typeID, _ := c.GetInt("id")

	if typeID == 0 {
		scoreType := new(models.ScoreTypeInfoI)
		scoreType.Name = strings.TrimSpace(c.GetString("name"))
		scoreType.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchScoreTypeInfoIByName(scoreType.Name)

		if err == nil {
			c.ajaxMsg("该模版名称已经存在", MSG_ERR)
		}

		if _, err := models.AddScoreTypeInfoI(scoreType); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
		c.ajaxMsg("", MSG_OK)
	}

	scoreType, _ := models.SearchScoreTypeInfoIByID(typeID)
	scoreType.Name = strings.TrimSpace(c.GetString("name"))

	if err := scoreType.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
