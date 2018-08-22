package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// 项目负责人评分发布记录
type ProjectorScoreSumInfo struct {
	ID          int     `orm:"column(id)"`
	ScoreUserID int     `orm:"column(scoreuser_id)"`
	UserID      int     `orm:"column(user_id)"`
	Score       float64 `orm:"column(score)"`
	Year        int     `orm:"column(year)"`
	Quarter     int     `orm:"column(quarter)"`
	ProjectID   int     `orm:"column(project_id)"`
}

func (p *ProjectorScoreSumInfo) TableName() string {
	return TableName("projector_score_sum_info")
}

func SearchProjectorScoreSumInfoList(page, pageSize int, filters ...interface{}) ([]*ProjectorScoreSumInfo, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ProjectorScoreSumInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_sum_info"))
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

func AddProjectorScoreSumInfo(p *ProjectorScoreSumInfo) (int64, error) {
	id, err := orm.NewOrm().Insert(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SearchProjectorScoreSumInfoByID(id int) (*ProjectorScoreSumInfo, error) {
	p := new(ProjectorScoreSumInfo)
	err := orm.NewOrm().QueryTable(TableName("projector_score_sum_info")).Filter("id", id).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *ProjectorScoreSumInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(p, fields...); err != nil {
		return err
	}
	return nil
}

func SearchProjectorScoreSumInfoByFilters(filters ...interface{}) []*ProjectorScoreSumInfo {
	list := make([]*ProjectorScoreSumInfo, 0)
	query := orm.NewOrm().QueryTable(TableName("projector_score_sum_info"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			fmt.Println(filters[k])
			fmt.Println(filters[k].(string))
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}
