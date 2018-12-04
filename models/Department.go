package models

import (
	"github.com/astaxie/beego/orm"
)

// 部门信息
type Department struct {
	ID     int    `orm:"column(id)"`
	Name   string `orm:"column(name);index;unique"`
	Leader int    `orm:"column(leader);null"`
	Status int    `orm:"column(status);default(1)"`
}

// 所有部门
func SearchAllDepartments() []*Department {
	departments := make([]*Department, 0)
	query := orm.NewOrm().QueryTable(TableName("department")).Filter("status", 1)
	query.OrderBy("id").All(&departments)
	return departments
}

// 根据部门负责人搜索部门信息
func SearchDepartmentsByLeader(leader int) []*Department {
	departments := make([]*Department, 0)
	query := orm.NewOrm().QueryTable(TableName("department")).Filter("status", 1).Filter("leader", leader)
	query.OrderBy("id").All(&departments)
	return departments
}

// 新增部门
func AddDepartment(d *Department) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(d, "Name")
	return err
}

// 根据ID搜索部门
func SearchDepartmentByID(id int) (Department, error) {
	var department Department
	err := orm.NewOrm().QueryTable(TableName("department")).Filter("id", id).One(&department)
	return department, err
}

// 更新部门
func (d *Department) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(d, fields...)
	return err
}

type DepartmentLeader struct {
	Department *Department
	User       *User
}

//  所有参与考核的部门负责人
func SearchAllDepartmentLeadersInDepartment() []*DepartmentLeader {
	leaders := make([]*DepartmentLeader, 0)

	departments := make([]*Department, 0)
	orm.NewOrm().QueryTable(TableName("department")).Filter("status", 1).OrderBy("id").All(&departments)

	for _, department := range departments {
		leader := new(DepartmentLeader)
		leader.Department = department
		user, err := SearchUserByID(department.Leader)
		if err == nil {
			leader.User = &user
			leaders = append(leaders, leader)
		}
	}

	return leaders
}

//  根据季度获取该季度下有效的部门负责人
func SearchAllDepartmentLeadersInActive(year, quarter int) []*DepartmentLeader {
	leaders := make([]*DepartmentLeader, 0)

	activeLeaders := SearchAllActiveQuarterDepartmentLeaders(year, quarter)

	for _, activeLeader := range activeLeaders {
		leader := new(DepartmentLeader)
		department, _ := SearchDepartmentByID(activeLeader.DepartmentID)
		user, _ := SearchUserByID(activeLeader.UID)

		if &department != nil && &user != nil {
			leader.Department = &department
			leader.User = &user
			leaders = append(leaders, leader)
		}
	}
	return leaders
}
