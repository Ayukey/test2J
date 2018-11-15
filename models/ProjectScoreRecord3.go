package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目三级评分记录
type ProjectScoreRecord3 struct {
	ID         int     `orm:"column(id)"`
	UID        int     `orm:"column(uid)"`
	PID        int     `orm:"column(pid)"`
	T1ID       int     `orm:"column(t1id)"`
	T2ID       int     `orm:"column(t2id)"`
	T3ID       int     `orm:"column(t3id)"`
	Score      float64 `orm:"column(score)"`
	Year       int     `orm:"column(year)"`
	Quarter    int     `orm:"column(quarter)"`
	UpdateDate float64 `orm:"column(update_date)"`
}

// 添加项目三级评分记录
func AddProjectScoreRecord3(s *ProjectScoreRecord3) error {
	_, err := orm.NewOrm().Insert(s)
	return err
}

// 更新项目三级评分记录
func (s *ProjectScoreRecord3) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}

// 查询项目三级评分记录
func SearchProjectScoreRecord3sByFilters(filters ...DBFilter) []*ProjectScoreRecord3 {
	records := make([]*ProjectScoreRecord3, 0)
	query := orm.NewOrm().QueryTable(TableName("project_score_record3"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}

// 根据ID查询项目三级评分记录
func SearchProjectScoreRecord3ByID(id int) (ProjectScoreRecord3, error) {
	var record ProjectScoreRecord3
	err := orm.NewOrm().QueryTable(TableName("project_score_record3")).Filter("id", id).One(&record)
	return record, err
}
