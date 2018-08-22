package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// 项目负责人评分记录
type ProjectorScoreRecords struct {
	ID          int     `orm:"column(id)"`
	ScoreUserID int     `orm:"column(scoreuser_id)"`
	UserID      int     `orm:"column(user_id)"`
	TID         int     `orm:"column(tid)"`
	Score       float64 `orm:"column(score)"`
	Year        int     `orm:"column(year)"`
	Quarter     int     `orm:"column(quarter)"`
	ProjectID   int     `orm:"column(project_id)"`
}

func (p *ProjectorScoreRecords) TableName() string {
	return TableName("projector_score_records")
}

func SearchProjectorScoreRecordsList(page, pageSize int, filters ...interface{}) ([]*ProjectorScoreRecords, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ProjectorScoreRecords, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_records"))
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

func AddProjectorScoreRecords(p *ProjectorScoreRecords) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchProjectorScoreRecordsByID(id int) (*ProjectorScoreRecords, error) {
	p := new(ProjectorScoreRecords)
	err := orm.NewOrm().QueryTable(TableName("projector_score_records")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *ProjectorScoreRecords) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func SearchProjectorScoreRecordsByFilters(filters ...interface{}) []*ProjectorScoreRecords {
	list := make([]*ProjectorScoreRecords, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_records"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			fmt.Println(filters[k])
			fmt.Println(filters[k+1])
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}
