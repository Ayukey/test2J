package logic

import (
	"errors"
	"jg2j_server/models"
)

// 部门负责人互评记录（基于单个有效项目负责人）
type DepartmentorScoreTypeRecordInfo struct {
	Type   *models.DepartmentorScoreTypeInfo
	Record *models.DepartmentorScoreRecords
}

// 部门负责人互评记录（基于所有有效项目负责人的平均分）
type DepartmentorScoreTypeRecordFinalInfo struct {
	Type  *models.DepartmentorScoreTypeInfo
	Score float64
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchDepartmentorScoreTypeRecordInfos(year, quarter, uid, suid int) []*DepartmentorScoreTypeRecordInfo {
	typeInfos := models.SearchAllDepartmentorScoreTypeInfoList()
	typeRecordInfos := make([]*DepartmentorScoreTypeRecordInfo, len(typeInfos))
	for i, t := range typeInfos {
		typeRecordInfo := new(DepartmentorScoreTypeRecordInfo)
		typeRecordInfo.Type = t
		// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度，评分模版ID，查询该季度下部门负责人的评分模版和对应的评分记录
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "tid", t.ID)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "scoreuser_id", suid)
		recordInfos := models.SearchDepartmentorScoreRecordsByFilters(filters...)

		if len(recordInfos) == 1 {
			typeRecordInfo.Record = recordInfos[0]
		}
		typeRecordInfos[i] = typeRecordInfo
	}
	return typeRecordInfos
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID，查询该季度下部门负责人的评分模版和对应的评分记录
func SearchSingleDepartmentorScoreTypeRecordInfoByTID(year, quarter, uid, suid, tid int) *DepartmentorScoreTypeRecordInfo {
	typeInfo, err := models.SearchDepartmentorScoreTypeInfoByID(tid)
	if err == nil {
		typeRecordInfo := new(DepartmentorScoreTypeRecordInfo)
		typeRecordInfo.Type = typeInfo
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "tid", typeInfo.ID)
		filters = append(filters, "scoreuser_id", suid)
		recordInfos := models.SearchDepartmentorScoreRecordsByFilters(filters...)
		if len(recordInfos) == 1 {
			typeRecordInfo.Record = recordInfos[0]
		}
		return typeRecordInfo
	} else {
		return nil
	}
}

// 被评分人（部门负责人ID）、年、季度，查询该季度下项目负责人的评分模版和所有有效项目负责人打分的平均分
func SearchDepartmentorScoreTypeRecordInfosBySumData(year, quarter, uid int) []*DepartmentorScoreTypeRecordFinalInfo {
	typeInfos := models.SearchAllDepartmentorScoreTypeInfoList()
	typeRecordInfos := make([]*DepartmentorScoreTypeRecordFinalInfo, len(typeInfos))
	for i, t := range typeInfos {
		typeRecordInfo := new(DepartmentorScoreTypeRecordFinalInfo)
		typeRecordInfo.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "tid", t.ID)
		recordInfos := models.SearchDepartmentorScoreRecordsByFilters(filters...)
		total := 0.0

		for _, r := range recordInfos {
			total += r.Score
		}

		typeRecordInfo.Score = total / float64(len(recordInfos))
		typeRecordInfos[i] = typeRecordInfo
	}
	return typeRecordInfos
}

// 根据评分人（项目负责人ID）、被评分人（部门负责人ID）、年、季度, 模版ID，保存当前季度的评分记录
func SaveDepartmentorScoreBySingleProjector(year, quarter, uid, suid int) error {
	typeRecordInfos := SearchDepartmentorScoreTypeRecordInfos(year, quarter, uid, suid)

	user, _ := models.SearchUserInfoByID(uid)
	if user == nil {
		return errors.New("未找到该部门负责人的信息")
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
	sumInfos := models.SearchDepartmentorScoreSumInfoByFilters(filters...)

	if len(sumInfos) == 0 {
		sumInfo := new(models.DepartmentorScoreSumInfo)
		sumInfo.UserID = uid
		sumInfo.ScoreUserID = suid
		sumInfo.Score = totalScore
		sumInfo.Year = year
		sumInfo.Quarter = quarter
		sumInfo.DepartmentID = user.DepartmentID

		if _, err := models.AddDepartmentorScoreSumInfo(sumInfo); err != nil {
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

// 发布部门负责人季度互评
func ReleaseDepartmentorScore(year, quarter, uid int) error {
	// 查询负责人信息
	user, _ := models.SearchUserInfoByID(uid)
	if user == nil {
		return errors.New("未找到该部门负责人的信息")
	}
	// 查询部门信息
	department, _ := models.SearchDepartmentInfoByID(user.DepartmentID)
	if department == nil {
		return errors.New("未找到该部门负责人的部门信息")
	}

	projectors := models.SearchAllProjectLeadersInUse(year, quarter)

	norecordProjectors := make([]*models.UserInfo, 0)

	totalScore := 0.0

	for _, user := range projectors {
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", uid)
		filters = append(filters, "scoreuser_id", user.ID)
		sumInfo := models.SearchDepartmentorScoreSumInfoByFilters(filters...)
		if len(sumInfo) == 0 {
			norecordProjectors = append(norecordProjectors, user)
		} else {
			totalScore += sumInfo[0].Score
		}
	}

	if totalScore == 0 {
		return errors.New("该部门负责人该季度无任何评分")
	}

	if len(norecordProjectors) > 0 {
		errorStr := "以下项目负责人还未完成互评: <br>"
		for _, u := range norecordProjectors {
			errorStr += u.Name + "(" + u.Account + ")" + "<br>"
		}
		return errors.New(errorStr)
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", uid)
	sumPubInfos := models.SearchDepartmentorSumPubInfoByFilters(filters...)
	if len(sumPubInfos) == 0 {
		// 没有发布记录，新增
		sumPubInfo := new(models.DepartmentorSumPubInfo)
		sumPubInfo.UserID = uid
		sumPubInfo.Score = totalScore / float64(len(projectors))
		sumPubInfo.Year = year
		sumPubInfo.Quarter = quarter
		sumPubInfo.DepartmentID = department.ID
		if _, err := models.AddDepartmentorSumPubInfo(sumPubInfo); err != nil {
			return err
		}
	} else {
		// 有发布记录，更新
		sumPubInfo := sumPubInfos[0]
		sumPubInfo.Score = totalScore / float64(len(projectors))
		if err := sumPubInfo.Update(); err != nil {
			return err
		}
	}

	return nil
}

// 检查部门负责人互评是否一发布
func CheckDepartmentorScoreReleaseStatus(year, quarter, uid int) bool {
	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "user_id", uid)
	sumPubInfos := models.SearchDepartmentorSumPubInfoByFilters(filters...)
	return len(sumPubInfos) > 0
}
