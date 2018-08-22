package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// 项目一级评分记录
type ScoreRecordInfoI struct {
	ID          int     `orm:"column(id)"`
	ProjectID   int     `orm:"column(project_id)"`
	ScoreTypeID int     `orm:"column(scoretype_id)"`
	Year        int     `orm:"column(year)"`
	Quarter     int     `orm:"column(quarter)"`
	TotalScore  float64 `orm:"column(total_score)"`
	UpdateDate  float64 `orm:"column(update_date)"`
	Remark      string  `orm:"column(remark)"`
}

func (s *ScoreRecordInfoI) TableName() string {
	return TableName("score_records1")
}

// 分页查询项目一级评分记录
func SearchScoreRecordInfoIList(page, pageSize int, filters ...interface{}) ([]*ScoreRecordInfoI, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreRecordInfoI, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records1"))
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

// 添加项目一级评分记录
func AddScoreRecordInfoI(s *ScoreRecordInfoI) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 更新项目一级评分记录
func (s *ScoreRecordInfoI) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

// 查询项目一级评分记录
func SearchScoreRecordInfoIByFilters(filters ...interface{}) []*ScoreRecordInfoI {
	list := make([]*ScoreRecordInfoI, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records1"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			fmt.Println(filters[k])
			fmt.Println(filters[k].(string))
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}
