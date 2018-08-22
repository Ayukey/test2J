package models

import (
	"github.com/astaxie/beego/orm"
)

type STUserMappingIII struct {
	ID     int `orm:"column(id)"`
	TID    int `orm:"column(tid)"`
	UserID int `orm:"column(user_id)"`
}

func (s *STUserMappingIII) TableName() string {
	return TableName("st_user_mapping3")
}

func SearchSTUserMappingIIIList(page, pageSize int, filters ...interface{}) ([]*STUserMappingIII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*STUserMappingIII, 0)
	query := orm.NewOrm().QueryTable(TableName("st_user_mapping3"))
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

func AddSTUserMappingIII(s *STUserMappingIII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchSTUserMappingIIIByID(id int) (*STUserMappingIII, error) {
	s := new(STUserMappingIII)
	err := orm.NewOrm().QueryTable(TableName("st_user_mapping3")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *STUserMappingIII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}
