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

	for _, p := range list {
		timeEnd := time.Unix(int64(p.EndDate/1000), 0)
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

		fmt.Printf("项目名称: %s\n 结束时间所在季度最后一天: %d\n", p.Name, timeQuarterLastMonthTS)

		if currentQuarterLastMonthTS <= timeQuarterLastMonthTS {
			inProjectList = append(inProjectList, p)
		}
	}
	return inProjectList
}
