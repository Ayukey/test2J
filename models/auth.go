package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Auth struct {
	ID         int `orm:"column(id)"`
	AuthName   string
	AuthURL    string `orm:"column(auth_url)"`
	UserID     int    `orm:"column(user_id)"`
	Pid        int
	Sort       int
	Icon       string
	IsShow     int
	Status     int
	CreateID   int `orm:"column(create_id)"`
	UpdateID   int `orm:"column(update_id)"`
	CreateTime int64
	UpdateTime int64
}

func (a *Auth) TableName() string {
	return TableName("sys_auth")
}

func AuthGetList(page, pageSize int, filters ...interface{}) ([]*Auth, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Auth, 0)
	query := orm.NewOrm().QueryTable(TableName("sys_auth"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("pid", "sort").Limit(pageSize, offset).All(&list)

	return list, total
}

func AuthGetListByIds(authIds string, userId int) ([]*Auth, error) {

	list1 := make([]*Auth, 0)
	var list []orm.Params
	//list:=[]orm.Params
	var err error
	if userId == 1 {
		//超级管理员
		_, err = orm.NewOrm().Raw("select id,auth_name,auth_url,pid,icon,is_show from jg2j_sys_auth where status=? order by pid asc,sort asc", 1).Values(&list)
	} else {
		_, err = orm.NewOrm().Raw("select id,auth_name,auth_url,pid,icon,is_show from jg2j_sys_auth where status=1 and id in("+authIds+") order by pid asc,sort asc", authIds).Values(&list)
	}

	for k, v := range list {
		fmt.Println(k, v)
	}

	fmt.Println(list)
	return list1, err
}

func AuthAdd(auth *Auth) (int64, error) {
	return orm.NewOrm().Insert(auth)
}

func AuthGetById(id int) (*Auth, error) {
	a := new(Auth)

	err := orm.NewOrm().QueryTable(TableName("sys_auth")).Filter("id", id).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Auth) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}
