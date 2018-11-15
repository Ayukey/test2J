package controllers

import (
	"fmt"
	"jg2j_server/logic"
	"jg2j_server/models"
	"time"
)

// 项目评分顶级模版
type ScoreRecordInfoIIIController struct {
	BaseController
}

func (c *ScoreRecordInfoIIIController) Score() {
	template1ID, _ := c.GetInt("template1_id", 0)
	template2ID, _ := c.GetInt("template2_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template2, _ := models.SearchProjectTemplate2ByID(template2ID)
	project, _ := models.SearchProjectByID(projectID)

	row := make(map[string]interface{})
	row["template1_id"] = template1ID
	row["template2_id"] = template2ID
	row["project_id"] = projectID
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = template2.Name + " (" + project.Name + ")"
	c.display()
}

func (c *ScoreRecordInfoIIIController) Search() {
	template1ID, _ := c.GetInt("template1_id", 0)
	template2ID, _ := c.GetInt("template2_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template3Records := logic.SearchProjectTemplate3Records(year, quarter, template2ID, template1ID, projectID)
	list := make([]map[string]interface{}, len(template3Records))

	for index, template3Record := range template3Records {
		row := make(map[string]interface{})
		row["template3_id"] = template3Record.Template.ID
		row["template1_id"] = template1ID
		row["template2_id"] = template2ID
		row["template3_name"] = template3Record.Template.Name
		row["template3_maxscore"] = template3Record.Template.MaxScore
		if template3Record.Record == nil {
			row["record3_score"] = "暂无评分"
		} else {
			row["record3_id"] = template3Record.Record.ID
			row["record3_score"] = template3Record.Record.Score
		}
		list[index] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ScoreRecordInfoIIIController) Edit() {
	template3ID, _ := c.GetInt("template3_id", 0)
	record3ID, _ := c.GetInt("record3_id", 0)
	template1ID, _ := c.GetInt("template1_id", 0)
	template2ID, _ := c.GetInt("template2_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template3, err := models.SearchProjectTemplate3ByID(template3ID)

	if err != nil {
		c.Ctx.WriteString("数据不存在")
		return
	}

	row := make(map[string]interface{})
	row["template3_id"] = template3ID
	row["record3_id"] = record3ID
	row["template1_id"] = template1ID
	row["template2_id"] = template2ID
	row["project_id"] = projectID
	row["year"] = year
	row["quarter"] = quarter
	row["template3_name"] = template3.Name
	row["template3_maxscore"] = template3.MaxScore

	filter1 := models.DBFilter{Key: "year", Value: year}        // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}  // 季度
	filter3 := models.DBFilter{Key: "t1id", Value: template1ID} // 对应一级模版ID
	filter4 := models.DBFilter{Key: "t2id", Value: template2ID} // 对应二级模版ID
	filter5 := models.DBFilter{Key: "t3id", Value: template3ID} // 对应二级模版ID
	filter6 := models.DBFilter{Key: "pid", Value: projectID}    // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

	result := models.SearchProjectScoreRecord3sByFilters(filters...)
	if len(result) == 0 {
		row["record3_score"] = "暂无评分"
	} else {
		record := result[0]
		row["record3_id"] = record.ID
		row["record3_score"] = record.Score
	}
	c.Data["Source"] = row
	fmt.Println(c.Data["Source"])
	c.Data["pageTitle"] = "编辑项目三级评分"
	c.display()
}

//存储资源
func (c *ScoreRecordInfoIIIController) AjaxSave() {
	template3ID, _ := c.GetInt("template3_id", 0)
	record3ID, _ := c.GetInt("record3_id", 0)
	record3Score, _ := c.GetFloat("record3_score")
	template1ID, _ := c.GetInt("template1_id", 0)
	template2ID, _ := c.GetInt("template2_id", 0)
	projectID, _ := c.GetInt("project_id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template2, _ := models.SearchProjectTemplate2ByID(template2ID)

	if template2.Percentage == 1 {
		if record3Score < 0 {
			c.ajaxMsg(MSG_ERR, "该评分项不能填写-1")
		}
	}

	if record3ID == 0 {
		record3 := new(models.ProjectScoreRecord3)
		record3.UID = 0
		record3.PID = projectID
		record3.T1ID = template1ID
		record3.T2ID = template2ID
		record3.T3ID = template3ID
		record3.Score = record3Score
		record3.Year = year
		record3.Quarter = quarter
		record3.UpdateDate = float64(time.Now().Unix())

		if err := models.AddProjectScoreRecord3(record3); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		record3, _ := models.SearchProjectScoreRecord3ByID(record3ID)
		record3.UID = 0
		record3.Score = record3Score
		record3.UpdateDate = float64(time.Now().Unix())
		// 更新三级评分数据
		if err := record3.Update(); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	// 更新二级评分数据
	template3Records := logic.SearchProjectTemplate3Records(year, quarter, template1ID, template2ID, projectID)
	realScoreII := 0.0
	maxScoreII := 0.0
	for _, tr3 := range template3Records {
		maxScoreII += tr3.Template.MaxScore
		if tr3.Record != nil {
			if tr3.Record.Score == -1 {
				maxScoreII -= tr3.Template.MaxScore
			} else {
				realScoreII += tr3.Record.Score
			}
		}
	}
	totalScoreII := realScoreII / maxScoreII * 100

	if template2.Percentage == 1 {
		totalScoreII = realScoreII
	}

	filter1 := models.DBFilter{Key: "year", Value: year}        // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter}  // 季度
	filter3 := models.DBFilter{Key: "t1id", Value: template1ID} // 一级模版ID
	filter4 := models.DBFilter{Key: "t2id", Value: template2ID} // 二级模版ID
	filter5 := models.DBFilter{Key: "pid", Value: projectID}    // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

	record2s := models.SearchProjectScoreRecord2sByFilters(filters...)

	if len(record2s) == 0 {
		record2 := new(models.ProjectScoreRecord2)
		record2.PID = projectID
		record2.T1ID = template1ID
		record2.T2ID = template2ID
		record2.TotalScore = totalScoreII
		record2.Year = year
		record2.Quarter = quarter
		record2.UpdateDate = float64(time.Now().Unix())
		models.AddProjectScoreRecord2(record2)
	} else {
		record2 := record2s[0]
		record2.TotalScore = totalScoreII
		record2.UpdateDate = float64(time.Now().Unix())
		record2.Update()
	}

	// 更新一级评分数据
	template2Records := logic.SearchProjectTemplate2Records(year, quarter, template1ID, projectID)
	totalScoreI := 0.0
	for _, tr := range template2Records {
		if tr.Record != nil {
			if tr.Record != nil {
				totalScoreI += tr.Template.Percentage * tr.Record.TotalScore
			}
		}
	}

	filter1 = models.DBFilter{Key: "year", Value: year}        // 年度
	filter2 = models.DBFilter{Key: "quarter", Value: quarter}  // 季度
	filter3 = models.DBFilter{Key: "t1id", Value: template1ID} // 一级模版ID
	filter4 = models.DBFilter{Key: "pid", Value: projectID}    // 项目ID
	filters = []models.DBFilter{filter1, filter2, filter3, filter4}

	record1s := models.SearchProjectScoreRecord1sByFilters(filters...)
	if len(record1s) == 0 {
		record1 := new(models.ProjectScoreRecord1)
		record1.PID = projectID
		record1.T1ID = template1ID
		record1.TotalScore = totalScoreI
		record1.Year = year
		record1.Quarter = quarter
		record1.UpdateDate = float64(time.Now().Unix())
		models.AddProjectScoreRecord1(record1)
	} else {
		record1 := record1s[0]
		record1.TotalScore = totalScoreI
		record1.UpdateDate = float64(time.Now().Unix())
		record1.Update()
	}

	c.ajaxMsg(MSG_OK, "success")
}
