package models

import (
	"github.com/astaxie/beego/orm"
)

// 已发布的项目评分记录
type ProjectReleaseRecord struct {
	ID      int `orm:"column(id)"`
	PID     int `orm:"column(pid)"`
	Year    int `orm:"column(year)"`
	Quarter int `orm:"column(quarter)"`
}

// 新增项目发布评分记录
func AddProjectReleaseRecord(p *ProjectReleaseRecord) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 更新项目发布评分记录
func (p *ProjectReleaseRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 根据条件查询项目发布评分记录
func SearchProjectReleaseRecordsByFilters(filters ...DBFilter) []*ProjectReleaseRecord {
	records := make([]*ProjectReleaseRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("project_release_record"))
	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}
	query.OrderBy("-id").All(&records)
	return records
}
