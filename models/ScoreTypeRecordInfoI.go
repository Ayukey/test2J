package models

type ScoreTypeRecordInfoI struct {
	Type   *ScoreTypeInfoI
	Record *ScoreRecordInfoI
}

func SerachScoreTypeRecordInfoIList(year, quarter, projectID int) []*ScoreTypeRecordInfoI {
	typeList := SearchAllScoreTypeInfoIList()
	typeRecordList := make([]*ScoreTypeRecordInfoI, len(typeList))
	for i, t := range typeList {
		typeRecord := new(ScoreTypeRecordInfoI)
		typeRecord.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "project_id", projectID)
		filters = append(filters, "scoretype_id", t.ID)
		recordList := SearchScoreRecordInfoIByFilters(filters...)
		if len(recordList) == 1 {
			typeRecord.Record = recordList[0]
		}
		typeRecordList[i] = typeRecord
	}
	return typeRecordList
}

func SerachAllScoreTypeRecordInfoIList(year, quarter int) []*ScoreTypeRecordInfoI {
	typeList := SearchAllScoreTypeInfoIList()
	typeRecordList := make([]*ScoreTypeRecordInfoI, len(typeList))
	for i, t := range typeList {
		typeRecord := new(ScoreTypeRecordInfoI)
		typeRecord.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "scoretype_id", t.ID)
		recordList := SearchScoreRecordInfoIByFilters(filters...)
		if len(recordList) == 1 {
			typeRecord.Record = recordList[0]
		}
		typeRecordList[i] = typeRecord
	}
	return typeRecordList
}
