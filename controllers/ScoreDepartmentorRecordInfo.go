package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
	"strings"
)

type ScoreDepartmentorRecordInfoController struct {
	BaseController
}

// 跳转部门负责人评分页
func (c *ScoreDepartmentorRecordInfoController) Score() {
	c.Data["pageTitle"] = "对部门负责人评分"
	projectLeaders := models.SearchAllProjectLeadersInProject()
	departmentLeaders := models.SearchAllDepartmentLeadersInDepartment()
	c.Data["projectLeaders"] = projectLeaders
	c.Data["departmentLeaders"] = departmentLeaders
	c.display()
}

// 查询评分信息
func (c *ScoreDepartmentorRecordInfoController) Search() {
	departmentLeader := strings.Split(c.GetString("departmentLeader"), "|") // 部门负责人（被打分）
	departmentLeaderID := 0
	departmentID := 0

	if len(departmentLeader) == 2 {
		departmentLeaderID, _ = strconv.Atoi(departmentLeader[0])
		departmentID, _ = strconv.Atoi(departmentLeader[1])
	}

	projectLeader := strings.Split(c.GetString("projectLeader"), "|") // 项目负责人（打分）
	projectLeaderID := 0
	if len(projectLeader) == 2 {
		projectLeaderID, _ = strconv.Atoi(projectLeader[0])
	}

	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecords := logic.SearchDepartmentLeaderScoreTemplateRecords(year, quarter, departmentLeaderID, projectLeaderID, departmentID)

	list := make([]map[string]interface{}, len(templateRecords))

	for i, tr := range templateRecords {
		row := make(map[string]interface{})
		row["projectLeader_id"] = projectLeaderID
		row["departmentLeader_id"] = departmentLeaderID
		row["department_id"] = departmentID
		row["template_id"] = tr.Template.ID
		row["template_name"] = tr.Template.Name
		row["template_maxscore"] = tr.Template.ScoreLimit
		if tr.Record == nil {
			row["record_score"] = "暂无评分"
		} else {
			row["record_score"] = tr.Record.Score
		}
		list[i] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ScoreDepartmentorRecordInfoController) Edit() {
	templateID, _ := c.GetInt("template_id", 0)
	departmentID, _ := c.GetInt("department_id", 0)
	projectLeaderID, _ := c.GetInt("projectLeader_id", 0)       // 项目负责人（打分）
	departmentLeaderID, _ := c.GetInt("departmentLeader_id", 0) // 部门负责人（被打分）
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecord := logic.SearchDepartmentLeaderScoreTemplateRecordByTID(year, quarter, departmentLeaderID, projectLeaderID, templateID, departmentID)

	if templateRecord == nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	row["projectLeader_id"] = projectLeaderID
	row["departmentLeader_id"] = departmentLeaderID
	row["department_id"] = departmentID
	row["template_id"] = templateRecord.Template.ID
	row["template_name"] = templateRecord.Template.Name
	row["template_maxscore"] = templateRecord.Template.ScoreLimit
	if templateRecord.Record != nil {
		row["record_score"] = templateRecord.Record.Score
	} else {
		row["record_score"] = "暂无评分"
	}

	c.Data["Source"] = row
	c.Data["pageTitle"] = "对部门负责人评分详情"
	c.display()
}

//存储资源
func (c *ScoreDepartmentorRecordInfoController) AjaxSave() {
	templateID, _ := c.GetInt("template_id", 0)
	recordScore, _ := c.GetFloat("record_score", 0)
	projectLeaderID, _ := c.GetInt("projectLeader_id", 0)       // 项目负责人（打分）
	departmentLeaderID, _ := c.GetInt("departmentLeader_id", 0) // 部门负责人（被打分）
	departmentID, _ := c.GetInt("department_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecord := logic.SearchDepartmentLeaderScoreTemplateRecordByTID(year, quarter, departmentLeaderID, projectLeaderID, templateID, departmentID)

	// 修改细分记录
	if templateRecord.Record == nil {
		record := new(models.DepartmentLeaderScoreRecord)
		record.UID = departmentLeaderID
		record.SUID = projectLeaderID
		record.Score = recordScore
		record.TID = templateID
		record.Year = year
		record.Quarter = quarter
		record.DepartmentID = departmentID

		if err := models.AddDepartmentLeaderScoreRecord(record); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		record := templateRecord.Record
		record.Score = recordScore
		if err := record.Update(); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	// 更新总记录
	err := logic.SaveDepartmentLeaderSumScoreRecord(year, quarter, departmentLeaderID, projectLeaderID, departmentID)
	if err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}

	// 如果已发布，更新发布信息
	filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
	filter3 := models.DBFilter{Key: "uid", Value: departmentLeaderID}     // 用户ID
	filter4 := models.DBFilter{Key: "department_id", Value: departmentID} // 部门ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchDepartmentLeaderReleaseRecordsByFilters(filters...)
	if len(records) != 0 {
		err := logic.ReleaseDepartmentLeaderScoreRecord(year, quarter, departmentLeaderID, departmentID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "")
}
