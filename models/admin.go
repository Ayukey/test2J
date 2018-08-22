package models

import (
	"github.com/astaxie/beego/orm"
)

type Admin struct {
	ID         int `orm:"column(id)"`
	LoginName  string
	RealName   string
	Password   string
	RoleIds    string
	Phone      string
	Email      string
	Salt       string
	LastLogin  int64
	LastIP     string `orm:"column(last_ip)"`
	Status     int
	CreateID   int `orm:"column(create_id)"`
	UpdateID   int `orm:"column(update_id)"`
	CreateTime int64
	UpdateTime int64
}

// 管理员表名
func (a *Admin) TableName() string {
	return TableName("sys_admin")
}

// 添加管理员
func AdminAdd(a *Admin) (int64, error) {
	return orm.NewOrm().Insert(a)
}

// 根据用户名查询管理员
func AdminGetByName(loginName string) (*Admin, error) {
	a := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("sys_admin")).Filter("login_name", loginName).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// 查询管理人员列表
func AdminGetList(page, pageSize int, filters ...interface{}) ([]*Admin, int64) {
	offset := (page - 1) * pageSize
	list := make([]*Admin, 0)
	query := orm.NewOrm().QueryTable(TableName("sys_admin"))
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

// 根据用户ID查询管理员
func AdminGetByID(id int) (*Admin, error) {
	a := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("sys_admin")).Filter("id", id).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// 更新管理员信息
func (a *Admin) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}
