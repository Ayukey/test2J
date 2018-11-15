package models

import "github.com/astaxie/beego/orm"

// 项目负责人互评（单个部门负责人季度评分）总分记录
type ProjectLeaderSumScoreRecord struct {
	ID        int     `orm:"column(id)"`
	SUID      int     `orm:"column(suid)"`
	UID       int     `orm:"column(uid)"`
	Score     float64 `orm:"column(score)"`
	Year      int     `orm:"column(year)"`
	Quarter   int     `orm:"column(quarter)"`
	ProjectID int     `orm:"column(project_id)"`
}

// 添加项目负责人评分总分记录（基于单个部门评分人）
func AddProjectLeaderSumScoreRecord(p *ProjectLeaderSumScoreRecord) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 根据ID搜索项目负责人评分总分记录（基于单个部门评分人）
func SearchProjectLeaderSumScoreRecordByID(id int) (ProjectLeaderSumScoreRecord, error) {
	var record ProjectLeaderSumScoreRecord
	err := orm.NewOrm().QueryTable(TableName("project_leader_sum_score_record")).Filter("id", id).One(&record)
	return record, err
}

// 更新项目负责人评分总分记录（基于单个部门评分人）
func (p *ProjectLeaderSumScoreRecord) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 根据条件搜索项目负责人评分总分记录（基于单个部门评分人）
func SearchProjectLeaderSumScoreRecordsByFilters(filters ...DBFilter) []*ProjectLeaderSumScoreRecord {
	records := make([]*ProjectLeaderSumScoreRecord, 0)
	query := orm.NewOrm().QueryTable(TableName("project_leader_sum_score_record"))

	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}

	query.OrderBy("-id").All(&records)
	return records
}
