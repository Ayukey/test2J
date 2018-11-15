package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 后台查看、修改部门负责人互评得分
type ScoreProjectorRecordInfoController struct {
	BaseController
}

// 跳转部门负责人评分页
func (c *ScoreProjectorRecordInfoController) Score() {
	c.Data["pageTitle"] = "部门负责人互评得分"
	projectors := models.SearchAllProjectLeaders()
	departmentors := models.SearchAllDepartmentLeaders()
	c.Data["projectors"] = projectors
	c.Data["departmentors"] = departmentors
	c.display()
}

// 查询评分信息
func (c *ScoreProjectorRecordInfoController) Search() {
	uid, _ := c.GetInt("uid", 0)  // 项目负责人（被打分）
	suid, _ := c.GetInt("uid", 0) // 部门负责人（打分）
	projectID, _ := c.GetInt("projectId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecords := logic.SearchProjectLeaderTemplateRecords(year, quarter, uid, suid, projectID)

	list := make([]map[string]interface{}, len(templateRecords))

	for i, tr := range templateRecords {
		row := make(map[string]interface{})
		row["id"] = tr.Template.ID
		row["name"] = tr.Template.Name
		row["max_score"] = tr.Template.ScoreLimit
		if tr.Record == nil {
			row["score"] = "暂无评分"
		} else {
			row["score"] = tr.Record.Score
		}
		list[i] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ScoreProjectorRecordInfoController) Edit() {
	tid, _ := c.GetInt("tid", 0)
	uid, _ := c.GetInt("uid", 0)  // 项目负责人（被打分）
	suid, _ := c.GetInt("uid", 0) // 部门负责人（打分）
	projectID, _ := c.GetInt("projectId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecord := logic.SearchProjectLeaderTemplateRecordByTID(year, quarter, uid, suid, tid, projectID)

	if templateRecord == nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	row["uid"] = uid
	row["suid"] = suid
	row["tid"] = templateRecord.Template.ID
	row["name"] = templateRecord.Template.Name
	row["maxscore"] = templateRecord.Template.ScoreLimit
	if templateRecord.Record != nil {
		row["score"] = templateRecord.Record.Score
	} else {
		row["score"] = "暂无评分"
	}

	c.Data["Source"] = row
	c.Data["pageTitle"] = "部门负责人评分详情"
	c.display()
}

//存储资源
func (c *ScoreProjectorRecordInfoController) AjaxSave() {
	tid, _ := c.GetInt("tid", 0)
	uid, _ := c.GetInt("uid", 0)  // 项目负责人（被打分）
	suid, _ := c.GetInt("uid", 0) // 部门负责人（打分）
	projectID, _ := c.GetInt("projectId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	score, _ := c.GetFloat("score", 0)

	templateRecord := logic.SearchProjectLeaderTemplateRecordByTID(year, quarter, uid, suid, tid, projectID)

	if templateRecord.Record == nil {
		record := new(models.ProjectLeaderScoreRecord)
		record.UID = uid
		record.SUID = suid
		record.Score = score
		record.TID = tid
		record.Year = year
		record.Quarter = quarter
		record.ProjectID = projectID

		if err := models.AddProjectLeaderScoreRecord(record); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		record := templateRecord.Record
		record.Score = score
		if err := record.Update(); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	// 更新总记录
	err := logic.SaveProjectLeaderSumScoreRecord(year, quarter, uid, suid, projectID)
	if err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}

	// 如果已发布，更新发布信息
	filter1 := models.DBFilter{Key: "year", Value: year}            // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}      // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}              // 用户ID
	filter4 := models.DBFilter{Key: "project_id", Value: projectID} // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchProjectLeaderReleaseRecordsByFilters(filters...)
	if len(records) != 0 {
		err := logic.ReleaseProjectLeaderScoreRecord(year, quarter, uid, projectID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "success")
}
