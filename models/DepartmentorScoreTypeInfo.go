package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门负责人评分模版
type DepartmentorScoreTypeInfo struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	ScoreLimit float64 `orm:"column(score_Limit)"`
	Status     int     `orm:"column(status)"`
}

func (d *DepartmentorScoreTypeInfo) TableName() string {
	return TableName("departmentor_score_type_info")
}

func SearchDepartmentorScoreTypeInfoList(page, pageSize int, filters ...interface{}) ([]*DepartmentorScoreTypeInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*DepartmentorScoreTypeInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_score_type_info"))
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

func AddDepartmentorScoreTypeInfo(d *DepartmentorScoreTypeInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(d)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchDepartmentorScoreTypeInfoByID(id int) (*DepartmentorScoreTypeInfo, error) {
	d := new(DepartmentorScoreTypeInfo)
	err := orm.NewOrm().QueryTable(TableName("departmentor_score_type_info")).Filter("id", id).One(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DepartmentorScoreTypeInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(d, fields...); err != nil {
		return err
	}
	return nil
}

func SearchAllDepartmentorScoreTypeInfoList() []*DepartmentorScoreTypeInfo {
	list := make([]*DepartmentorScoreTypeInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_score_type_info"))
	query.OrderBy("id").All(&list)
	return list
}
