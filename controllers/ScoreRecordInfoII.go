package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 项目评分顶级模版
type ScoreRecordInfoIIController struct {
	BaseController
}

func (c *ScoreRecordInfoIIController) Score() {
	template1ID, _ := c.GetInt("template1_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	template1, _ := models.SearchProjectTemplate1ByID(template1ID)
	project, _ := models.SearchProjectByID(projectID)
	row := make(map[string]interface{})
	row["template1_id"] = template1ID
	row["project_id"] = projectID
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = template1.Name + " (" + project.Name + ")"
	c.display()
}

func (c *ScoreRecordInfoIIController) Search() {
	template1ID, _ := c.GetInt("template1_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template2Records := logic.SearchProjectTemplate2Records(year, quarter, template1ID, projectID)
	list := make([]map[string]interface{}, len(template2Records))

	for index, template2Record := range template2Records {
		row := make(map[string]interface{})
		row["template2_id"] = template2Record.Template.ID
		row["template2_name"] = template2Record.Template.Name

		if template2Record.Record == nil {
			row["record2_score"] = "暂无评分"
		} else {
			row["record2_id"] = template2Record.Record.ID
			row["record2_score"] = libs.Float64ToStringWithNoZero(template2Record.Record.TotalScore)
		}

		list[index] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}
