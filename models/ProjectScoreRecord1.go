package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目一级评分记录
type ProjectScoreRecord1 struct {
	ID         int     `orm:"column(id)"`
	PID        int     `orm:"column(pid)"`
	T1ID       int     `orm:"column(t1id)"`
	Year       int     `orm:"column(year)"`
	Quarter    int     `orm:"column(quarter)"`
	TotalScore float64 `orm:"column(total_score)"`
	UpdateDate float64 `orm:"column(update_date)"`
}

// 添加项目一级评分记录
func AddProjectScoreRecord1(s *ProjectScoreRecord1) error {
	_, err := orm.NewOrm().Insert(s)
	return err
}

// 更新项目一级评分记录
func (s *ProjectScoreRecord1) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}

// 查询项目一级评分记录
func SearchProjectScoreRecord1sByFilters(filters ...DBFilter) []*ProjectScoreRecord1 {
	records := make([]*ProjectScoreRecord1, 0)
	query := orm.NewOrm().QueryTable(TableName("project_score_record1"))

	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}

	query.OrderBy("-id").All(&records)
	return records
}
