package models

import "github.com/astaxie/beego/orm"

type QuarterActiveProjectLeader struct {
	ID        int `orm:"column(id)"`
	UID       int `orm:"column(uid)"`
	ProjectID int `orm:"column(project_id)"`
	Year      int `orm:"column(year)"`
	Quarter   int `orm:"column(quarter)"`
}

func AddQuarterActiveProjectLeader(year, quarter, uid, projectId int) error {
	o := orm.NewOrm()
	exist := o.QueryTable(TableName("quarter_active_project_leader")).Filter("year", year).Filter("quarter", quarter).Filter("uid", uid).Filter("project_id", projectId).Exist()
	if !exist {
		q := new(QuarterActiveProjectLeader)
		q.Year = year
		q.Quarter = quarter
		q.UID = uid
		q.ProjectID = projectId
		_, err := o.Insert(q)
		return err
	}
	return nil
}

func DeleteQuarterActiveProjectLeader(year, quarter, uid, projectId int) error {
	o := orm.NewOrm()
	err := o.Begin()

	_, err = o.QueryTable(TableName("quarter_active_project_leader")).Filter("year", year).Filter("quarter", quarter).Filter("uid", uid).Filter("project_id", projectId).Delete()

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}

func SearchAllActiveQuarterProjectLeaders(year, quarter int) []*QuarterActiveProjectLeader {
	leaders := make([]*QuarterActiveProjectLeader, 0)
	orm.NewOrm().QueryTable(TableName("quarter_active_project_leader")).Filter("year", year).Filter("quarter", quarter).All(&leaders)
	return leaders
}

func ActiveQuarterAllProjectLeaders(year, quarter int) error {
	o := orm.NewOrm()
	err := o.Begin()
	currentProjectLeaders := make([]*QuarterActiveProjectLeader, 0)
	count, err := o.QueryTable(TableName("quarter_active_project_leader")).Filter("year", year).Filter("quarter", quarter).All(&currentProjectLeaders)

	if count == 0 {
		leaders := SearchAllProjectLeadersInProject()
		projectLeaders := make([]QuarterActiveProjectLeader, 0)

		for _, leader := range leaders {
			projectLeader := QuarterActiveProjectLeader{
				UID:       leader.User.ID,
				ProjectID: leader.Project.ID,
				Year:      year,
				Quarter:   quarter,
			}
			projectLeaders = append(projectLeaders, projectLeader)
		}

		_, err = o.InsertMulti(len(projectLeaders), projectLeaders)
	}

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}
