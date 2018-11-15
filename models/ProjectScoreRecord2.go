package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目二级评分记录
type ProjectScoreRecord2 struct {
	ID         int     `orm:"column(id)"`
	PID        int     `orm:"column(pid)"`
	T1ID       int     `orm:"column(t1id)"`
	T2ID       int     `orm:"column(t2id)"`
	Year       int     `orm:"column(Year)"`
	Quarter    int     `orm:"column(Quarter)"`
	TotalScore float64 `orm:"column(total_score)"`
	UpdateDate float64 `orm:"column(update_date)"`
	Remark     string  `orm:"column(remark);null"`
}

// 添加项目二级评分记录
func AddProjectScoreRecord2(s *ProjectScoreRecord2) error {
	_, err := orm.NewOrm().Insert(s)
	return err
}

// 更新项目二级评分记录
func (s *ProjectScoreRecord2) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}

// 查询项目二级评分记录
func SearchProjectScoreRecord2sByFilters(filters ...DBFilter) []*ProjectScoreRecord2 {
	records := make([]*ProjectScoreRecord2, 0)
	query := orm.NewOrm().QueryTable(TableName("project_score_record2"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}
