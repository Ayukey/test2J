package logic

import "jg2j_server/models"

// 项目一级评分模版及记录
type ProjectTemplate1Record struct {
	Template *models.ProjectTemplate1
	Record   *models.ProjectScoreRecord1
}

// 项目二级评分模版及记录
type ProjectTemplate2Record struct {
	Template *models.ProjectTemplate2
	Record   *models.ProjectScoreRecord2
}

// 项目三级评分模版及记录
type ProjectTemplate3Record struct {
	Template *models.ProjectTemplate3
	Record   *models.ProjectScoreRecord3
}

// 根据年、季度、项目ID查询一级项目模版评分记录
func SearchProjectTemplate1Records(year, quarter, pid int) []*ProjectTemplate1Record {
	templates := models.SearchAllProjectTemplate1s()
	templateRecords := make([]*ProjectTemplate1Record, len(templates))
	for i, t := range templates {
		templateRecord := new(ProjectTemplate1Record)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}
		filter3 := models.DBFilter{Key: "pid", Value: pid}
		filter4 := models.DBFilter{Key: "t1id", Value: t.ID}
		filters := []models.DBFilter{filter1, filter2, filter3, filter4}

		records := models.SearchProjectScoreRecord1sByFilters(filters...)

		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据年、季度、项目ID、一级模版ID，查询二级项目模版评分记录
func SearchProjectTemplate2Records(year, quarter, t1id, pid int) []*ProjectTemplate2Record {
	templates := models.SearchProjectTemplate2sByTID(t1id)
	templateRecords := make([]*ProjectTemplate2Record, len(templates))
	for i, t := range templates {
		templateRecord := new(ProjectTemplate2Record)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}
		filter3 := models.DBFilter{Key: "t1id", Value: t1id}
		filter4 := models.DBFilter{Key: "pid", Value: pid}
		filter5 := models.DBFilter{Key: "t2id", Value: t.ID}
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

		records := models.SearchProjectScoreRecord2sByFilters(filters...)
		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		templateRecords[i] = templateRecord
	}
	return templateRecords
}

// 根据年、季度、项目ID、一级模版ID、二级模版ID,查询三级项目模版评分记录
func SearchProjectTemplate3Records(year, quarter, t1id, t2id, pid int) []*ProjectTemplate3Record {
	templates := models.SearchProjectTemplate3sByTID(t2id)
	templateRecords := make([]*ProjectTemplate3Record, len(templates))
	for i, t := range templates {
		templateRecord := new(ProjectTemplate3Record)
		templateRecord.Template = t

		filter1 := models.DBFilter{Key: "year", Value: year}
		filter2 := models.DBFilter{Key: "quarter", Value: quarter}
		filter3 := models.DBFilter{Key: "t1id", Value: t1id}
		filter4 := models.DBFilter{Key: "t2id", Value: t2id}
		filter5 := models.DBFilter{Key: "t3id", Value: t.ID}
		filter6 := models.DBFilter{Key: "pid", Value: pid}
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

		records := models.SearchProjectScoreRecord3sByFilters(filters...)
		if len(records) == 1 {
			templateRecord.Record = records[0]
		}
		templateRecords[i] = templateRecord
	}
	return templateRecords
}
