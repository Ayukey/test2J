package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 项目评分二级模版
type ScoreTypeInfoIIController struct {
	BaseController
}

func (c *ScoreTypeInfoIIController) List() {
	tid, _ := c.GetInt("tid", 0)
	scoreTypeI, _ := models.SearchScoreTypeInfoIByID(tid)
	row := make(map[string]interface{})
	row["tid"] = tid
	c.Data["Source"] = row
	c.Data["pageTitle"] = "项目二级评分模版 (" + scoreTypeI.Name + ")"
	c.display()
}

func (c *ScoreTypeInfoIIController) Table() {
	tid, _ := c.GetInt("tid", 0)
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
	filters = append(filters, "tid", tid)
	result, count := models.SearchScoreTypeInfoIIList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["percentage"] = v.Percentage
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ScoreTypeInfoIIController) Add() {
	tid, _ := c.GetInt("tid", 0)
	scoreTypeI, _ := models.SearchScoreTypeInfoIByID(tid)
	row := make(map[string]interface{})
	row["tid"] = tid
	row["t_name"] = scoreTypeI.Name
	c.Data["Source"] = row
	c.Data["pageTitle"] = "新增二级模版"
	c.display()
}

func (c *ScoreTypeInfoIIController) Edit() {
	typeID, _ := c.GetInt("id", 0)
	scoreType, err := models.SearchScoreTypeInfoIIByID(typeID)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = scoreType.ID
	row["name"] = scoreType.Name
	row["percentage"] = scoreType.Percentage
	c.Data["Source"] = row
	c.Data["pageTitle"] = "编辑项目评分模版"
	c.display()
}

//存储资源
func (c *ScoreTypeInfoIIController) AjaxSave() {
	typeID, _ := c.GetInt("id")

	if typeID == 0 {
		scoreType := new(models.ScoreTypeInfoII)
		tid, _ := c.GetInt("tid")
		scoreType.TID = tid
		scoreType.Name = strings.TrimSpace(c.GetString("name"))
		scoreType.Percentage, _ = c.GetFloat("percentage")
		scoreType.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchScoreTypeInfoIIByName(scoreType.Name)

		if err == nil {
			c.ajaxMsg("该模版名称已经存在", MSG_ERR)
		}

		if _, err := models.AddScoreTypeInfoII(scoreType); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
		c.ajaxMsg("", MSG_OK)
	}

	scoreType, _ := models.SearchScoreTypeInfoIIByID(typeID)
	scoreType.Name = strings.TrimSpace(c.GetString("name"))
	percentage, _ := c.GetFloat("percentage")
	scoreType.Percentage = percentage

	if err := scoreType.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
