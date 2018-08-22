package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目二级模版信息
type ScoreTypeInfoII struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	Percentage float64 `orm:"column(percentage)"`
	TID        int     `orm:"column(tid)"`
	Status     int     `orm:"column(status)"`
}

func (s *ScoreTypeInfoII) TableName() string {
	return TableName("score_type_info2")
}

// 分页查询项目二级模版信息
func SearchScoreTypeInfoIIList(page, pageSize int, filters ...interface{}) ([]*ScoreTypeInfoII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreTypeInfoII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_type_info2"))
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

// 添加项目二级模版信息
func AddScoreTypeInfoII(s *ScoreTypeInfoII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 根据ID查询项目二级模版信息
func SearchScoreTypeInfoIIByID(id int) (*ScoreTypeInfoII, error) {
	s := new(ScoreTypeInfoII)
	err := orm.NewOrm().QueryTable(TableName("score_type_info2")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 根据TID(项目一级模版信息ID)查询所有项目二级模版信息
func SearchScoreTypeInfoIIByTID(tid int) []*ScoreTypeInfoII {
	list := make([]*ScoreTypeInfoII, 0)
	orm.NewOrm().QueryTable(TableName("score_type_info2")).Filter("tid", tid).All(&list)
	return list
}

// 根据Name查询项目二级模版信息
func SearchScoreTypeInfoIIByName(name string) (*ScoreTypeInfoII, error) {
	s := new(ScoreTypeInfoII)
	err := orm.NewOrm().QueryTable(TableName("score_type_info2")).Filter("name", name).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 更新项目二级模版信息
func (s *ScoreTypeInfoII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}
