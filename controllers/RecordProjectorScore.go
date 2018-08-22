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
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", uid)
	recordList := models.SearchProjectorSumPubInfoByOrder(filters...)
	if len(recordList) == 0 {
		c.ajaxMsg("该季度项目负责人评分未发布", MSG_ERR)
	}

	list := make([]map[string]interface{}, len(recordList))

	for i, v := range recordList {
		row := make(map[string]interface{})
		row["id"] = i + 1
		user, _ := models.SearchUserInfoByID(v.UserID)
		if user != nil {
			row["name"] = user.Name
		}

		project, _ := models.SearchProjectInfoByID(v.ProjectID)
		if project != nil {
			row["pname"] = project.Name
		}
		row["score"] = v.Score
		list[i] = row
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
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
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecords := logic.SearchProjectorScoreTypeRecordInfosBySumData(year, quarter, uid)
	list := make([]map[string]interface{}, len(typeRecords))
	for i, tr := range typeRecords {
		info := make(map[string]interface{})
		info["ID"] = tr.Type.ID
		info["Name"] = tr.Type.Name
		info["ScoreLimit"] = tr.Type.ScoreLimit
		info["Score"] = strconv.FormatFloat(tr.Score, 'f', 2, 64)
		list[i] = info
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}
