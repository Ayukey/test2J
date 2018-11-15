package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目信息
type Project struct {
	ID           int    `orm:"column(id)"`
	Name         string `orm:"column(name);index;unique"`
	Active       int    `orm:"column(active);default(1)"`
	Leader       int    `orm:"column(leader);null"`
	LeaderActive int    `orm:"column(leader_active);default(1)"`
	Status       int    `orm:"column(status);default(1)"`
}

// 所有项目信息
func SearchAllProjects() []*Project {
	projects := make([]*Project, 0)
	query := orm.NewOrm().QueryTable(TableName("project")).Filter("status", 1)
	query.OrderBy("id").All(&projects)
	return projects
}

// 新增项目信息
func AddProject(p *Project) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(p, "Name")
	return err
}

// 根据ID搜索项目信息
func SearchProjectByID(id int) (Project, error) {
	var project Project
	err := orm.NewOrm().QueryTable(TableName("project")).Filter("id", id).One(&project)
	return project, err
}

// 更新项目信息
func (p *Project) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(p, fields...)
	return err
}

//  所有参与考核的项目
func SearchAllProjectsInActive() []*Project {
	projects := make([]*Project, 0)
	query := orm.NewOrm().QueryTable(TableName("project")).Filter("status", 1).Filter("active", 1)
	query.OrderBy("id").All(&projects)
	return projects
}

type ProjectLeader struct {
	Project *Project
	User    *User
}

//  所有参与考核的项目负责人
func SearchAllProjectLeadersInActive() []*ProjectLeader {
	leaders := make([]*ProjectLeader, 0)

	projects := make([]*Project, 0)
	orm.NewOrm().QueryTable(TableName("project")).Filter("status", 1).Filter("leader_active", 1).OrderBy("id").All(&projects)

	for _, project := range projects {
		leader := new(ProjectLeader)
		leader.Project = project
		user, err := SearchUserByID(project.Leader)
		if err == nil {
			leader.User = &user
			leaders = append(leaders, leader)
		}
	}

	return leaders
}
