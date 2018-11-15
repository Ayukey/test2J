package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门负责人互评（单个项目负责人季度评分）总分记录
type DepartmentLeaderSumScoreRecord struct {
	ID           int     `orm:"column(id)"`
	SUID         int     `orm:"column(suid)"`
	UID          int     `orm:"column(uid)"`
	Score        float64 `orm:"column(score)"`
	Year         int     `orm:"column(year)"`
	Quarter      int     `orm:"column(quarter)"`
	DepartmentID int     `orm:"column(department_id)"`
}

// 添加部门负责人评分总分记录（基于单个项目评分人）
func AddDepartmentLeaderSumScoreRecord(p *DepartmentLeaderSumScoreRecord) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 根据ID搜索部门负责人评分总分记录（基于单个项目评分人）
func SearchDepartmentLeaderSumScoreRecordByID(id int) (DepartmentLeaderSumScoreRecord, error) {
	var record DepartmentLeaderSumScoreRecord
	err := orm.NewOrm().QueryTable(TableName("department_leader_sum_score_record")).Filter("id", id).One(&record)
	return record, err
}

// 更新部门负责人评分总分记录（基于单个项目评分人）
func (p *DepartmentLeaderSumScoreRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 根据条件搜索部门负责人评分总分记录（基于单个项目评分人）
func SearchDepartmentLeaderSumScoreRecordsByFilters(filters ...DBFilter) []*DepartmentLeaderSumScoreRecord {
	records := make([]*DepartmentLeaderSumScoreRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("department_leader_sum_score_record"))

	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}

	query.OrderBy("-id").All(&records)
	return records
}
