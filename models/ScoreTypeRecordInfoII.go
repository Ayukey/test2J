package models

type ScoreTypeRecordInfoII struct {
	Type   *ScoreTypeInfoII
	Record *ScoreRecordInfoII
}

func SerachScoreTypeRecordInfoIIList(year, quarter, tid, projectID int) []*ScoreTypeRecordInfoII {
	typeList := SearchScoreTypeInfoIIByTID(tid)
	typeRecordList := make([]*ScoreTypeRecordInfoII, len(typeList))
	for i, t := range typeList {
		typeRecord := new(ScoreTypeRecordInfoII)
		typeRecord.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "tid", tid)
		filters = append(filters, "project_id", projectID)
		filters = append(filters, "scoretype_id", t.ID)
		recordList := SearchScoreRecordInfoIIByFilters(filters...)
		if len(recordList) == 1 {
			typeRecord.Record = recordList[0]
		}
		typeRecordList[i] = typeRecord
	}
	return typeRecordList
}
