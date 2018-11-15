package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
)

// 部门信息
type RecordDepartmentorScoreController struct {
	BaseController
}

func (c *RecordDepartmentorScoreController) Score() {
	c.Data["pageTitle"] = "部门负责人互评记录"
	leaders := models.SearchAllDepartmentLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

func (c *RecordDepartmentorScoreController) Search() {
	uid, _ := c.GetInt("uid", 0)
	departmentId, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}                 // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}           // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}                   // 用户ID
	filter4 := models.DBFilter{Key: "departmen_id", Value: departmentId} // 用户ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchDepartmentLeaderReleaseRecordsByOrder(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度部门负责人评分未发布")
	}

	list := make([]map[string]interface{}, len(records))

	for i, r := range records {
		row := make(map[string]interface{})
		row["id"] = i + 1

		user, err := models.SearchUserByID(r.UID)
		if err == nil {
			row["name"] = user.Name
		}
		department, err := models.SearchDepartmentByID(r.DepartmentID)
		if err == nil {
			row["dname"] = department.Name
		}
		row["score"] = r.Score
		list[i] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *RecordDepartmentorScoreController) ScoreDetails() {
	uid, _ := c.GetInt("uid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	row := make(map[string]interface{})
	row["uid"] = uid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = "部门负责人互评记录详情"
	c.display()
}

func (c *RecordDepartmentorScoreController) SearchDetails() {
	uid, _ := c.GetInt("uid", 0)
	departmentId, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecords := logic.SearchDepartmentLeaderTemplateAverageRecords(year, quarter, uid, departmentId)
	list := make([]map[string]interface{}, len(templateRecords))
	for i, tr := range templateRecords {
		info := make(map[string]interface{})
		info["ID"] = tr.Template.ID
		info["Name"] = tr.Template.Name
		info["ScoreLimit"] = tr.Template.ScoreLimit
		info["Score"] = strconv.FormatFloat(tr.Score, 'f', 2, 64)
		list[i] = info
	}

	c.ajaxList(MSG_OK, "成功", list)
}
