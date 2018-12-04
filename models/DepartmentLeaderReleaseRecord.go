package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门负责人评分发布记录
type DepartmentLeaderReleaseRecord struct {
	ID           int     `orm:"column(id)"`
	UID          int     `orm:"column(uid)"`
	Year         int     `orm:"column(year)"`
	Quarter      int     `orm:"column(quarter)"`
	Score        float64 `orm:"column(score)"`
	DepartmentID int     `orm:"column(department_id)"`
}

// 新增部门负责人评分发布记录
func AddDepartmentLeaderReleaseRecord(d *DepartmentLeaderReleaseRecord) error {
	_, err := orm.NewOrm().Insert(d)
	return err
}

// 更新部门负责人评分发布记录
func (d *DepartmentLeaderReleaseRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(d, fields...)
	return err
}

// 根据条件查询部门负责人季度评分发布记录
func SearchDepartmentLeaderReleaseRecordsByFilters(filters ...DBFilter) []*DepartmentLeaderReleaseRecord {
	records := make([]*DepartmentLeaderReleaseRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("department_leader_release_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}

// 根据条件查询部门负责人季度评分发布记录(按分数排序)
func SearchDepartmentLeaderReleaseRecordsByOrder(filters ...DBFilter) []*DepartmentLeaderReleaseRecord {
	records := make([]*DepartmentLeaderReleaseRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("department_leader_release_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-score").All(&records)
	return records
}
