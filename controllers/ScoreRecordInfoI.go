package controllers

import (
	"jg2j_server/models"
)

// 项目评分顶级模版
type ScoreRecordInfoIController struct {
	BaseController
}

func (c *ScoreRecordInfoIController) Score() {
	c.Data["pageTitle"] = "项目评分"
	projects := models.SearchAllProjectInfo()
	c.Data["projects"] = projects
	c.display()
}

func (c *ScoreRecordInfoIController) Search() {
	//列表
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	pid, _ := c.GetInt("pid", 0)

	typeRecordList := models.SerachScoreTypeRecordInfoIList(year, quarter, pid)
	list := make([]map[string]interface{}, len(typeRecordList))

	for k, v := range typeRecordList {
		row := make(map[string]interface{})
		row["id"] = v.Type.ID
		row["name"] = v.Type.Name

		if v.Record == nil {
			row["score"] = "暂无评分"
		} else {
			row["rid"] = v.Record.ID
			row["score"] = v.Record.TotalScore
		}

		list[k] = row
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}
