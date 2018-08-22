package controllers

import (
	"jg2j_server/models"
	"time"
)

// 项目评分顶级模版
type ScoreRecordInfoIIIController struct {
	BaseController
}

func (c *ScoreRecordInfoIIIController) Score() {
	tid, _ := c.GetInt("tid", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	scoreTypeII, _ := models.SearchScoreTypeInfoIIByID(tid)
	project, _ := models.SearchProjectInfoByID(pid)
	row := make(map[string]interface{})
	row["tid"] = tid
	row["ttid"] = ttid
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = scoreTypeII.Name + " (" + project.Name + ")"
	c.display()
}

func (c *ScoreRecordInfoIIIController) Search() {
	tid, _ := c.GetInt("tid", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordList := models.SerachScoreTypeRecordInfoIIIList(year, quarter, tid, ttid, pid)
	list := make([]map[string]interface{}, len(typeRecordList))

	for k, v := range typeRecordList {
		row := make(map[string]interface{})
		row["id"] = v.Type.ID
		row["tid"] = tid
		row["ttid"] = ttid
		row["name"] = v.Type.Name
		row["max_score"] = v.Type.MaxScore
		if v.Record == nil {
			row["score"] = "暂无评分"
		} else {
			row["rid"] = v.Record.ID
			row["score"] = v.Record.Score
		}
		list[k] = row
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *ScoreRecordInfoIIIController) Edit() {
	id, _ := c.GetInt("id", 0)
	rid, _ := c.GetInt("rid", 0)
	tid, _ := c.GetInt("tid", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	t, err := models.SearchScoreTypeInfoIIIByID(id)

	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["id"] = id
	row["rid"] = rid
	row["tid"] = tid
	row["ttid"] = ttid
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	row["name"] = t.Name
	row["maxscore"] = t.MaxScore

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "scoretype_id", id)
	filters = append(filters, "project_id", pid)
	filters = append(filters, "tid", tid)
	result := models.SearchScoreRecordInfoIIIByFilters(filters...)
	if len(result) == 0 {
		row["score"] = "暂无评分"
	} else {
		record := result[0]
		row["rid"] = record.ID
		row["score"] = record.Score
	}
	c.Data["Source"] = row
	c.Data["pageTitle"] = "编辑项目三级评分"
	c.display()
}

//存储资源
func (c *ScoreRecordInfoIIIController) AjaxSave() {
	id, _ := c.GetInt("id", 0)
	rid, _ := c.GetInt("rid", 0)
	tid, _ := c.GetInt("tid", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	score, _ := c.GetFloat("score")

	typeII, _ := models.SearchScoreTypeInfoIIByID(tid)

	if typeII.Percentage == 1 {
		if score < 0 {
			c.ajaxMsg("该评分项不能填写-1", MSG_ERR)
		}
	}

	if rid == 0 {
		recordIII := new(models.ScoreRecordInfoIII)
		recordIII.UserID = 0
		recordIII.ProjectID = pid
		recordIII.ScoreTypeID = id
		recordIII.Score = score
		recordIII.Year = year
		recordIII.Quarter = quarter
		recordIII.TID = tid
		recordIII.UpdateDate = float64(time.Now().Unix())

		if _, err := models.AddScoreRecordInfoIII(recordIII); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		recordIII, _ := models.SearchScoreRecordInfoIIIByID(rid)
		recordIII.Score = score
		recordIII.UpdateDate = float64(time.Now().Unix())
		// 更新三级评分数据
		if err := recordIII.Update(); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
	}

	// 更新二级评分数据
	typeRecordIIIList := models.SerachScoreTypeRecordInfoIIIList(year, quarter, tid, ttid, pid)
	realScoreII := 0.0
	maxScoreII := 0.0
	for _, v := range typeRecordIIIList {
		maxScoreII += v.Type.MaxScore
		if v.Record != nil {
			if v.Record.Score == -1 {
				maxScoreII -= v.Type.MaxScore
			} else {
				realScoreII += v.Record.Score
			}
		}
	}
	totalScoreII := realScoreII / maxScoreII * 100

	if typeII.Percentage == 1 {
		totalScoreII = realScoreII
	}

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "tid", ttid)
	filters = append(filters, "project_id", pid)
	filters = append(filters, "scoretype_id", tid)
	recordIIList := models.SearchScoreRecordInfoIIByFilters(filters...)

	if len(recordIIList) == 0 {
		recordII := new(models.ScoreRecordInfoII)
		recordII.ProjectID = pid
		recordII.ScoreTypeID = tid

		recordII.TotalScore = totalScoreII
		recordII.Year = year
		recordII.Quarter = quarter
		recordII.TID = ttid
		recordII.UpdateDate = float64(time.Now().Unix())
		models.AddScoreRecordInfoII(recordII)
	} else {
		recordII := recordIIList[0]
		recordII.TotalScore = totalScoreII
		recordII.UpdateDate = float64(time.Now().Unix())
		recordII.Update()
	}

	// 更新一级评分数据
	typeRecordIIList := models.SerachScoreTypeRecordInfoIIList(year, quarter, ttid, pid)
	totalScoreI := 0.0
	for _, v := range typeRecordIIList {
		if v.Record != nil {
			if v.Record != nil {
				totalScoreI += v.Type.Percentage * v.Record.TotalScore
			}
		}
	}

	filters = make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "project_id", pid)
	filters = append(filters, "scoretype_id", ttid)
	recordIList := models.SearchScoreRecordInfoIByFilters(filters...)
	if len(recordIList) == 0 {
		recordI := new(models.ScoreRecordInfoI)
		recordI.ProjectID = pid
		recordI.ScoreTypeID = ttid
		recordI.TotalScore = totalScoreI
		recordI.Year = year
		recordI.Quarter = quarter
		recordI.UpdateDate = float64(time.Now().Unix())
		models.AddScoreRecordInfoI(recordI)
	} else {
		recordI := recordIList[0]
		recordI.TotalScore = totalScoreI
		recordI.UpdateDate = float64(time.Now().Unix())
		recordI.Update()
	}

	c.ajaxMsg("", MSG_OK)
}
