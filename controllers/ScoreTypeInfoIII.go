package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 项目评分三级模版
type ScoreTypeInfoIIIController struct {
	BaseController
}

func (c *ScoreTypeInfoIIIController) List() {
	tid, _ := c.GetInt("tid", 0)
	scoreTypeII, _ := models.SearchScoreTypeInfoIIByID(tid)
	row := make(map[string]interface{})
	row["tid"] = tid
	c.Data["Source"] = row
	c.Data["pageTitle"] = "项目三级评分模版 (" + scoreTypeII.Name + ")"
	c.display()
}

func (c *ScoreTypeInfoIIIController) Table() {
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
	result, count := models.SearchScoreTypeInfoIIIList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["score"] = v.MaxScore
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ScoreTypeInfoIIIController) Add() {
	tid, _ := c.GetInt("tid", 0)
	scoreTypeII, _ := models.SearchScoreTypeInfoIIByID(tid)
	row := make(map[string]interface{})
	row["tid"] = tid
	row["t_name"] = scoreTypeII.Name
	c.Data["Source"] = row
	c.Data["pageTitle"] = "新增三级模版"
	c.display()
}

func (c *ScoreTypeInfoIIIController) Edit() {
	typeID, _ := c.GetInt("id", 0)
	scoreType, err := models.SearchScoreTypeInfoIIIByID(typeID)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = scoreType.ID
	row["name"] = scoreType.Name
	row["score"] = scoreType.MaxScore
	c.Data["Source"] = row
	c.Data["pageTitle"] = "编辑项目三级评分模版"
	c.display()
}

//存储资源
func (c *ScoreTypeInfoIIIController) AjaxSave() {
	typeID, _ := c.GetInt("id")

	if typeID == 0 {
		scoreType := new(models.ScoreTypeInfoIII)
		tid, _ := c.GetInt("tid")
		scoreType.TID = tid
		scoreType.Name = strings.TrimSpace(c.GetString("name"))
		scoreType.MaxScore, _ = c.GetFloat("score")
		scoreType.Status = 1

		// 检查登录名是否已经存在
		_, err := models.SearchScoreTypeInfoIIIByName(scoreType.Name)

		if err == nil {
			c.ajaxMsg("该模版名称已经存在", MSG_ERR)
		}

		if _, err := models.AddScoreTypeInfoIII(scoreType); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
		c.ajaxMsg("", MSG_OK)
	}

	scoreType, _ := models.SearchScoreTypeInfoIIIByID(typeID)
	// 修改
	scoreType.Name = strings.TrimSpace(c.GetString("name"))
	score, _ := c.GetFloat("score")
	scoreType.MaxScore = score

	if err := scoreType.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
