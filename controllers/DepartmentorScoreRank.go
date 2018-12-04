package controllers

import (
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 部门负责人评分季度汇总排名
type DepartmentorScoreRankController struct {
	BaseController
}

func (c *DepartmentorScoreRankController) List() {
	c.Data["pageTitle"] = "部门负责人评分季度汇总排名"
	c.display()
}

func (c *DepartmentorScoreRankController) Search() {
	templates := models.SearchAllDepartmentLeaderTemplates()

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
	col1["width"] = "15%"

	colList = append(colList, col, col1)

	for _, t := range templates {
		col := make(map[string]interface{})
		col["field"] = "score" + strconv.Itoa(t.ID)
		col["align"] = "center"
		col["title"] = t.Name + "(" + libs.Float64ToStringWithNoZero(t.ScoreLimit) + "分)"
		col["width"] = strconv.Itoa(eachWidth) + "%"
		colList = append(colList, col)
	}

	col2 := make(map[string]interface{})
	col2["field"] = "totalscore"
	col2["align"] = "center"
	col2["title"] = "总分"
	col2["width"] = "15%"

	colList = append(colList, col2)

	out := make(map[string]interface{})
	out["status"] = MSG_OK
	out["col"] = colList
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

func (c *DepartmentorScoreRankController) Table() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}

	records := models.SearchDepartmentLeaderReleaseRecordsByOrder(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度没有任何部门负责人评分发布")
	}

	list := make([]map[string]interface{}, len(records))

	for i, record := range records {
		col := make(map[string]interface{})
		user, _ := models.SearchUserByID(record.UID)
		department, _ := models.SearchDepartmentByID(record.DepartmentID)
		templateRecords := logic.SearchDepartmentLeaderTemplateAverageRecords(year, quarter, record.UID, record.DepartmentID)

		col["id"] = i + 1
		col["name"] = user.Name + " (" + department.Name + ")"

		for _, templateRecord := range templateRecords {
			key := "score" + strconv.Itoa(templateRecord.Template.ID)
			col[key] = libs.Float64ToStringWithNoZero(templateRecord.Score)
		}

		col["totalscore"] = libs.Float64ToStringWithNoZero(record.Score)
		list[i] = col
	}

	c.ajaxList(MSG_OK, "成功", list)
}

// 下载
func (c *DepartmentorScoreRankController) Download() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	title := "部门负责人" + strconv.Itoa(year) + "第" + strconv.Itoa(quarter) + "季度互评汇总"

	filePath := "static/excel/" + title + ".xlsx"

	templates := models.SearchAllDepartmentLeaderTemplates()

	xlsx := excelize.NewFile()
	centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)

	// 标题
	colList := make([]string, 0)
	colList = append(colList, "排名", "姓名")
	for _, t := range templates {
		colList = append(colList, t.Name+"("+libs.Float64ToStringWithNoZero(t.ScoreLimit)+"分)")
	}
	colList = append(colList, "总分")

	// 内容
	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filters := []models.DBFilter{filter1, filter2}
	records := models.SearchDepartmentLeaderReleaseRecordsByOrder(filters...)

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

	for i, record := range records {
		templateRecords := logic.SearchDepartmentLeaderTemplateAverageRecords(year, quarter, record.UID, record.DepartmentID)
		user, _ := models.SearchUserByID(record.UID)
		department, _ := models.SearchDepartmentByID(record.DepartmentID)

		idex++
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), user.Name+" ("+department.Name+")")
		xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[0].Score))
		xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[1].Score))
		xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[2].Score))
		xlsx.SetCellStyle("Sheet1", "E"+strconv.Itoa(idex), "E"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[3].Score))
		xlsx.SetCellStyle("Sheet1", "F"+strconv.Itoa(idex), "F"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[4].Score))
		xlsx.SetCellStyle("Sheet1", "G"+strconv.Itoa(idex), "G"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(templateRecords[4].Score))
		xlsx.SetCellStyle("Sheet1", "H"+strconv.Itoa(idex), "H"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(idex), libs.Float64ToStringWithNoZero(record.Score))
		xlsx.SetCellStyle("Sheet1", "I"+strconv.Itoa(idex), "I"+strconv.Itoa(idex), centerStyle)
	}

	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Println(err)
	}

	// https://scapi.sh2j.com/
	c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
	// c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
}
