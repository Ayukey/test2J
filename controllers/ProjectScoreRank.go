package controllers

import (
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
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

func (c *ProjectScoreRankController) Download() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	t1id, _ := c.GetInt("t1id", 0)

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

	template1, _ := models.SearchProjectTemplate1ByID(t1id)

	title := strconv.Itoa(year) + "年-第" + strconv.Itoa(quarter) + "季度评分汇总(" + template1.Name + ")"

	filePath := "static/excel/" + title + ".xlsx"

	template2s := models.SearchProjectTemplate2sByTID(t1id)

	xlsx := excelize.NewFile()
	centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)

	// 标题
	colList := make([]string, 0)
	colList = append(colList, "排名", "项目")
	for _, t := range template2s {
		colList = append(colList, t.Name+"("+libs.Float64ToStringWithNoZero(t.Percentage*100)+"%)")
	}
	colList = append(colList, "总分")

	// 内容

	if t1id == 1 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)
		xlsx.SetColWidth("Sheet1", "E", "E", 20)
		xlsx.SetColWidth("Sheet1", "F", "F", 20)
		xlsx.SetColWidth("Sheet1", "G", "G", 20)
		xlsx.SetColWidth("Sheet1", "H", "H", 20)

		xlsx.MergeCell("Sheet1", "A1", "H1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)
		xlsx.SetCellValue("Sheet1", "E2", colList[4])
		xlsx.SetCellStyle("Sheet1", "E2", "E2", centerStyle)
		xlsx.SetCellValue("Sheet1", "F2", colList[5])
		xlsx.SetCellStyle("Sheet1", "F2", "F2", centerStyle)
		xlsx.SetCellValue("Sheet1", "G2", colList[6])
		xlsx.SetCellStyle("Sheet1", "G2", "G2", centerStyle)
		xlsx.SetCellValue("Sheet1", "H2", colList[7])
		xlsx.SetCellStyle("Sheet1", "H2", "H2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["score1"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), r["score2"])
			xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), r["score3"])
			xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), r["score4"])
			xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "H"+strconv.Itoa(idex), "H"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 2 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)
		xlsx.SetColWidth("Sheet1", "E", "E", 20)
		xlsx.SetColWidth("Sheet1", "F", "F", 20)
		xlsx.SetColWidth("Sheet1", "G", "G", 20)
		xlsx.SetColWidth("Sheet1", "H", "H", 20)
		xlsx.SetColWidth("Sheet1", "I", "I", 20)

		xlsx.MergeCell("Sheet1", "A1", "I1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)
		xlsx.SetCellValue("Sheet1", "E2", colList[4])
		xlsx.SetCellStyle("Sheet1", "E2", "E2", centerStyle)
		xlsx.SetCellValue("Sheet1", "F2", colList[5])
		xlsx.SetCellStyle("Sheet1", "F2", "F2", centerStyle)
		xlsx.SetCellValue("Sheet1", "G2", colList[6])
		xlsx.SetCellStyle("Sheet1", "G2", "G2", centerStyle)
		xlsx.SetCellValue("Sheet1", "H2", colList[7])
		xlsx.SetCellStyle("Sheet1", "H2", "H2", centerStyle)
		xlsx.SetCellValue("Sheet1", "I2", colList[8])
		xlsx.SetCellStyle("Sheet1", "I2", "I2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["score1"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), r["score2"])
			xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), r["score3"])
			xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), r["score4"])
			xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(idex), r["score5"])
			xlsx.SetCellStyle("Sheet1", "H"+strconv.Itoa(idex), "H"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "I"+strconv.Itoa(idex), "I"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 3 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)
		xlsx.SetColWidth("Sheet1", "E", "E", 20)
		xlsx.SetColWidth("Sheet1", "F", "F", 20)
		xlsx.SetColWidth("Sheet1", "G", "G", 20)

		xlsx.MergeCell("Sheet1", "A1", "G1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)
		xlsx.SetCellValue("Sheet1", "E2", colList[4])
		xlsx.SetCellStyle("Sheet1", "E2", "E2", centerStyle)
		xlsx.SetCellValue("Sheet1", "F2", colList[5])
		xlsx.SetCellStyle("Sheet1", "F2", "F2", centerStyle)
		xlsx.SetCellValue("Sheet1", "G2", colList[6])
		xlsx.SetCellStyle("Sheet1", "G2", "G2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["score1"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), r["score2"])
			xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), r["score3"])
			xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 4 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)
		xlsx.SetColWidth("Sheet1", "E", "E", 20)
		xlsx.SetColWidth("Sheet1", "F", "F", 20)
		xlsx.SetColWidth("Sheet1", "G", "G", 20)

		xlsx.MergeCell("Sheet1", "A1", "G1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)
		xlsx.SetCellValue("Sheet1", "E2", colList[4])
		xlsx.SetCellStyle("Sheet1", "E2", "E2", centerStyle)
		xlsx.SetCellValue("Sheet1", "F2", colList[5])
		xlsx.SetCellStyle("Sheet1", "F2", "F2", centerStyle)
		xlsx.SetCellValue("Sheet1", "G2", colList[6])
		xlsx.SetCellStyle("Sheet1", "G2", "G2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["score1"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), r["score2"])
			xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), r["score3"])
			xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 5 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)
		xlsx.SetColWidth("Sheet1", "E", "E", 20)
		xlsx.SetColWidth("Sheet1", "F", "F", 20)
		xlsx.SetColWidth("Sheet1", "G", "G", 20)
		xlsx.SetColWidth("Sheet1", "H", "H", 20)
		xlsx.SetColWidth("Sheet1", "I", "I", 20)

		xlsx.MergeCell("Sheet1", "A1", "I1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)
		xlsx.SetCellValue("Sheet1", "E2", colList[4])
		xlsx.SetCellStyle("Sheet1", "E2", "E2", centerStyle)
		xlsx.SetCellValue("Sheet1", "F2", colList[5])
		xlsx.SetCellStyle("Sheet1", "F2", "F2", centerStyle)
		xlsx.SetCellValue("Sheet1", "G2", colList[6])
		xlsx.SetCellStyle("Sheet1", "G2", "G2", centerStyle)
		xlsx.SetCellValue("Sheet1", "H2", colList[7])
		xlsx.SetCellStyle("Sheet1", "H2", "H2", centerStyle)
		xlsx.SetCellValue("Sheet1", "I2", colList[8])
		xlsx.SetCellStyle("Sheet1", "I2", "I2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["score1"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), r["score2"])
			xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), r["score3"])
			xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), r["score4"])
			xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(idex), r["score5"])
			xlsx.SetCellStyle("Sheet1", "H"+strconv.Itoa(idex), "H"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "I"+strconv.Itoa(idex), "I"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 6 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)

		xlsx.MergeCell("Sheet1", "A1", "D1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 7 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)

		xlsx.MergeCell("Sheet1", "A1", "D1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
		}
	} else if t1id == 8 {
		xlsx.SetColWidth("Sheet1", "A", "A", 20)
		xlsx.SetColWidth("Sheet1", "B", "B", 20)
		xlsx.SetColWidth("Sheet1", "C", "C", 20)
		xlsx.SetColWidth("Sheet1", "D", "D", 20)

		xlsx.MergeCell("Sheet1", "A1", "D1")
		xlsx.SetCellValue("Sheet1", "A1", title)
		xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

		xlsx.SetCellValue("Sheet1", "A2", colList[0])
		xlsx.SetCellStyle("Sheet1", "A2", "A2", centerStyle)
		xlsx.SetCellValue("Sheet1", "B2", colList[1])
		xlsx.SetCellStyle("Sheet1", "B2", "B2", centerStyle)
		xlsx.SetCellValue("Sheet1", "C2", colList[2])
		xlsx.SetCellStyle("Sheet1", "C2", "C2", centerStyle)
		xlsx.SetCellValue("Sheet1", "D2", colList[3])
		xlsx.SetCellStyle("Sheet1", "D2", "D2", centerStyle)

		idex := 2

		for i, r := range list {
			idex++
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
			xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), r["name"])
			xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), r["score0"])
			xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), r["totalscore"])
			xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
		}
	}

	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Println(err)
	}
	// c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
	// https://scapi.sh2j.com/
	c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
}
