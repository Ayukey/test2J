package logic

import (
	"errors"
	"jg2j_server/models"
)

// 项目负责人互评记录（基于单个部门负责人）
type ProjectLeaderTemplateRecord struct {
	Template *models.ProjectLeaderTemplate
	Record   *models.ProjectLeaderScoreRecord
}

// 项目负责人互评记录（基于所有部门负责人的平均分）
type ProjectLeaderTemplateAverageRecord struct {
	Template *models.ProjectLeaderTemplate
	Score    float64
}

// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，查询该季度下项目负责人的评分模版和对应的评分记录
func SearchProjectLeaderTemplateRecords(year, quarter, uid, suid, projectId int) []*ProjectLeaderTemplateRecord {
	templates := models.SearchAllProjectLeaderTemplates()
	templateRecords := make([]*ProjectLeaderTemplateRecord, len(templates))
	for i, t := range templates {
		templateRecord := new(ProjectLeaderTemplateRecord)
		templateRecord.Template = t
		// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，评分模版ID，查询该季度下项目负责人的评分模版和对应的评分记录

		filter1 := models.DBFilter{Key: "year", Value: year}            // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}      // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}              // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "suid", Value: suid}            // 评分者（项目负责人）
		filter5 := models.DBFilter{Key: "tid", Value: t.ID}             // 模版ID
		filter6 := models.DBFilter{Key: "project_id", Value: projectId} // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

		records := models.SearchProjectLeaderScoreRecordsByFilters(filters...)

		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchProjectLeaderTemplateRecordByTID(year, quarter, uid, suid, tid, projectId int) *ProjectLeaderTemplateRecord {
	template, err := models.SearchProjectLeaderTemplateByID(tid)
	if err == nil {
		templateRecord := new(ProjectLeaderTemplateRecord)
		templateRecord.Template = &template

		filter1 := models.DBFilter{Key: "year", Value: year}            // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}      // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}              // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "suid", Value: suid}            // 评分者（项目负责人）
		filter5 := models.DBFilter{Key: "tid", Value: template.ID}      // 模版ID
		filter6 := models.DBFilter{Key: "project_id", Value: projectId} // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

		records := models.SearchProjectLeaderScoreRecordsByFilters(filters...)
		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		return templateRecord
	} else {
		return nil
	}
}

// 被评分人（项目负责人ID）、年、季度，查询该季度下项目负责人的评分模版和所有部门负责人打分的平均分
func SearchProjectLeaderTemplateAverageRecords(year, quarter, uid, projectId int) []*ProjectLeaderTemplateAverageRecord {
	templates := models.SearchAllProjectLeaderTemplates()
	templateRecords := make([]*ProjectLeaderTemplateAverageRecord, len(templates))
	for i, t := range templates {
		templateRecord := new(ProjectLeaderTemplateAverageRecord)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}            // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}      // 季度
		filter3 := models.DBFilter{Key: "uid", Value: uid}              // 被评分者（部门负责人）
		filter4 := models.DBFilter{Key: "tid", Value: t.ID}             // 模版ID
		filter5 := models.DBFilter{Key: "project_id", Value: projectId} // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

		records := models.SearchProjectLeaderScoreRecordsByFilters(filters...)
		count := float64(len(records))

		total := 0.0

		for _, r := range records {
			total += r.Score
		}

		templateRecord.Score = total / count
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，保存当前季度的评分记录
func SaveProjectLeaderSumScoreRecord(year, quarter, uid, suid, projectId int) error {
	templateRecords := SearchProjectLeaderTemplateRecords(year, quarter, uid, suid, projectId)

	totalScore := 0.0
	for _, templateRecord := range templateRecords {
		if templateRecord.Record != nil {
			totalScore += templateRecord.Record.Score
		}
	}

	filter1 := models.DBFilter{Key: "year", Value: year}           // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}     // 季度
	filter3 := models.DBFilter{Key: "suid", Value: suid}           // 评分者（项目负责人）
	filter4 := models.DBFilter{Key: "uid", Value: uid}             // 被评分者（部门负责人）
	filter5 := models.DBFilter{Key: "projectId", Value: projectId} // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

	records := models.SearchProjectLeaderSumScoreRecordsByFilters(filters...)

	if len(records) == 0 {
		record := new(models.ProjectLeaderSumScoreRecord)
		record.UID = uid
		record.SUID = suid
		record.Score = totalScore
		record.Year = year
		record.Quarter = quarter
		record.ProjectID = projectId

		if err := models.AddProjectLeaderSumScoreRecord(record); err != nil {
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

// 发布项目负责人季度互评
func ReleaseProjectLeaderScoreRecord(year, quarter, uid, projectId int) error {

	departmentLeaders := models.SearchAllDepartmentLeadersInActive(year, quarter)

	notScoreDepartmentLeaders := make([]*models.DepartmentLeader, 0)

	totalScore := 0.0

	for _, departmentLeader := range departmentLeaders {
		user := departmentLeader.User

		filter1 := models.DBFilter{Key: "year", Value: year}           // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}     // 季度
		filter3 := models.DBFilter{Key: "suid", Value: user.ID}        // 评分者（项目负责人）
		filter4 := models.DBFilter{Key: "uid", Value: uid}             // 被评分者（部门负责人）
		filter5 := models.DBFilter{Key: "projectId", Value: projectId} // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

		records := models.SearchProjectLeaderSumScoreRecordsByFilters(filters...)
		if len(records) == 0 {
			notScoreDepartmentLeaders = append(notScoreDepartmentLeaders, departmentLeader)
		} else {
			totalScore += records[0].Score
		}
	}

	if totalScore == 0 {
		return errors.New("该项目负责人该季度无任何评分")
	}

	if len(notScoreDepartmentLeaders) > 0 {
		errorStr := "以下部门负责人还未完成互评: <br>"
		for _, departmentLeader := range notScoreDepartmentLeaders {
			errorStr += departmentLeader.User.Name + "(" + departmentLeader.User.Account + ")" + "<br>"
		}
		return errors.New(errorStr)
	}

	filter1 := models.DBFilter{Key: "year", Value: year}           // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}     // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}             // 被评分者（部门负责人）
	filter4 := models.DBFilter{Key: "projectId", Value: projectId} // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchProjectLeaderReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		// 没有发布记录，新增
		record := new(models.ProjectLeaderReleaseRecord)
		record.UID = uid
		record.Score = totalScore / float64(len(departmentLeaders))
		record.Year = year
		record.Quarter = quarter
		record.ProjectID = projectId
		if err := models.AddProjectLeaderReleaseRecord(record); err != nil {
			return err
		}
	} else {
		// 有发布记录，更新
		record := records[0]
		record.Score = totalScore / float64(len(departmentLeaders))
		if err := record.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 检查部门负责人互评是否一发布
func CheckProjectLeaderReleaseReacord(year, quarter, uid, projectId int) bool {
	filter1 := models.DBFilter{Key: "year", Value: year}           // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}     // 季度
	filter3 := models.DBFilter{Key: "uid", Value: uid}             // 被评分者（部门负责人）
	filter4 := models.DBFilter{Key: "projectId", Value: projectId} // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4}

	records := models.SearchProjectLeaderReleaseRecordsByFilters(filters...)
	return len(records) > 0
}
