package controllers

import (
	"jg2j_server/models"
)

// 项目评分顶级模版
type ScoreRecordInfoIIController struct {
	BaseController
}

func (c *ScoreRecordInfoIIController) Score() {
	tid, _ := c.GetInt("tid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	scoreTypeI, _ := models.SearchScoreTypeInfoIByID(tid)
	project, _ := models.SearchProjectInfoByID(pid)
	row := make(map[string]interface{})
	row["tid"] = tid
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = scoreTypeI.Name + " (" + project.Name + ")"
	c.display()
}

func (c *ScoreRecordInfoIIController) Search() {
	tid, _ := c.GetInt("tid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordList := models.SerachScoreTypeRecordInfoIIList(year, quarter, tid, pid)
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
