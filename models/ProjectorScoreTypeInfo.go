package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目负责人评分模版
type ProjectorScoreTypeInfo struct {
	ID         int     `orm:"column(id)"`
	Name       string  `orm:"column(name)"`
	ScoreLimit float64 `orm:"column(score_Limit)"`
	Status     int     `orm:"column(status)"`
}

func (p *ProjectorScoreTypeInfo) TableName() string {
	return TableName("projector_score_type_info")
}

func SearchProjectorScoreTypeInfoList(page, pageSize int, filters ...interface{}) ([]*ProjectorScoreTypeInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ProjectorScoreTypeInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_type_info"))
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

func AddProjectorScoreTypeInfo(p *ProjectorScoreTypeInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchProjectorScoreTypeInfoByID(id int) (*ProjectorScoreTypeInfo, error) {
	p := new(ProjectorScoreTypeInfo)
	err := orm.NewOrm().QueryTable(TableName("projector_score_type_info")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *ProjectorScoreTypeInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func SearchAllProjectorScoreTypeInfoList() []*ProjectorScoreTypeInfo {
	list := make([]*ProjectorScoreTypeInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_type_info"))
	query.OrderBy("id").All(&list)
	return list
}
