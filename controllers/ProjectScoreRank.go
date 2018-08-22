package controllers

import (
	"fmt"
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
	types := models.SearchAllScoreTypeInfoIList()
	c.Data["types"] = types
	c.display()
}

func (c *ProjectScoreRankController) Search() {
	tid, _ := c.GetInt("tid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	projectList := models.SearchProjectSumPubInfoByFilters(filters...)
	if len(projectList) == 0 {
		c.ajaxMsg("该季度没有任何项目发布评分", MSG_ERR)
	}

	typeList := models.SearchScoreTypeInfoIIByTID(tid)

	eachCount := 70 / len(typeList)

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

	for i, t := range typeList {
		col := make(map[string]interface{})
		col["field"] = "score" + strconv.Itoa(i)
		col["align"] = "center"
		col["title"] = t.Name + " (" + strconv.FormatFloat(t.Percentage*100, 'f', 1, 64) + "%)"
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
	tid, _ := c.GetInt("tid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	projectList := models.SearchProjectSumPubInfoByFilters(filters...)
	if len(projectList) == 0 {
		c.ajaxMsg("该季度没有任何项目发布评分", MSG_ERR)
	}

	list := make([]map[string]string, 0)

	for _, p := range projectList {
		col := make(map[string]string)
		project, _ := models.SearchProjectInfoByID(p.ProjectID)
		col["name"] = project.Name

		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "scoretype_id", tid)
		filters = append(filters, "project_id", p.ProjectID)

		recordIList := models.SearchScoreRecordInfoIByFilters(filters...)
		if len(recordIList) > 0 {
			recordI := recordIList[0]
			col["totalscore"] = strconv.FormatFloat(recordI.TotalScore, 'f', 1, 64)
		}

		recordIIList := models.SerachScoreTypeRecordInfoIIList(year, quarter, tid, p.ProjectID)
		for i, r := range recordIIList {
			key := "score" + strconv.Itoa(i)
			if r.Record != nil {
				col[key] = strconv.FormatFloat(r.Record.TotalScore*r.Type.Percentage, 'f', 1, 64)
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

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}
