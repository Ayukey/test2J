package controllers

import (
	"fmt"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 部门信息
type ProjectorScoreRankController struct {
	BaseController
}

func (c *ProjectorScoreRankController) List() {
	c.Data["pageTitle"] = "项目负责人评分季度汇总排名"
	c.display()
}

func (c *ProjectorScoreRankController) Search() {
	templates := models.SearchAllProjectLeaderTemplates()

	eachWidth := 70 / len(templates)

	colList := make([]map[string]interface{}, 0)
	col := make(map[string]interface{})
	col["field"] = "id"
	col["align"] = "center"
	col["title"] = "排名"
	col["width"] = "10%"

	col1 := make(map[string]interface{})
	col1["field"] = "name"
	col1["align"] = "center"
	col1["title"] = "姓名"
	col1["width"] = "10%"

	colList = append(colList, col, col1)

	for _, t := range templates {
		col := make(map[string]interface{})
		col["field"] = "score" + strconv.Itoa(t.ID)
		col["align"] = "center"
		col["title"] = t.Name + "(" + strconv.FormatFloat(t.ScoreLimit, 'f', 0, 64) + "分)"
		col["width"] = strconv.Itoa(eachWidth) + "%"
		colList = append(colList, col)
	}

	col2 := make(map[string]interface{})
	col2["field"] = "totalscore"
	col2["align"] = "center"
	col2["title"] = "总分"
	col2["width"] = "10%"

	colList = append(colList, col2)

	fmt.Println(colList)

	out := make(map[string]interface{})
	out["status"] = MSG_OK
	out["col"] = colList
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

func (c *ProjectorScoreRankController) Table() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}

	records := models.SearchProjectLeaderReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度没有任何项目负责人评分发布")
	}

	list := make([]map[string]interface{}, len(records))

	for i, r := range records {
		col := make(map[string]interface{})
		user, _ := models.SearchUserByID(r.UID)
		project, _ := models.SearchProjectByID(r.ProjectID)
		templateRecords := logic.SearchProjectLeaderTemplateAverageRecords(year, quarter, r.UID, r.ProjectID)

		col["id"] = i + 1
		col["name"] = user.Name + " (" + project.Name + ")"

		for _, templateRecord := range templateRecords {
			key := "score" + strconv.Itoa(templateRecord.Template.ID)
			col[key] = strconv.FormatFloat(templateRecord.Score, 'f', 2, 64)
		}
		col["totalscore"] = strconv.FormatFloat(r.Score, 'f', 2, 64)
		list[i] = col
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *ProjectorScoreRankController) Download() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filePath := "static/excel/ProjectorScore" + "_" + strconv.Itoa(year) + "_" + strconv.Itoa(quarter) + ".xlsx"
	title := "项目负责人" + strconv.Itoa(year) + "第" + strconv.Itoa(quarter) + "季度互评汇总"

	templates := models.SearchAllProjectLeaderTemplates()

	xlsx := excelize.NewFile()
	centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)

	// 标题

	colList := make([]string, 0)
	colList = append(colList, "排名", "姓名")
	for _, t := range templates {
		colList = append(colList, t.Name+"("+strconv.FormatFloat(t.ScoreLimit, 'f', 0, 64)+"分)")
	}
	colList = append(colList, "总分")

	// 内容

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}
	records := models.SearchProjectLeaderReleaseRecordsByOrder(filters...)

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

	for i, r := range records {
		templateRecords := logic.SearchProjectLeaderTemplateAverageRecords(year, quarter, r.UID, r.ProjectID)
		user, _ := models.SearchUserByID(r.UID)
		project, _ := models.SearchProjectByID(r.ProjectID)

		idex++
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), user.Name+" ("+project.Name+")")
		xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), strconv.FormatFloat(templateRecords[0].Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), strconv.FormatFloat(templateRecords[1].Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), strconv.FormatFloat(templateRecords[2].Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), strconv.FormatFloat(templateRecords[3].Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), strconv.FormatFloat(templateRecords[4].Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(idex), strconv.FormatFloat(r.Score, 'f', 2, 64))
		xlsx.SetCellStyle("Sheet1", "H"+strconv.Itoa(idex), "H"+strconv.Itoa(idex), centerStyle)
	}

	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
}
