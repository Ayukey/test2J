package controllers

import (
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"

	"github.com/bradfitz/slice"
)

// 部门信息
type ProjectScoreRankController struct {
	BaseController
}

type ProjectScoreRank struct {
	name       string
	totalscore float64
}

func (c *ProjectScoreRankController) List() {
	c.Data["pageTitle"] = "项目季度汇总排名"
	templates := models.SearchAllProjectTemplate1s()
	c.Data["templates"] = templates
	c.display()
}

func (c *ProjectScoreRankController) Search() {
	t1id, _ := c.GetInt("t1id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}

	records := models.SearchProjectReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度没有任何项目发布评分")
	}

	templates := models.SearchProjectTemplate2sByTID(t1id)

	eachCount := 70 / len(templates)

	colList := make([]map[string]interface{}, 0)
	col := make(map[string]interface{})
	col["field"] = "id"
	col["align"] = "center"
	col["title"] = "排名"
	col["width"] = "10%"

	col1 := make(map[string]interface{})
	col1["field"] = "name"
	col1["align"] = "center"
	col1["title"] = "项目"
	col1["width"] = "10%"

	colList = append(colList, col, col1)

	for i, t := range templates {
		col := make(map[string]interface{})
		col["field"] = "score" + strconv.Itoa(i)
		col["align"] = "center"
		col["title"] = t.Name + " (" + libs.Float64ToStringWithNoZero(t.Percentage*100) + "%)"
		col["width"] = strconv.Itoa(eachCount) + "%"
		colList = append(colList, col)
	}

	col2 := make(map[string]interface{})
	col2["field"] = "totalscore"
	col2["align"] = "center"
	col2["title"] = "总分"
	col2["width"] = "10%"

	colList = append(colList, col2)

	out := make(map[string]interface{})
	out["status"] = MSG_OK
	out["col"] = colList
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

func (c *ProjectScoreRankController) Table() {
	t1id, _ := c.GetInt("t1id", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}

	records := models.SearchProjectReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度没有任何项目发布评分")
	}

	list := make([]map[string]string, 0)

	for _, r := range records {
		col := make(map[string]string)
		project, _ := models.SearchProjectByID(r.PID)
		col["name"] = project.Name

		filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filter3 := models.DBFilter{Key: "t1id", Value: t1id}       // 模版ID
		filter4 := models.DBFilter{Key: "pid", Value: r.PID}       // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4}

		record1s := models.SearchProjectScoreRecord1sByFilters(filters...)
		if len(record1s) > 0 {
			record1 := record1s[0]
			col["totalscore"] = libs.Float64ToStringWithNoZero(record1.TotalScore)
		}

		template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, r.PID)
		for i, tr := range template2Records {
			key := "score" + strconv.Itoa(i)
			if tr.Record != nil {
				col[key] = libs.Float64ToStringWithNoZero(tr.Record.TotalScore * tr.Template.Percentage)
			}
		}

		list = append(list, col)

	}

	slice.Sort(list[:], func(i, j int) bool {
		iScore, _ := strconv.ParseFloat(list[i]["totalscore"], 64)
		jScore, _ := strconv.ParseFloat(list[j]["totalscore"], 64)
		return iScore > jScore
	})

	for i, c := range list {
		c["id"] = fmt.Sprintf("%d", i+1)
	}

	c.ajaxList(MSG_OK, "成功", list)
}
