package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目一级模版用户权限
type Pt1Permission struct {
	ID  int `orm:"column(id)"`
	TID int `orm:"column(tid)"`
	UID int `orm:"column(uid)"`
}

// 新增项目一级模版用户权限
func AddPt1Permission(p *Pt1Permission) error {
	_, err := orm.NewOrm().Insert(p)
	return err
}

// 根据模版ID查询对应的一级项目评分模版权限
func SearchPt1PermissionsByTID(tid int) []*Pt1Permission {
	permissions := make([]*Pt1Permission, 0)
	orm.NewOrm().QueryTable(TableName("pt1_permission")).Filter("tid", tid).All(&permissions)
	return permissions
}

// 根据用户ID查询对应的一级项目评分模版权限
func SearchPt1PermissionsByUID(uid int) []*Pt1Permission {
	permissions := make([]*Pt1Permission, 0)
	orm.NewOrm().QueryTable(TableName("pt1_permission")).Filter("uid", uid).All(&permissions)
	return permissions
}

// 根据用户ID、模版ID查询对应的一级项目评分模版权限
func SearchPt1PermissionsByTUID(tid, uid int) (Pt1Permission, error) {
	var permission Pt1Permission
	err := orm.NewOrm().QueryTable(TableName("pt1_permission")).Filter("tid", tid).Filter("uid", uid).One(&permission)
	return permission, err
}

// 根据模版ID清除对应的一级项目评分模版权限
func ClearPt1PermissionsByTID(tid int) error {
	_, err := orm.NewOrm().QueryTable(TableName("pt1_permission")).Filter("tid", tid).Delete()
	return err
}
