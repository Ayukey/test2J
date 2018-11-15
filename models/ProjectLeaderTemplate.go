package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目负责人评分模版
type ProjectLeaderTemplate struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	ScoreLimit float64 `orm:"column(score_limit)"`
	Status     int     `orm:"column(status);default(1)"`
}

// 根据ID搜索项目负责人评分模版
func SearchProjectLeaderTemplateByID(id int) (ProjectLeaderTemplate, error) {
	var template ProjectLeaderTemplate
	err := orm.NewOrm().QueryTable(TableName("project_leader_template")).Filter("id", id).One(&template)
	return template, err
}

// 更新项目负责人评分模版
func (p *ProjectLeaderTemplate) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

// 搜索所有的项目评分模版
func SearchAllProjectLeaderTemplates() []*ProjectLeaderTemplate {
	templates := make([]*ProjectLeaderTemplate, 0)
	query := orm.NewOrm().QueryTable(TableName("project_leader_template"))
	query.OrderBy("id").All(&templates)
	return templates
}
