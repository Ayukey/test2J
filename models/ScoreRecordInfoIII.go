package models

import (
	"github.com/astaxie/beego/orm"
)

// 项目三级评分记录
type ScoreRecordInfoIII struct {
	ID          int     `orm:"column(id)"`
	UserID      int     `orm:"column(user_id)"`
	ProjectID   int     `orm:"column(project_id)"`
	ScoreTypeID int     `orm:"column(scoretype_id)"`
	Score       float64 `orm:"column(score)"`
	Year        int     `orm:"column(year)"`
	Quarter     int     `orm:"column(quarter)"`
	TID         int     `orm:"column(tid)"`
	UpdateDate  float64 `orm:"column(update_date)"`
	Remark      string  `orm:"column(remark)"`
}

func (s *ScoreRecordInfoIII) TableName() string {
	return TableName("score_records3")
}

// 分页查询项目三级评分记录
func SearchScoreRecordInfoIIIList(page, pageSize int, filters ...interface{}) ([]*ScoreRecordInfoIII, int64) {
	offset := (page - 1) * pageSize
	list := make([]*ScoreRecordInfoIII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records3"))
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

// 添加项目三级评分记录
func AddScoreRecordInfoIII(s *ScoreRecordInfoIII) (int64, error) {
	id, err := orm.NewOrm().Insert(s)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 更新项目三级评分记录
func (s *ScoreRecordInfoIII) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

// 查询项目三级评分记录
func SearchScoreRecordInfoIIIByFilters(filters ...interface{}) []*ScoreRecordInfoIII {
	list := make([]*ScoreRecordInfoIII, 0)
	query := orm.NewOrm().QueryTable(TableName("score_records3"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-id").All(&list)
	return list
}

// 根据ID查询项目三级评分记录
func SearchScoreRecordInfoIIIByID(id int) (*ScoreRecordInfoIII, error) {
	s := new(ScoreRecordInfoIII)
	err := orm.NewOrm().QueryTable(TableName("score_records3")).Filter("id", id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
