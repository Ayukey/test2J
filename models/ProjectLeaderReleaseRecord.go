package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目负责人评分发布记录
type ProjectLeaderReleaseRecord struct {
	ID        int     `orm:"column(id)"`
	UID       int     `orm:"column(uid)"`
	Score     float64 `orm:"column(score)"`
	Year      int     `orm:"column(year)"`
	Quarter   int     `orm:"column(quarter)"`
	ProjectID int     `orm:"column(project_id)"`
}

// 新增项目负责人评分发布记录
func AddProjectLeaderReleaseRecord(p *ProjectLeaderReleaseRecord) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 更新部门负责人评分发布记录
func (p *ProjectLeaderReleaseRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 根据条件查询项目负责人季度评分发布记录
func SearchProjectLeaderReleaseRecordsByFilters(filters ...DBFilter) []*ProjectLeaderReleaseRecord {
	records := make([]*ProjectLeaderReleaseRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("project_leader_release_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}

// 根据条件查询项目负责人季度评分发布记录(按分数排序)
func SearchProjectLeaderReleaseRecordsByOrder(filters ...DBFilter) []*ProjectLeaderReleaseRecord {
	records := make([]*ProjectLeaderReleaseRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("project_leader_release_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("score").All(&records)
	return records
}
