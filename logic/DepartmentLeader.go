package logic

import (
	"errors"
	"jg2j_server/models"
)

// 部门负责人互评记录（基于单个有效项目负责人）
type DepartmentLeaderTemplateRecord struct {
	Template *models.DepartmentLeaderTemplate    // 模版
	Record   *models.DepartmentLeaderScoreRecord // 基于单用户的评分记录
}

// 部门负责人互评记录（基于所有有效项目负责人的平均分）
type DepartmentLeaderTemplateAverageRecord struct {
	Template *models.DepartmentLeaderTemplate // 模版
	Score    float64                          // 计算所有该模版的评分的平均分
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度、部门ID，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchDepartmentLeaderScoreTemplateRecords(year, quarter, uid, suid, departmentId int) []*DepartmentLeaderTemplateRecord {
	templates := models.SearchAllDepartmentLeaderTemplates()
	templateRecords := make([]*DepartmentLeaderTemplateRecord, len(templates))
	for i, t := range templates {
		templateRecord := new(DepartmentLeaderTemplateRecord)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "suid", Value: suid}                  // 评分者（项目负责人）
		filter5 := models.DBFilter{Key: "tid", Value: t.ID}                   // 模版ID
		filter6 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

		records := models.SearchDepartmentLeaderScoreRecordsByFilters(filters...)

		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID、部门ID，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchDepartmentLeaderScoreTemplateRecordByTID(year, quarter, uid, suid, tid, departmentId int) *DepartmentLeaderTemplateRecord {
	template, err := models.SearchDepartmentLeaderTemplateByID(tid)
	if err == nil {
		templateRecord := new(DepartmentLeaderTemplateRecord)
		templateRecord.Template = &template

		filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "suid", Value: suid}                  // 评分者（项目负责人）
		filter5 := models.DBFilter{Key: "tid", Value: template.ID}            // 模版ID
		filter6 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

		records := models.SearchDepartmentLeaderScoreRecordsByFilters(filters...)
		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		return templateRecord
	} else {
		return nil
	}
}

// 被评分人（部门负责人ID）、年、季度、部门ID，查询该季度下项目负责人的评分模版和所有有效项目负责人打分的平均分
func SearchDepartmentLeaderTemplateAverageRecords(year, quarter, uid, departmentId int) []*DepartmentLeaderTemplateAverageRecord {
	templates := models.SearchAllDepartmentLeaderTemplates()
	templateRecords := make([]*DepartmentLeaderTemplateAverageRecord, len(templates))
	for i, t := range templates {
		templateRecord := new(DepartmentLeaderTemplateAverageRecord)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "tid", Value: t.ID}                   // 模版ID
		filter5 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

		records := models.SearchDepartmentLeaderScoreRecordsByFilters(filters...)
		count := float64(len(records))

		// 计算平均数
		total := 0.0
		for _, r := range records {
			total += r.Score
		}

		templateRecord.Score = total / count
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID、部门ID，保存当前季度的评分记录
func SaveDepartmentLeaderSumScoreRecord(year, quarter, uid, suid, departmentId int) error {
	templateRecords := SearchDepartmentLeaderScoreTemplateRecords(year, quarter, uid, suid, departmentId)

	totalScore := 0.0
	for _, templateRecord := range templateRecords {
		if templateRecord.Record != nil {
			totalScore += templateRecord.Record.Score
		}
	}

	filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
	filter3 := models.DBFilter{Key: "suid", Value: suid}                  // 评分者（项目负责人）
	filter4 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
	filter5 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

	records := models.SearchDepartmentLeaderSumScoreRecordsByFilters(filters...)

	if len(records) == 0 {
		record := new(models.DepartmentLeaderSumScoreRecord)
		record.UID = uid
		record.SUID = suid
		record.Score = totalScore
		record.Year = year
		record.Quarter = quarter
		record.DepartmentID = departmentId

		if err := models.AddDepartmentLeaderSumScoreRecord(record); err != nil {
			return err
		}
	} else {
		record := records[0]
		record.Score = totalScore
		if err := record.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 发布部门负责人季度互评
func ReleaseDepartmentLeaderScoreRecord(year, quarter, uid, departmentId int) error {

	projectLeaders := models.SearchAllProjectLeadersInActive()

	notScoreProjectLeaders := make([]*models.ProjectLeader, 0)

	totalScore := 0.0

	for _, projectLeader := range projectLeaders {
		user := projectLeader.User

		filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
		filter3 := models.DBFilter{Key: "suid", Value: user.ID}               // 评分者（项目负责人）
		filter4 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
		filter5 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

		records := models.SearchDepartmentLeaderSumScoreRecordsByFilters(filters...)
		if len(records) == 0 {
			notScoreProjectLeaders = append(notScoreProjectLeaders, projectLeader)
		} else {
			totalScore += records[0].Score
		}
	}

	if totalScore == 0 {
		return errors.New("该部门负责人该季度无任何评分")
	}

	if len(notScoreProjectLeaders) > 0 {
		errorStr := "以下项目负责人还未完成互评: <br>"
		for _, projectLeader := range notScoreProjectLeaders {
			errorStr += projectLeader.User.Name + "(" + projectLeader.User.Account + ")" + "<br>"
		}
		return errors.New(errorStr)
	}

	filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
	filter4 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchDepartmentLeaderReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		// 没有发布记录，新增
		record := new(models.DepartmentLeaderReleaseRecord)
		record.UID = uid
		record.Score = totalScore / float64(len(projectLeaders))
		record.Year = year
		record.Quarter = quarter
		record.DepartmentID = departmentId
		if err := models.AddDepartmentLeaderReleaseRecord(record); err != nil {
			return err
		}
	} else {
		// 有发布记录，更新
		record := records[0]
		record.Score = totalScore / float64(len(projectLeaders))
		if err := record.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 检查部门负责人互评是否一发布
func CheckDepartmentLeaderReleaseReacord(year, quarter, uid, departmentId int) bool {
	filter1 := models.DBFilter{Key: "year", Value: year}                  // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}            // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}                    // 被评分者（部门负责人）
	filter4 := models.DBFilter{Key: "department_id", Value: departmentId} // 部门ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchDepartmentLeaderReleaseRecordsByFilters(filters...)
	return len(records) > 0
}
