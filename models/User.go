package models

import (
	"github.com/astaxie/beego/orm"
)

// 用户
type User struct {
	ID         int    `orm:"column(id)"`
	Account    string `orm:"column(account);index;unique"`
	Password   string `orm:"column(password)"`
	Name       string `orm:"column(name)"`
	PositionID int    `orm:"column(position_id)"`
	Status     int    `orm:"column(status);default(1)"`
}

// 所有用户
func SearchAllUsers(filters ...DBFilter) []*User {
	users := make([]*User, 0)
	query := orm.NewOrm().QueryTable(TableName("user")).Filter("status", 1)

	for _, filter := range filters {
		query = query.Filter(filter.Key, filter.Value)
	}

	query.OrderBy("id").All(&users)
	return users
}

// 新增用户
func AddUser(u *User) error {
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(u, "Account")
	return err
}

// 根据ID搜索用户数据
func SearchUserByID(id int) (User, error) {
	var user User
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(&user)
	return user, err
}

// 根据账号搜索用户数据
func SearchUserByAccount(account string) (User, error) {
	var user User
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("account", account).One(&user)
	return user, err
}

// 更新用户数据
func (u *User) Update(fields ...string) error {
	_, err := orm.NewOrm().Update(u, fields...)
	return err
}

// 所有部门负责人
func SearchAllDepartmentLeaders() []*User {
	users := make([]*User, 0)
	query := orm.NewOrm().QueryTable(TableName("user")).Filter("position_id", 3).Filter("status", 1)
	query.OrderBy("id").All(&users)
	return users
}

// 所有项目负责人
func SearchAllProjectLeaders() []*User {
	users := make([]*User, 0)
	query := orm.NewOrm().QueryTable(TableName("user")).Filter("position_id", 5).Filter("status", 1)
	query.OrderBy("id").All(&users)
	return users
}
