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
	projectorID, _ := c.GetInt("projectorID", 0)       // 项目负责人（打分）
	departmentorID, _ := c.GetInt("departmentorID", 0) // 部门负责人（被打分）
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordInfos := logic.SearchDepartmentorScoreTypeRecordInfos(year, quarter, departmentorID, projectorID)

	list := make([]map[string]interface{}, len(typeRecordInfos))

	for i, typeRecordInfo := range typeRecordInfos {
		row := make(map[string]interface{})
		row["id"] = typeRecordInfo.Type.ID
		row["name"] = typeRecordInfo.Type.Name
		row["max_score"] = typeRecordInfo.Type.ScoreLimit
		if typeRecordInfo.Record == nil {
			row["score"] = "暂无评分"
		} else {
			row["score"] = typeRecordInfo.Record.Score
		}
		list[i] = row
	}
	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ScoreDepartmentorRecordInfoController) Edit() {
	tid, _ := c.GetInt("tid", 0)
	projectorID, _ := c.GetInt("projectorID", 0)       // 项目负责人（打分）
	departmentorID, _ := c.GetInt("departmentorID", 0) // 部门负责人（被打分）
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordInfo := logic.SearchSingleDepartmentorScoreTypeRecordInfoByTID(year, quarter, departmentorID, projectorID, tid)

	if typeRecordInfo == nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	row["projectorID"] = projectorID
	row["departmentorID"] = departmentorID
	row["tid"] = typeRecordInfo.Type.ID
	row["name"] = typeRecordInfo.Type.Name
	row["maxscore"] = typeRecordInfo.Type.ScoreLimit
	if typeRecordInfo.Record != nil {
		row["score"] = typeRecordInfo.Record.Score
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
	projectorID, _ := c.GetInt("projectorID", 0)       // 项目负责人（打分）
	departmentorID, _ := c.GetInt("departmentorID", 0) // 部门负责人（被打分）
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	score, _ := c.GetFloat("score", 0)

	user, _ := models.SearchUserInfoByID(departmentorID)

	typeRecordInfo := logic.SearchSingleDepartmentorScoreTypeRecordInfoByTID(year, quarter, departmentorID, projectorID, tid)

	// 修改细分记录
	if typeRecordInfo.Record == nil {
		record := new(models.DepartmentorScoreRecords)
		record.UserID = departmentorID
		record.ScoreUserID = projectorID
		record.Score = score
		record.TID = tid
		record.Year = year
		record.Quarter = quarter
		record.DepartmentID = user.DepartmentID

		if _, err := models.AddDepartmentorScoreRecords(record); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		record := typeRecordInfo.Record
		record.Score = score
		if err := record.Update(); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
	}

	// 更新总记录
	err := logic.SaveDepartmentorScoreBySingleProjector(year, quarter, departmentorID, projectorID)
	if err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}

	// 如果已发布，更新发布信息
	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", departmentorID)
	sumPubInfos := models.SearchDepartmentorSumPubInfoByFilters(filters...)
	if len(sumPubInfos) != 0 {
		err := logic.ReleaseDepartmentorScore(year, quarter, departmentorID)
		if err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
	}

	c.ajaxMsg("", MSG_OK)
}
