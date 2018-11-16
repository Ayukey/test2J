package models

import "github.com/astaxie/beego/orm"

type QuarterActiveDepartmentLeader struct {
	ID           int `orm:"column(id)"`
	UID          int `orm:"column(uid)"`
	Year         int `orm:"column(year)"`
	Quarter      int `orm:"column(quarter)"`
	DepartmentID int `orm:"column(department_id)"`
}

func AddQuarterActiveDepartmentLeader(year, quarter, uid, departmentId int) error {
	o := orm.NewOrm()
	exist := o.QueryTable("quarter_active_department_leader").Filter("year", year).Filter("quarter", quarter).Filter("uid", uid).Filter("department_id", departmentId).Exist()
	if !exist {
		q := new(QuarterActiveDepartmentLeader)
		q.Year = year
		q.Quarter = quarter
		q.UID = uid
		q.DepartmentID = departmentId
		_, err := o.Insert(q)
		return err
	}
	return nil
}

func DeleteQuarterActiveDepartmentLeader(year, quarter, uid, departmentId int) error {
	o := orm.NewOrm()
	err := o.Begin()

	_, err = o.QueryTable(TableName("quarter_active_department_leader")).Filter("year", year).Filter("quarter", quarter).Filter("uid", uid).Filter("department_id", departmentId).Delete()

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}

func SearchAllActiveQuarterDepartmentLeaders(year, quarter int) []*QuarterActiveDepartmentLeader {
	leaders := make([]*QuarterActiveDepartmentLeader, 0)
	orm.NewOrm().QueryTable(TableName("quarter_active_department_leader")).Filter("year", year).Filter("quarter", quarter).All(&leaders)
	return leaders
}

func ActiveQuarterAllDepartmentLeaders(year, quarter int) error {
	o := orm.NewOrm()
	err := o.Begin()
	currentDepartmentLeaders := make([]*QuarterActiveDepartmentLeader, 0)
	count, err := o.QueryTable(TableName("quarter_active_department_leader")).Filter("year", year).Filter("quarter", quarter).All(&currentDepartmentLeaders)

	if count == 0 {
		leaders := SearchAllDepartmentLeadersInDepartment()
		departmentLeaders := make([]QuarterActiveDepartmentLeader, 0)

		for _, leader := range leaders {
			departmentLeader := QuarterActiveDepartmentLeader{
				UID:          leader.User.ID,
				DepartmentID: leader.Department.ID,
				Year:         year,
				Quarter:      quarter,
			}
			departmentLeaders = append(departmentLeaders, departmentLeader)
		}

		_, err = o.InsertMulti(len(departmentLeaders), departmentLeaders)
	}

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}
