package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目二级模版信息
type ProjectTemplate2 struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	Percentage float64 `orm:"column(percentage)"`
	TID        int     `orm:"column(tid)"`
	Status     int     `orm:"column(status);default(1)"`
}

// 添加项目二级模版信息
func AddProjectTemplate2(s *ProjectTemplate2) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(s, "Name")
	return err
}

// 根据ID查询项目二级模版信息
func SearchProjectTemplate2ByID(id int) (ProjectTemplate2, error) {
	var template ProjectTemplate2
	err := orm.NewOrm().QueryTable(TableName("project_template2")).Filter("id", id).One(&template)
	return template, err
}

// 根据TID(项目一级模版信息ID)查询所有项目二级模版信息
func SearchProjectTemplate2sByTID(tid int) []*ProjectTemplate2 {
	templates := make([]*ProjectTemplate2, 0)
	orm.NewOrm().QueryTable(TableName("project_template2")).Filter("tid", tid).Filter("status", 1).All(&templates)
	return templates
}

// 更新项目二级模版信息
func (s *ProjectTemplate2) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}
