package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目二级模版用户权限
type Pt2Permission struct {
	ID  int `orm:"column(id)"`
	TID int `orm:"column(tid)"`
	UID int `orm:"column(uid)"`
}

// 新增项目二级模版用户权限
func AddPt2Permission(s *Pt2Permission) error {
	_, err := orm.NewOrm().Insert(s)
	return err
}

// 根据模版ID查询对应的二级项目评分模版权限
func SearchPt2PermissionsByTID(tid int) []*Pt2Permission {
	permissions := make([]*Pt2Permission, 0)
	orm.NewOrm().QueryTable(TableName("pt2_permission")).Filter("tid", tid).All(&permissions)
	return permissions
}

// 根据用户ID查询对应的二级项目评分模版权限
func SearchPt2PermissionsByUID(uid int) []*Pt2Permission {
	permissions := make([]*Pt2Permission, 0)
	orm.NewOrm().QueryTable(TableName("pt2_permission")).Filter("uid", uid).All(&permissions)
	return permissions
}

// 根据用户ID、模版ID查询对应的二级项目评分模版权限
func SearchPt2PermissionsByTUID(tid, uid int) (Pt2Permission, error) {
	var permission Pt2Permission
	err := orm.NewOrm().QueryTable(TableName("pt2_permission")).Filter("tid", tid).Filter("uid", uid).One(&permission)
	return permission, err
}

// 根据模版ID清除对应的二级项目评分模版权限
func ClearPt2PermissionsByTID(tid int) error {
	_, err := orm.NewOrm().QueryTable(TableName("pt2_permission")).Filter("tid", tid).Delete()
	return err
}
