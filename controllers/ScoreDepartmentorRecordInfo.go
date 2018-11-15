package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 后台查看、修改部门负责人互评得分
type ScoreDepartmentorRecordInfoController struct {
	BaseController
}

// 跳转部门负责人评分页
func (c *ScoreDepartmentorRecordInfoController) Score() {
	c.Data["pageTitle"] = "部门负责人互评得分"
	projectors := models.SearchAllProjectLeaders()
	departmentors := models.SearchAllDepartmentLeaders()
	c.Data["projectors"] = projectors
	c.Data["departmentors"] = departmentors
	c.display()
}

// 查询评分信息
func (c *ScoreDepartmentorRecordInfoController) Search() {
	uid, _ := c.GetInt("uid", 0)  // 部门负责人（被打分）
	suid, _ := c.GetInt("uid", 0) // 项目负责人（打分）
	departmentID, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecords := logic.SearchDepartmentLeaderScoreTemplateRecords(year, quarter, uid, suid, departmentID)

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

func (c *ScoreDepartmentorRecordInfoController) Edit() {
	tid, _ := c.GetInt("tid", 0)
	uid, _ := c.GetInt("uid", 0)   // 部门负责人（被打分）
	suid, _ := c.GetInt("suid", 0) // 项目负责人（打分）
	departmentID, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateReocrd := logic.SearchDepartmentLeaderScoreTemplateRecordByTID(year, quarter, uid, suid, tid, departmentID)

	if templateReocrd == nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	row["uid"] = uid
	row["suid"] = suid
	row["tid"] = templateReocrd.Template.ID
	row["name"] = templateReocrd.Template.Name
	row["maxscore"] = templateReocrd.Template.ScoreLimit
	if templateReocrd.Record != nil {
		row["score"] = templateReocrd.Record.Score
	} else {
		row["score"] = "暂无评分"
	}

	c.Data["Source"] = row
	c.Data["pageTitle"] = "部门负责人评分详情"
	c.display()
}

//存储资源
func (c *ScoreDepartmentorRecordInfoController) AjaxSave() {
	tid, _ := c.GetInt("tid", 0)
	uid, _ := c.GetInt("uid", 0)   // 部门负责人（被打分）
	suid, _ := c.GetInt("suid", 0) // 项目负责人（打分）
	departmentID, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	score, _ := c.GetFloat("score", 0)

	templateRecord := logic.SearchDepartmentLeaderScoreTemplateRecordByTID(year, quarter, uid, suid, tid, departmentID)

	// 修改细分记录
	if templateRecord.Record == nil {
		record := new(models.DepartmentLeaderScoreRecord)
		record.UID = uid
		record.SUID = suid
		record.Score = score
		record.TID = tid
		record.Year = year
		record.Quarter = quarter
		record.DepartmentID = departmentID

		if err := models.AddDepartmentLeaderScoreRecord(record); err != nil {
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
	err := logic.SaveDepartmentLeaderSumScoreRecord(year, quarter, uid, suid, departmentID)
	if err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}

	// 如果已发布，更新发布信息
	filter1 := models.DBFilter{Key: "year", Value: year}                    // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}              // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}                      // 用户ID
	filter4 := models.DBFilter{Key: "departmentor_id", Value: departmentID} // 部门ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchDepartmentLeaderReleaseRecordsByFilters(filters...)
	if len(records) != 0 {
		err := logic.ReleaseDepartmentLeaderScoreRecord(year, quarter, uid, departmentID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "")
}
