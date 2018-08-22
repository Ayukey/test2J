package models

import (
	"github.com/astaxie/beego/orm"
)

type STUserMappingI struct {
	ID     int `orm:"column(id)"`
	TID    int `orm:"column(tid)"`
	UserID int `orm:"column(user_id)"`
}

func (s *STUserMappingI) TableName() string {
	return TableName("st_user_mapping1")
}

func SearchSTUserMappingIList(page, pageSize int, filters ...interface{}) ([]*STUserMappingI, int64) {
	offset := (page - 1) * pageSize
	list := make([]*STUserMappingI, 0)
	query := orm.NewOrm().QueryTable(TableName("st_user_mapping1"))
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

func AddSTUserMappingI(s *STUserMappingI) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchSTUserMappingIByTUID(tid, uid int) (*STUserMappingI, error) {
	s := new(STUserMappingI)
	err := orm.NewOrm().QueryTable(TableName("st_user_mapping1")).Filter("tid", tid).Filter("user_id", uid).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func SearchSTUserMappingIByUID(uid int) []*STUserMappingI {
	list := make([]*STUserMappingI, 0)
	orm.NewOrm().QueryTable(TableName("st_user_mapping1")).Filter("user_id", uid).All(&list)
	return list
}

func (s *STUserMappingI) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

func ClearSTUserMappingIListByTID(tid int) error {
	_, err := orm.NewOrm().QueryTable(TableName("st_user_mapping1")).Filter("tid", tid).Delete()
	return err
}
