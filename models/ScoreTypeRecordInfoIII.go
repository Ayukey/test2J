package models

type ScoreTypeRecordInfoIII struct {
	TTID   int
	Type   *ScoreTypeInfoIII
	Record *ScoreRecordInfoIII
}

func SerachScoreTypeRecordInfoIIIList(year, quarter, tid, ttid, projectID int) []*ScoreTypeRecordInfoIII {
	typeList := SearchScoreTypeInfoIIIByTID(tid)
	typeRecordList := make([]*ScoreTypeRecordInfoIII, len(typeList))
	for i, t := range typeList {
		typeRecord := new(ScoreTypeRecordInfoIII)
		typeRecord.TTID = ttid
		typeRecord.Type = t
		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "tid", tid)
		filters = append(filters, "project_id", projectID)
		filters = append(filters, "scoretype_id", t.ID)
		recordList := SearchScoreRecordInfoIIIByFilters(filters...)
		if len(recordList) == 1 {
			typeRecord.Record = recordList[0]
		}
		typeRecordList[i] = typeRecord
	}
	return typeRecordList
}
