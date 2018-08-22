package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目一级模版信息
type ScoreTypeInfoI struct {
	ID     int    `orm:"column(id)"`
	Name   string `orm:"column(name)"`
	Status int    `orm:"column(status)"`
}

func (s *ScoreTypeInfoI) TableName() string {
	return TableName("score_type_info1")
}

// 分页查询项目一级模版信息
func SearchScoreTypeInfoIList(page, pageSize int, filters ...interface{}) ([]*ScoreTypeInfoI, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreTypeInfoI, 0)
	query := orm.NewOrm().QueryTable(TableName("score_type_info1"))
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

// 添加项目一级模版信息
func AddScoreTypeInfoI(s *ScoreTypeInfoI) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 根据ID查询项目一级模版信息
func SearchScoreTypeInfoIByID(id int) (*ScoreTypeInfoI, error) {
	s := new(ScoreTypeInfoI)
	err := orm.NewOrm().QueryTable(TableName("score_type_info1")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 根据Name查询项目一级模版信息
func SearchScoreTypeInfoIByName(name string) (*ScoreTypeInfoI, error) {
	s := new(ScoreTypeInfoI)
	err := orm.NewOrm().QueryTable(TableName("score_type_info1")).Filter("name", name).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 查询所有项目一级模版信息
func SearchAllScoreTypeInfoIList() []*ScoreTypeInfoI {
	list := make([]*ScoreTypeInfoI, 0)
	query := orm.NewOrm().QueryTable(TableName("score_type_info1"))
	query.OrderBy("id").All(&list)
	return list
}

// 更新项目一级模版信息
func (s *ScoreTypeInfoI) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}
