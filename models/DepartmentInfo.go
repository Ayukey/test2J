package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门信息
type DepartmentInfo struct {
	ID     int    `orm:"column(id)"`
	Name   string `orm:"column(name)"`
	UserID int    `orm:"column(user_id)"`
	Status int    `orm:"column(status)"`
}

func (d *DepartmentInfo) TableName() string {
	return TableName("department_info")
}

func SearchDepartmentInfoList(page, pageSize int, filters ...interface{}) ([]*DepartmentInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*DepartmentInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("department_info"))
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

func SearchAllDepartmentInfo() []*DepartmentInfo {
	list := make([]*DepartmentInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("department_info"))
	query.OrderBy("id").All(&list)
	return list
}

func AddDepartmentInfo(d *DepartmentInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(d)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchDepartmentInfoByID(id int) (*DepartmentInfo, error) {
	d := new(DepartmentInfo)
	err := orm.NewOrm().QueryTable(TableName("department_info")).Filter("id", id).One(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func SearchDepartmentInfoByName(name string) (*DepartmentInfo, error) {
	d := new(DepartmentInfo)
	err := orm.NewOrm().QueryTable(TableName("department_info")).Filter("name", name).One(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DepartmentInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(d, fields...); err != nil {
		return err
	}
	return nil
}
