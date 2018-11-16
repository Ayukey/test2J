package models

import "github.com/astaxie/beego/orm"

type QuarterActive struct {
	ID      int `orm:"column(id)"`
	Year    int `orm:"column(year)"`
	Quarter int `orm:"column(quarter)"`
	Status  int `orm:"column(status)"`
}

// 查询当前的评分季度
func SearchQuarterInActive() (QuarterActive, error) {
	var quarterActive QuarterActive
	err := orm.NewOrm().QueryTable(TableName("quarter_active")).Filter("status", 1).One(&quarterActive)
	return quarterActive, err
}

func UnActiveQuarter(year, quarter int) error {
	_, err := orm.NewOrm().QueryTable(TableName("quarter_active")).Filter("year", year).Filter("quarter", quarter).Update(orm.Params{
		"status": 0,
	})
	return err
}

func ActiveQuarter(year, quarter int) error {
	o := orm.NewOrm()
	err := o.Begin()
	_, err = o.QueryTable(TableName("quarter_active")).Filter("status", 1).Update(orm.Params{
		"status": 0,
	})
	exist := o.QueryTable(TableName("quarter_active")).Filter("year", year).Filter("quarter", quarter).Exist()

	if exist {
		_, err = o.QueryTable(TableName("quarter_active")).Filter("year", year).Filter("quarter", quarter).Update(orm.Params{
			"status": 1,
		})
	} else {
		q := new(QuarterActive)
		q.Year = year
		q.Quarter = quarter
		q.Status = 1
		_, err = o.Insert(q)
	}

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}
