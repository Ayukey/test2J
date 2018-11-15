package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目三级模版信息
type ProjectTemplate3 struct {
	ID       int     `orm:"column(id)"`
	Name     string  `orm:"column(name)"`
	MaxScore float64 `orm:"column(max_score)"`
	TID      int     `orm:"column(tid)"`
	Status   int     `orm:"column(status);default(1)"`
}

// 添加项目三级模版信息
func AddProjectTemplate3(s *ProjectTemplate3) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(s, "Name")
	return err
}

// 根据ID查询项目三级模版信息
func SearchProjectTemplate3ByID(id int) (ProjectTemplate3, error) {
	var template ProjectTemplate3
	err := orm.NewOrm().QueryTable(TableName("project_template3")).Filter("id", id).One(&template)
	return template, err
}

// 根据TID(项目二级模版信息ID)查询所有项目三级模版信息
func SearchProjectTemplate3sByTID(tid int) []*ProjectTemplate3 {
	template := make([]*ProjectTemplate3, 0)
	orm.NewOrm().QueryTable(TableName("project_template3")).Filter("tid", tid).All(&template)
	return template
}

// 更新项目三级模版信息
func (s *ProjectTemplate3) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(s, fields...)
	return err
}
