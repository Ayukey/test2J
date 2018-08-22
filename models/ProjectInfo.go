package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

// 项目信息
type ProjectInfo struct {
	ID        int     `orm:"column(id)"`
	Name      string  `orm:"column(name)"`
	UserID    int     `orm:"column(user_id)"`
	BeginDate float64 `orm:"column(begin_date)"`
	EndDate   float64 `orm:"column(end_date)"`
	Status    int     `orm:"column(status)"`
}

func (p *ProjectInfo) TableName() string {
	return TableName("project_info")
}

func SearchProjectInfoList(page, pageSize int, filters ...interface{}) ([]*ProjectInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ProjectInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("project_info"))
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

func SearchAllProjectInfo() []*ProjectInfo {
	list := make([]*ProjectInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("project_info"))
	query.OrderBy("id").All(&list)
	return list
}

func AddProjectInfo(p *ProjectInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchProjectInfoByID(id int) (*ProjectInfo, error) {
	p := new(ProjectInfo)
	err := orm.NewOrm().QueryTable(TableName("project_info")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func SearchProjectInfoByName(name string) (*ProjectInfo, error) {
	p := new(ProjectInfo)
	err := orm.NewOrm().QueryTable(TableName("project_info")).Filter("name", name).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *ProjectInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func SearchAllProjectsInUse(year, quarter int) []*ProjectInfo {
	inProjectList := make([]*ProjectInfo, 0)
	list := make([]*ProjectInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("project_info"))
	query.OrderBy("id").All(&list)

	for _, p := range list {
		tmBegin := time.Unix(int64(p.BeginDate/1000), 0)
		tmEnd := time.Unix(int64(p.EndDate/1000), 0)

		fmt.Printf("项目名称 %s\n开始时间 年：%d月：%d\n结束时间 年：%d月：%d\n", p.Name, tmBegin.Year(), tmBegin.Month(), tmEnd.Year(), tmEnd.Month())

		if year >= tmBegin.Year() && year <= tmEnd.Year() {
			switch quarter {
			case 1:
				if (tmBegin.Month() >= 1 && tmBegin.Month() <= 3) || (tmBegin.Month() >= 1 && tmBegin.Month() <= 3) {
					inProjectList = append(inProjectList, p)
				}
			case 2:
				if (tmBegin.Month() >= 4 && tmBegin.Month() <= 6) || (tmBegin.Month() >= 4 && tmBegin.Month() <= 6) {
					inProjectList = append(inProjectList, p)
				}
			case 3:
				if (tmBegin.Month() >= 7 && tmBegin.Month() <= 9) || (tmBegin.Month() >= 7 && tmBegin.Month() <= 9) {
					inProjectList = append(inProjectList, p)
				}
			case 4:
				if (tmBegin.Month() >= 10 && tmBegin.Month() <= 12) || (tmBegin.Month() >= 10 && tmBegin.Month() <= 12) {
					inProjectList = append(inProjectList, p)
				}
			}
		}
	}
	return inProjectList
}
