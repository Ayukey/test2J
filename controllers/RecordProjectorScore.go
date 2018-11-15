package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
)

// 部门信息
type RecordProjectorScoreController struct {
	BaseController
}

func (c *RecordProjectorScoreController) Score() {
	c.Data["pageTitle"] = "项目负责人互评记录"
	leaders := models.SearchAllProjectLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *RecordProjectorScoreController) Search() {
	uid, _ := c.GetInt("uid", 0)
	projectId, _ := c.GetInt("projectId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}            // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}      // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}              // 用户ID
	filter4 := models.DBFilter{Key: "project_id", Value: projectId} // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchProjectLeaderReleaseRecordsByOrder(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度项目负责人评分未发布")
	}

	list := make([]map[string]interface{}, len(records))

	for i, r := range records {
		row := make(map[string]interface{})
		row["id"] = i + 1
		user, err := models.SearchUserByID(r.UID)
		if err == nil {
			row["name"] = user.Name
		}

		project, err := models.SearchProjectByID(r.ProjectID)
		if err == nil {
			row["pname"] = project.Name
		}
		row["score"] = r.Score
		list[i] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *RecordProjectorScoreController) ScoreDetails() {
	uid, _ := c.GetInt("uid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	row := make(map[string]interface{})
	row["uid"] = uid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = "项目负责人互评记录详情"
	c.display()
}

func (c *RecordProjectorScoreController) SearchDetails() {
	uid, _ := c.GetInt("uid", 0)
	projectId, _ := c.GetInt("projectId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecords := logic.SearchProjectLeaderTemplateAverageRecords(year, quarter, uid, projectId)
	list := make([]map[string]interface{}, len(typeRecords))
	for i, tr := range typeRecords {
		info := make(map[string]interface{})
		info["ID"] = tr.Template.ID
		info["Name"] = tr.Template.Name
		info["ScoreLimit"] = tr.Template.ScoreLimit
		info["Score"] = strconv.FormatFloat(tr.Score, 'f', 2, 64)
		list[i] = info
	}

	c.ajaxList(MSG_OK, "成功", list)
}
