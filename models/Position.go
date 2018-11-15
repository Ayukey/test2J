package models

import (
	"github.com/astaxie/beego/orm"
)

// 角色职位信息
type Position struct {
	ID     int    `orm:"column(id)"`
	Name   string `orm:"column(name)"`
	Status int    `orm:"column(status);default(1)"`
}

// 所有职位
func SearchAllPositions() []*Position {
	positions := make([]*Position, 0)
	query := orm.NewOrm().QueryTable(TableName("position")).Filter("status", 1)
	query.OrderBy("id").All(&positions)
	return positions
}

func AddPosition(p *Position) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(&p, "Name")
	return err
}

// 根据ID查找职位信息
func SearchPositionByID(id int) (Position, error) {
	var position Position
	err := orm.NewOrm().QueryTable(TableName("position")).Filter("id", id).One(&position)
	return position, err
}
