package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// 部门负责人评分发布记录
type DepartmentorSumPubInfo struct {
	ID           int     `orm:"column(id)"`
	UserID       int     `orm:"column(user_id)"`
	Year         int     `orm:"column(year)"`
	Quarter      int     `orm:"column(quarter)"`
	Score        float64 `orm:"column(score)"`
	DepartmentID int     `orm:"column(department_id)"`
}

func (d *DepartmentorSumPubInfo) TableName() string {
	return TableName("departmentor_sumpub_info")
}

func SearchDepartmentorSumPubInfoList(page, pageSize int, filters ...interface{}) ([]*DepartmentorSumPubInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*DepartmentorSumPubInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_sumpub_info"))
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

func AddDepartmentorSumPubInfo(d *DepartmentorSumPubInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(d)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchDepartmentorSumPubInfoByID(id int) (*DepartmentorSumPubInfo, error) {
	d := new(DepartmentorSumPubInfo)
	err := orm.NewOrm().QueryTable(TableName("departmentor_sumpub_info")).Filter("id", id).One(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DepartmentorSumPubInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(d, fields...); err != nil {
		return err
	}
	return nil
}

func SearchDepartmentorSumPubInfoByFilters(filters ...interface{}) []*DepartmentorSumPubInfo {
	list := make([]*DepartmentorSumPubInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_sumpub_info"))
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

func SearchDepartmentorSumPubInfoByOrder(filters ...interface{}) []*DepartmentorSumPubInfo {
	list := make([]*DepartmentorSumPubInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("departmentor_sumpub_info"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			fmt.Println(filters[k])
			fmt.Println(filters[k].(string))
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("score").All(&list)
	return list
}
