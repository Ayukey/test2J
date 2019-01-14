package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 项目评分顶级模版
type ScoreRecordInfoIController struct {
	BaseController
}

func (c *ScoreRecordInfoIController) Score() {
	c.Data["pageTitle"] = "项目评分"
	projects := models.SearchAllProjects()
	c.Data["projects"] = projects
	c.display()
}

func (c *ScoreRecordInfoIController) Search() {
	//列表
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	projectID, _ := c.GetInt("project_id", 0)

	template1Records := logic.SearchProjectTemplate1Records(year, quarter, projectID)
	list := make([]map[string]interface{}, len(template1Records))

	for index, template1Record := range template1Records {
		row := make(map[string]interface{})
		row["template1_id"] = template1Record.Template.ID
		row["template1_name"] = template1Record.Template.Name

		if template1Record.Record == nil {
			row["record1_score"] = "暂无评分"
		} else {
			row["record1_id"] = template1Record.Record.ID
			row["record1_score"] = libs.Float64ToStringWithNoZero(template1Record.Record.TotalScore)
		}

		list[index] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}
