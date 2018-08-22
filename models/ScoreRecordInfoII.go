package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目二级评分记录
type ScoreRecordInfoII struct {
	ID          int     `orm:"column(id)"`
	ProjectID   int     `orm:"column(project_id)"`
	ScoreTypeID int     `orm:"column(scoretype_id)"`
	Year        int     `orm:"column(year)"`
	Quarter     int     `orm:"column(quarter)"`
	TotalScore  float64 `orm:"column(total_score)"`
	TID         int     `orm:"column(tid)"`
	UpdateDate  float64 `orm:"column(update_date)"`
	Remark      string  `orm:"column(remark)"`
}

func (s *ScoreRecordInfoII) TableName() string {
	return TableName("score_records2")
}

// 分页查询项目二级评分记录
func SearchScoreRecordInfoIIList(page, pageSize int, filters ...interface{}) ([]*ScoreRecordInfoII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreRecordInfoII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records2"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

// 添加项目二级评分记录
func AddScoreRecordInfoII(s *ScoreRecordInfoII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 更新项目二级评分记录
func (s *ScoreRecordInfoII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

// 查询项目二级评分记录
func SearchScoreRecordInfoIIByFilters(filters ...interface{}) []*ScoreRecordInfoII {
	list := make([]*ScoreRecordInfoII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records2"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}
