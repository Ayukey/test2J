package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门负责人评分模版
type DepartmentLeaderTemplate struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	ScoreLimit float64 `orm:"column(score_limit)"`
	Status     int     `orm:"column(status);default(1)"`
}

// 根据ID搜索部门负责人评分模版
func SearchDepartmentLeaderTemplateByID(id int) (DepartmentLeaderTemplate, error) {
	var template DepartmentLeaderTemplate
	err := orm.NewOrm().QueryTable(TableName("department_leader_template")).Filter("id", id).One(&template)
	return template, err
}

// 更新部门负责人评分模版
func (d *DepartmentLeaderTemplate) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(d, fields...)
	return err
}

// 所有部门负责人评分模版
func SearchAllDepartmentLeaderTemplates() []*DepartmentLeaderTemplate {
	templates := make([]*DepartmentLeaderTemplate, 0)
	query := orm.NewOrm().QueryTable(TableName("department_leader_template")).Filter("status", 1)
	query.OrderBy("id").All(&templates)
	return templates
}
