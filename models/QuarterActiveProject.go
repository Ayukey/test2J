package models

import "github.com/astaxie/beego/orm"

type QuarterActiveProject struct {
	ID      int `orm:"column(id)"`
	PID     int `orm:"column(pid)"`
	Year    int `orm:"column(year)"`
	Quarter int `orm:"column(quarter)"`
}

func AddQuarterActiveProject(year, quarter, projectId int) error {
	o := orm.NewOrm()
	exist := o.QueryTable(TableName("quarter_active_project")).Filter("year", year).Filter("quarter", quarter).Filter("pid", projectId).Exist()
	if !exist {
		q := new(QuarterActiveProject)
		q.Year = year
		q.Quarter = quarter
		q.PID = projectId
		_, err := o.Insert(q)
		return err
	}
	return nil
}

func DeleteQuarterActiveProject(year, quarter, projectId int) error {
	o := orm.NewOrm()
	err := o.Begin()

	_, err = o.QueryTable(TableName("quarter_active_project")).Filter("year", year).Filter("quarter", quarter).Filter("pid", projectId).Delete()

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}

func ExistActiveQuarterProject(year, quarter, projectId int) bool {
	o := orm.NewOrm()
	isExist := o.QueryTable(TableName("quarter_active_project")).Filter("year", year).Filter("quarter", quarter).Filter("pid", projectId).Exist()
	return isExist
}

func SearchAllActiveQuarterProjects(year, quarter int) []*QuarterActiveProject {
	projects := make([]*QuarterActiveProject, 0)
	orm.NewOrm().QueryTable(TableName("quarter_active_project")).Filter("year", year).Filter("quarter", quarter).All(&projects)
	return projects
}

func ActiveQuarterAllProjects(year, quarter int) error {
	o := orm.NewOrm()
	err := o.Begin()
	currentProjects := make([]*QuarterActiveProject, 0)
	count, err := o.QueryTable(TableName("quarter_active_project")).Filter("year", year).Filter("quarter", quarter).All(&currentProjects)

	if count == 0 {
		projects := SearchAllProjects()
		activeProjects := make([]QuarterActiveProject, 0)

		for _, project := range projects {
			activeProject := QuarterActiveProject{
				PID:     project.ID,
				Year:    year,
				Quarter: quarter,
			}
			activeProjects = append(activeProjects, activeProject)
		}

		_, err = o.InsertMulti(len(activeProjects), activeProjects)
	}

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}
