package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// 部门负责人评分记录
type DepartmentorScoreRecords struct {
	ID           int     `orm:"column(id)"`
	ScoreUserID  int     `orm:"column(scoreuser_id)"`
	UserID       int     `orm:"column(user_id)"`
	TID          int     `orm:"column(tid)"`
	Score        float64 `orm:"column(score)"`
	Year         int     `orm:"column(year)"`
	Quarter      int     `orm:"column(quarter)"`
	DepartmentID int     `orm:"column(department_id)"`
}

func (d *DepartmentorScoreRecords) TableName() string {
	return TableName("departmentor_score_records")
}

func SearchDepartmentorScoreRecordsList(page, pageSize int, filters ...interface{}) ([]*DepartmentorScoreRecords, int64) {
	offset := (page - 1) * pageSize
	list := make([]*DepartmentorScoreRecords, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_score_records"))
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

func AddDepartmentorScoreRecords(d *DepartmentorScoreRecords) (int64, error) {
	id, err := orm.NewOrm().Insert(d)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchDepartmentorScoreRecordsByID(id int) (*DepartmentorScoreRecords, error) {
	d := new(DepartmentorScoreRecords)
	err := orm.NewOrm().QueryTable(TableName("departmentor_score_records")).Filter("id", id).One(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DepartmentorScoreRecords) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(d, fields...); err != nil {
		return err
	}
	return nil
}

func SearchDepartmentorScoreRecordsByFilters(filters ...interface{}) []*DepartmentorScoreRecords {
	list := make([]*DepartmentorScoreRecords, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_score_records"))
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
