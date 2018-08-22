package logic

import (
	"errors"
	"jg2j_server/models"
)

// 项目负责人互评记录（基于单个部门负责人）
type ProjectorScoreTypeRecordInfo struct {
	Type   *models.ProjectorScoreTypeInfo
	Record *models.ProjectorScoreRecords
}

// 项目负责人互评记录（基于所有部门负责人的平均分）
type ProjectorScoreTypeRecordFinalInfo struct {
	Type  *models.ProjectorScoreTypeInfo
	Score float64
}

// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，查询该季度下项目负责人的评分模版和对应的评分记录
func SearchProjectorScoreTypeRecordInfos(year, quarter, uid, suid int) []*ProjectorScoreTypeRecordInfo {
	typeInfos := models.SearchAllProjectorScoreTypeInfoList()
	typeRecordInfos := make([]*ProjectorScoreTypeRecordInfo, len(typeInfos))
	for i, t := range typeInfos {
		typeRecordInfo := new(ProjectorScoreTypeRecordInfo)
		typeRecordInfo.Type = t
		// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，评分模版ID，查询该季度下项目负责人的评分模版和对应的评分记录
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "tid", t.ID)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "scoreuser_id", suid)
		recordInfos := models.SearchProjectorScoreRecordsByFilters(filters...)

		if len(recordInfos) == 1 {
			typeRecordInfo.Record = recordInfos[0]
		}
		typeRecordInfos[i] = typeRecordInfo
	}
	return typeRecordInfos
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchSingleProjectorScoreTypeRecordInfoByTID(year, quarter, uid, suid, tid int) *ProjectorScoreTypeRecordInfo {
	typeInfo, err := models.SearchProjectorScoreTypeInfoByID(tid)
	if err == nil {
		typeRecordInfo := new(ProjectorScoreTypeRecordInfo)
		typeRecordInfo.Type = typeInfo
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "tid", typeInfo.ID)
		filters = append(filters, "scoreuser_id", suid)
		recordInfos := models.SearchProjectorScoreRecordsByFilters(filters...)
		if len(recordInfos) == 1 {
			typeRecordInfo.Record = recordInfos[0]
		}
		return typeRecordInfo
	} else {
		return nil
	}
}

// 被评分人（项目负责人ID）、年、季度，查询该季度下项目负责人的评分模版和所有部门负责人打分的平均分
func SearchProjectorScoreTypeRecordInfosBySumData(year, quarter, uid int) []*ProjectorScoreTypeRecordFinalInfo {
	typeInfos := models.SearchAllProjectorScoreTypeInfoList()
	typeRecordInfos := make([]*ProjectorScoreTypeRecordFinalInfo, len(typeInfos))
	for i, t := range typeInfos {
		typeRecordInfo := new(ProjectorScoreTypeRecordFinalInfo)
		typeRecordInfo.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "tid", t.ID)
		recordInfos := models.SearchProjectorScoreRecordsByFilters(filters...)
		total := 0.0

		for _, r := range recordInfos {
			total += r.Score
		}

		typeRecordInfo.Score = total / float64(len(recordInfos))
		typeRecordInfos[i] = typeRecordInfo
	}
	return typeRecordInfos
}

// 根据评分人（部门负责人ID）、被评分人（项目负责人ID）、年、季度，保存当前季度的评分记录
func SaveProjectorScoreBySingleDepartmentor(year, quarter, uid, suid int) error {
	typeRecordInfos := SearchProjectorScoreTypeRecordInfos(year, quarter, uid, suid)

	user, _ := models.SearchUserInfoByID(uid)
	if user == nil {
		return errors.New("未找到该项目负责人的信息")
	}

	totalScore := 0.0
	for _, typeRecordInfo := range typeRecordInfos {
		if typeRecordInfo.Record != nil {
			totalScore += typeRecordInfo.Record.Score
		}
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "scoreuser_id", suid)
	filters = append(filters, "user_id", uid)
	sumInfos := models.SearchProjectorScoreSumInfoByFilters(filters...)

	if len(sumInfos) == 0 {
		sumInfo := new(models.ProjectorScoreSumInfo)
		sumInfo.UserID = uid
		sumInfo.ScoreUserID = suid
		sumInfo.Score = totalScore
		sumInfo.Year = year
		sumInfo.Quarter = quarter
		sumInfo.ProjectID = user.ProjectID

		if _, err := models.AddProjectorScoreSumInfo(sumInfo); err != nil {
			return err
		}
	} else {
		sumInfo := sumInfos[0]
		sumInfo.Score = totalScore
		if err := sumInfo.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 发布项目负责人季度互评
func ReleaseProjectorScore(year, quarter, uid int) error {
	// 查询负责人信息
	user, _ := models.SearchUserInfoByID(uid)
	if user == nil {
		return errors.New("未找到该项目负责人的信息")
	}
	// 查询项目信息
	project, _ := models.SearchProjectInfoByID(user.ProjectID)
	if project == nil {
		return errors.New("未找到该项目负责人的项目信息")
	}

	departmentors := models.SearchAllDepartmentLeaders()

	norecordDepartmentors := make([]*models.UserInfo, 0)

	totalScore := 0.0

	for _, user := range departmentors {
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "scoreuser_id", user.ID)
		sumInfo := models.SearchProjectorScoreSumInfoByFilters(filters...)
		if len(sumInfo) == 0 {
			norecordDepartmentors = append(norecordDepartmentors, user)
		} else {
			totalScore += sumInfo[0].Score
		}
	}

	if totalScore == 0 {
		return errors.New("该项目负责人该季度无任何评分")
	}

	if len(norecordDepartmentors) > 0 {
		errorStr := "以下部门负责人还未完成互评: <br>"
		for _, u := range norecordDepartmentors {
			errorStr += u.Name + "(" + u.Account + ")" + "<br>"
		}
		return errors.New(errorStr)
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", uid)
	sumPubInfos := models.SearchProjectorSumPubInfoByFilters(filters...)
	if len(sumPubInfos) == 0 {
		// 没有发布记录，新增
		sumPubInfo := new(models.ProjectorSumPubInfo)
		sumPubInfo.UserID = uid
		sumPubInfo.Score = totalScore / float64(len(departmentors))
		sumPubInfo.Year = year
		sumPubInfo.Quarter = quarter
		sumPubInfo.ProjectID = project.ID
		if _, err := models.AddProjectorSumPubInfo(sumPubInfo); err != nil {
			return err
		}
	} else {
		// 有发布记录，更新
		sumPubInfo := sumPubInfos[0]
		sumPubInfo.Score = totalScore / float64(len(departmentors))
		if err := sumPubInfo.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 检查部门负责人互评是否一发布
func CheckProjectorScoreReleaseStatus(year, quarter, uid int) bool {
	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", uid)
	sumPubInfos := models.SearchProjectorSumPubInfoByFilters(filters...)
	return len(sumPubInfos) > 0
}
