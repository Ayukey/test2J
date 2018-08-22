package models

import (
	"github.com/astaxie/beego/orm"
)

// 角色职位信息
type PositionRoleInfo struct {
	ID   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
}

func (p *PositionRoleInfo) TableName() string {
	return TableName("position_role_info")
}

func SearchPositionRoleInfoList(page, pageSize int, filters ...interface{}) ([]*PositionRoleInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*PositionRoleInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("position_role_info"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("id").Limit(pageSize, offset).All(&list)
	return list, total
}

func SearchAllPositionRoleInfo() []*PositionRoleInfo {
	list := make([]*PositionRoleInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("position_role_info"))
	query.OrderBy("id").All(&list)
	return list
}

func AddPositionRoleInfo(p *PositionRoleInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchPositionRoleInfoByID(id int) (*PositionRoleInfo, error) {
	p := new(PositionRoleInfo)
	err := orm.NewOrm().QueryTable(TableName("position_role_info")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PositionRoleInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}
