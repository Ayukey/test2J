package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目三级模版信息
type ScoreTypeInfoIII struct {
	ID       int     `orm:"column(id)"`
	Name     string  `orm:"column(name)"`
	MaxScore float64 `orm:"column(max_score)"`
	TID      int     `orm:"column(tid)"`
	Status   int     `orm:"column(status)"`
}

func (s *ScoreTypeInfoIII) TableName() string {
	return TableName("score_type_info3")
}

// 分页查询项目三级模版信息
func SearchScoreTypeInfoIIIList(page, pageSize int, filters ...interface{}) ([]*ScoreTypeInfoIII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreTypeInfoIII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_type_info3"))
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

// 添加项目三级模版信息
func AddScoreTypeInfoIII(s *ScoreTypeInfoIII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 根据ID查询项目三级模版信息
func SearchScoreTypeInfoIIIByID(id int) (*ScoreTypeInfoIII, error) {
	s := new(ScoreTypeInfoIII)
	err := orm.NewOrm().QueryTable(TableName("score_type_info3")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 根据TID(项目二级模版信息ID)查询所有项目三级模版信息
func SearchScoreTypeInfoIIIByTID(tid int) []*ScoreTypeInfoIII {
	list := make([]*ScoreTypeInfoIII, 0)
	orm.NewOrm().QueryTable(TableName("score_type_info3")).Filter("tid", tid).All(&list)
	return list
}

// 根据Name查询项目三级模版信息
func SearchScoreTypeInfoIIIByName(name string) (*ScoreTypeInfoII, error) {
	s := new(ScoreTypeInfoII)
	err := orm.NewOrm().QueryTable(TableName("score_type_info3")).Filter("name", name).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 更新项目三级模版信息
func (s *ScoreTypeInfoIII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}
