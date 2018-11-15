package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目负责人评分记录
type ProjectLeaderScoreRecord struct {
	ID        int     `orm:"column(id)"`
	SUID      int     `orm:"column(suid)"` // 评分者ID
	UID       int     `orm:"column(uid)"`  // 被评分者ID
	TID       int     `orm:"column(tid)"`
	Score     float64 `orm:"column(score)"`
	Year      int     `orm:"column(year)"`
	Quarter   int     `orm:"column(quarter)"`
	ProjectID int     `orm:"column(project_id)"`
}

// 添加项目负责人季度评分记录
func AddProjectLeaderScoreRecord(p *ProjectLeaderScoreRecord) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 更新项目负责人季度评分记录
func (p *ProjectLeaderScoreRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 根据条件查询项目负责人季度评分记录
func SearchProjectLeaderScoreRecordsByFilters(filters ...DBFilter) []*ProjectLeaderScoreRecord {
	records := make([]*ProjectLeaderScoreRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("project_leader_score_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}
