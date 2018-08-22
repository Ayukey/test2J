package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目负责人信息
type ProjectorInfo struct {
	ID        int     `orm:"column(id)"`
	UserID    int     `orm:"column(user_id)"`
	BeginDate float64 `orm:"column(begin_date)"`
	EndDate   float64 `orm:"column(end_date)"`
}

func (p *ProjectorInfo) TableName() string {
	return TableName("projector_info")
}

func SearchProjectorInfoList(page, pageSize int, filters ...interface{}) ([]*ProjectorInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ProjectorInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_info"))
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

func SearchAllProjectorInfo() []*ProjectorInfo {
	list := make([]*ProjectorInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_info"))
	query.OrderBy("id").All(&list)
	return list
}

func AddPProjectorInfo(p *ProjectorInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchProjectorInfoByID(id int) (*ProjectorInfo, error) {
	p := new(ProjectorInfo)
	err := orm.NewOrm().QueryTable(TableName("projector_info")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func SearchProjectorInfoByUserID(id int) (*ProjectorInfo, error) {
	p := new(ProjectorInfo)
	err := orm.NewOrm().QueryTable(TableName("projector_info")).Filter("user_id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *ProjectorInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}
