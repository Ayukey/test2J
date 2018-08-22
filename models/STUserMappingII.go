package models

import (
	"github.com/astaxie/beego/orm"
)

type STUserMappingII struct {
	ID     int `orm:"column(id)"`
	TID    int `orm:"column(tid)"`
	UserID int `orm:"column(user_id)"`
}

func (s *STUserMappingII) TableName() string {
	return TableName("st_user_mapping2")
}

func SearchSTUserMappingIIList(page, pageSize int, filters ...interface{}) ([]*STUserMappingII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*STUserMappingII, 0)
	query := orm.NewOrm().QueryTable(TableName("st_user_mapping2"))
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

func AddSTUserMappingII(s *STUserMappingII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchSTUserMappingIIByID(id int) (*STUserMappingII, error) {
	s := new(STUserMappingII)
	err := orm.NewOrm().QueryTable(TableName("st_user_mapping2")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *STUserMappingII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

func SearchSTUserMappingIIByTUID(tid, uid int) (*STUserMappingII, error) {
	s := new(STUserMappingII)
	err := orm.NewOrm().QueryTable(TableName("st_user_mapping2")).Filter("tid", tid).Filter("user_id", uid).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func SearchSTUserMappingIIByUID(uid int) []*STUserMappingII {
	list := make([]*STUserMappingII, 0)
	orm.NewOrm().QueryTable(TableName("st_user_mapping2")).Filter("user_id", uid).All(&list)
	return list
}

func ClearSTUserMappingIIListByTID(tid int) error {
	_, err := orm.NewOrm().QueryTable(TableName("st_user_mapping2")).Filter("tid", tid).Delete()
	return err
}
