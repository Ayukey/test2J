package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

// 用户
type UserInfo struct {
	ID           int    `orm:"column(id)"`
	Account      string `orm:"column(account)"`
	Password     string `orm:"column(password)"`
	Name         string `orm:"column(name)"`
	Phone        string `orm:"column(phone)"`
	RoleID       int    `orm:"column(role_id)"`
	DepartmentID int    `orm:"column(department_id)"`
	ProjectID    int    `orm:"column(project_id)"`
	Status       int    `orm:"column(status)"`
}

func (u *UserInfo) TableName() string {
	return TableName("user_info")
}

func SearchUserInfoList(page, pageSize int, filters ...interface{}) ([]*UserInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

func SearchAllUserInfoList(filters ...interface{}) []*UserInfo {
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("id").All(&list)
	return list
}

func AddUserInfo(u *UserInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(u)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchUserInfoByID(id int) (*UserInfo, error) {
	u := new(UserInfo)
	err := orm.NewOrm().QueryTable(TableName("user_info")).Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func SearchUserInfoByAccount(account string) (*UserInfo, error) {
	u := new(UserInfo)
	err := orm.NewOrm().QueryTable(TableName("user_info")).Filter("account", account).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

// 搜索所有部门负责人
func SearchAllDepartmentLeaders() []*UserInfo {
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("role_id", 3)
	query.OrderBy("id").All(&list)
	return list
}

// 根据部门ID搜索所有部门成员
func SearchAllDepartmentMembers(departmentID int) []*UserInfo {
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("department_id", departmentID)
	query.OrderBy("id").All(&list)
	return list
}

// 搜索所有项目负责人
func SearchAllProjectLeaders() []*UserInfo {
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("role_id", 5)
	query.OrderBy("id").All(&list)
	return list
}

// 根据部门ID搜索所有部门成员
func SearchAllProjectMembers(projectID int) []*UserInfo {
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("project_id", projectID)
	query.OrderBy("id").All(&list)
	return list
}

// 搜索所有有效项目负责人
func SearchAllProjectLeadersInUse(year, quarter int) []*UserInfo {
	inUserList := make([]*UserInfo, 0)
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("role_id", 5)
	query.OrderBy("id").All(&list)

	for _, user := range list {
		projectorInfo, _ := SearchProjectorInfoByUserID(user.ID)
		tmBegin := time.Unix(int64(projectorInfo.BeginDate/1000), 0)
		tmEnd := time.Unix(int64(projectorInfo.EndDate/1000), 0)

		fmt.Printf("项目负责人 %s\n开始时间 年：%d月：%d\n结束时间 年：%d月：%d\n", user.Name, tmBegin.Year(), tmBegin.Month(), tmEnd.Year(), tmEnd.Month())

		if year >= tmBegin.Year() && year <= tmEnd.Year() {
			switch quarter {
			case 1:
				if (tmBegin.Month() >= 1 && tmBegin.Month() <= 3) || (tmBegin.Month() >= 1 && tmBegin.Month() <= 3) {
					inUserList = append(inUserList, user)
				}
			case 2:
				if (tmBegin.Month() >= 4 && tmBegin.Month() <= 6) || (tmBegin.Month() >= 4 && tmBegin.Month() <= 6) {
					inUserList = append(inUserList, user)
				}
			case 3:
				if (tmBegin.Month() >= 7 && tmBegin.Month() <= 9) || (tmBegin.Month() >= 7 && tmBegin.Month() <= 9) {
					inUserList = append(inUserList, user)
				}
			case 4:
				if (tmBegin.Month() >= 10 && tmBegin.Month() <= 12) || (tmBegin.Month() >= 10 && tmBegin.Month() <= 12) {
					inUserList = append(inUserList, user)
				}
			}
		}
	}
	return inUserList
}

// 项目负责人
type Projector struct {
	UserInfo *UserInfo
	Project  *ProjectInfo
	Score    *ProjectorScoreSumInfo
}

// 搜索所有有效项目负责人
func SearchAllProjectLeadersInUseWithScore(year, quarter, uid int) []*Projector {
	inUserList := make([]*Projector, 0)
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("role_id", 5)
	query.OrderBy("id").All(&list)

	now := time.Now()
	currentYear := year
	currentLocation := now.Location()
	currentMonth := time.Month(3)
	switch quarter {
	case 1:
		currentMonth = time.Month(3)
	case 2:
		currentMonth = time.Month(6)
	case 3:
		currentMonth = time.Month(9)
	case 4:
		currentMonth = time.Month(12)
	}

	currentQuarterFirstMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentQuarterLastMonth := currentQuarterFirstMonth.AddDate(0, 1, -1)
	currentQuarterLastMonthTS := currentQuarterLastMonth.Unix()

	fmt.Printf("当前季度时间戳: %d\n", currentQuarterLastMonthTS)

	for _, user := range list {
		project, _ := SearchProjectInfoByID(user.ProjectID)
		projectorInfo, _ := SearchProjectorInfoByUserID(user.ID)

		if project != nil && projectorInfo != nil {
			timeEnd := time.Unix(int64(projectorInfo.EndDate/1000), 0)
			timeEndMonth := time.Month(3)
			if 1 <= timeEnd.Month() && timeEnd.Month() <= 3 {
				timeEndMonth = time.Month(3)
			} else if 4 <= timeEnd.Month() && timeEnd.Month() <= 6 {
				timeEndMonth = time.Month(6)
			} else if 7 <= timeEnd.Month() && timeEnd.Month() <= 9 {
				timeEndMonth = time.Month(9)
			} else if 10 <= timeEnd.Month() && timeEnd.Month() <= 12 {
				timeEndMonth = time.Month(12)
			}
			timeQuarterFirstMonth := time.Date(timeEnd.Year(), timeEndMonth, 1, 0, 0, 0, 0, currentLocation)
			timeQuarterLastMonth := timeQuarterFirstMonth.AddDate(0, 1, -1)
			timeQuarterLastMonthTS := timeQuarterLastMonth.Unix()

			projector := new(Projector)
			projector.UserInfo = user
			projector.Project = project

			if currentQuarterLastMonthTS <= timeQuarterLastMonthTS {
				filters := make([]interface{}, 0)
				filters = append(filters, "year", year)
				filters = append(filters, "quarter", quarter)
				filters = append(filters, "user_id", user.ID)
				filters = append(filters, "project_id", project.ID)
				filters = append(filters, "scoreuser_id", uid)
				recordList := SearchProjectorScoreSumInfoByFilters(filters...)
				if len(recordList) == 1 {
					projector.Score = recordList[0]
				}
				inUserList = append(inUserList, projector)
			}

		}
	}
	return inUserList
}

// 项目负责人
type Departmentor struct {
	UserInfo   *UserInfo
	Department *DepartmentInfo
	Score      *DepartmentorScoreSumInfo
}

// 搜索所有有效项目负责人
func SearchAllDepartmentLeadersInUseWithScore(year, quarter, uid int) []*Departmentor {
	inUserList := make([]*Departmentor, 0)
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("user_info")).Filter("role_id", 3)
	query.OrderBy("id").All(&list)

	for _, u := range list {
		d, _ := SearchDepartmentInfoByID(u.DepartmentID)

		departmentor := new(Departmentor)
		departmentor.UserInfo = u
		departmentor.Department = d

		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "user_id", u.ID)
		filters = append(filters, "department_id", d.ID)
		filters = append(filters, "scoreuser_id", uid)
		recordList := SearchDepartmentorScoreSumInfoByFilters(filters...)
		if len(recordList) == 1 {
			departmentor.Score = recordList[0]
		}
		inUserList = append(inUserList, departmentor)
	}
	return inUserList
}
