package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目一级模版信息
type ProjectTemplate1 struct {
	ID     int    `orm:"column(id)"`
	Name   string `orm:"column(name)"`
	Status int    `orm:"column(status);default(1)"`
}

// 添加项目一级模版信息
func AddProjectTemplate1(s *ProjectTemplate1) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(s, "Name")
	return err
}

// 根据ID查询项目一级模版信息
func SearchProjectTemplate1ByID(id int) (ProjectTemplate1, error) {
	var template ProjectTemplate1
	err := orm.NewOrm().QueryTable(TableName("project_template1")).Filter("id", id).One(&template)
	return template, err
}

// 查询所有项目一级模版信息
func SearchAllProjectTemplate1s() []*ProjectTemplate1 {
	templates := make([]*ProjectTemplate1, 0)
	query := orm.NewOrm().QueryTable(TableName("project_template1")).Filter("status", 1)
	query.OrderBy("id").All(&templates)
	return templates
}

// 更新项目一级模版信息
func (s *ProjectTemplate1) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}
