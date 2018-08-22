package controllers

import (
	"jg2j_server/models"
	"strings"
)

// 部门负责人评分模版
type DepartmentorScoreTypeInfoController struct {
	BaseController
}

func (c *DepartmentorScoreTypeInfoController) List() {
	c.Data["pageTitle"] = "部门负责人评分模版"
	c.display()
}

func (c *DepartmentorScoreTypeInfoController) Table() {
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
	result, count := models.SearchDepartmentorScoreTypeInfoList(page, c.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["score_limit"] = v.ScoreLimit
		list[k] = row
	}
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *DepartmentorScoreTypeInfoController) Edit() {
	c.Data["pageTitle"] = "编辑部门负责人评分模版"

	id, _ := c.GetInt("id", 0)
	t, err := models.SearchDepartmentorScoreTypeInfoByID(id)
	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}
	row := make(map[string]interface{})
	row["id"] = t.ID
	row["name"] = t.Name
	row["score_limit"] = t.ScoreLimit
	c.Data["Source"] = row
	c.display()
}

//存储资源
func (c *DepartmentorScoreTypeInfoController) AjaxSave() {
	tid, _ := c.GetInt("id")
	template, _ := models.SearchDepartmentorScoreTypeInfoByID(tid)
	// 修改
	template.Name = strings.TrimSpace(c.GetString("name"))
	template.ScoreLimit, _ = c.GetFloat("score_limit")

	if err := template.Update(); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("", MSG_OK)
}
