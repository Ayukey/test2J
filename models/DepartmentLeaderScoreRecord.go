package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门负责人评分记录
type DepartmentLeaderScoreRecord struct {
	ID           int     `orm:"column(id)"`
	SUID         int     `orm:"column(suid)"` // 评分者ID
	UID          int     `orm:"column(uid)"`  // 被评分者ID
	TID          int     `orm:"column(tid)"`
	Score        float64 `orm:"column(score)"`
	Year         int     `orm:"column(year)"`
	Quarter      int     `orm:"column(quarter)"`
	DepartmentID int     `orm:"column(department_id)"`
}

// 添加部门负责人季度评分记录
func AddDepartmentLeaderScoreRecord(d *DepartmentLeaderScoreRecord) error {
	_, err := orm.NewOrm().Insert(d)
	return err
}

// 更新部门负责人季度评分记录
func (d *DepartmentLeaderScoreRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(d, fields...)
	return err
}

// 根据条件查询部门负责人季度评分记录
func SearchDepartmentLeaderScoreRecordsByFilters(filters ...DBFilter) []*DepartmentLeaderScoreRecord {
	records := make([]*DepartmentLeaderScoreRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("department_leader_score_record"))

	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}

	query.OrderBy("-id").All(&records)
	return records
}
