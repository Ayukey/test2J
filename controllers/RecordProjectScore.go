package controllers

import (
	"fmt"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 部门信息
type RecordProjectScoreController struct {
	BaseController
}

func (c *RecordProjectScoreController) Score1() {
	c.Data["pageTitle"] = "项目评分记录"
	projects := models.SearchAllProjects()
	c.Data["projects"] = projects
	c.display()
}

func (c *RecordProjectScoreController) Search1() {
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filter3 := models.DBFilter{Key: "pid", Value: pid}         // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3}

	records := models.SearchProjectReleaseRecordsByFilters(filters...)
	if len(records) == 0 {
		c.ajaxMsg(MSG_ERR, "该季度项目评分未发布")
	}

	templateRecords := logic.SearchProjectTemplate1Records(year, quarter, pid)
	list := make([]map[string]interface{}, len(templateRecords))

	for i, tr := range templateRecords {
		row := make(map[string]interface{})
		row["t1id"] = tr.Template.ID
		row["template_name"] = tr.Template.Name

		if tr.Record == nil {
			row["record_score"] = "暂无评分"
		} else {
			row["record_score"] = tr.Record.TotalScore
		}
		list[i] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *RecordProjectScoreController) Score2() {
	t1id, _ := c.GetInt("t1id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template1, _ := models.SearchProjectTemplate1ByID(t1id)
	project, _ := models.SearchProjectByID(pid)
	row := make(map[string]interface{})
	row["t1id"] = t1id
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = strconv.Itoa(year) + "年 第" + strconv.Itoa(quarter) + "季度  " + project.Name + " (" + template1.Name + ")"
	c.display()
}

func (c *RecordProjectScoreController) Search2() {
	t1id, _ := c.GetInt("t1id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, pid)
	list := make([]map[string]interface{}, len(template2Records))

	for k, tr := range template2Records {
		row := make(map[string]interface{})
		row["t2id"] = tr.Template.ID
		row["tempalte2_name"] = tr.Template.Name

		if tr.Record == nil {
			row["record2_score"] = "暂无评分"
		} else {
			row["record2_score"] = tr.Record.TotalScore
		}

		list[k] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *RecordProjectScoreController) Score3() {
	t1id, _ := c.GetInt("t1id", 0)
	t2id, _ := c.GetInt("t2id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	template1, _ := models.SearchProjectTemplate1ByID(t1id)
	template2, _ := models.SearchProjectTemplate2ByID(t2id)
	project, _ := models.SearchProjectByID(pid)
	row := make(map[string]interface{})
	row["t1id"] = t1id
	row["t2id"] = t2id
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row

	c.Data["pageTitle"] = strconv.Itoa(year) + "年 第" + strconv.Itoa(quarter) + "季度 " + project.Name + " (" + template1.Name + " -- " + template2.Name + ")"
	c.display()
}

func (c *RecordProjectScoreController) Search3() {
	t1id, _ := c.GetInt("t1id", 0)
	t2id, _ := c.GetInt("t2id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template3Records := logic.SearchProjectTemplate3Records(year, quarter, t2id, t1id, pid)
	list := make([]map[string]interface{}, len(template3Records))

	for k, tr := range template3Records {
		row := make(map[string]interface{})
		row["t3id"] = tr.Template.ID
		row["t1id"] = t1id
		row["t2id"] = t2id
		row["template3_name"] = tr.Template.Name
		row["template3_maxscore"] = tr.Template.MaxScore
		if tr.Record == nil {
			row["record3_score"] = "暂无评分"
		} else {
			row["record3_score"] = tr.Record.Score
		}
		list[k] = row
	}

	c.ajaxList(MSG_OK, "成功", list)
}

func (c *RecordProjectScoreController) Download() {
	tid, _ := c.GetInt("tid", 0)
	tname := c.GetString("tname", "")
	tscore, _ := c.GetFloat("tscore", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	templateRecords := logic.SearchProjectTemplate3Records(year, quarter, tid, ttid, pid)
	project, _ := models.SearchProjectByID(pid)

	filePath := "static/excel/ProjectScore" + "_" + strconv.Itoa(year) + "_" + strconv.Itoa(quarter) + "_" + strconv.Itoa(pid) + "_" + strconv.Itoa(tid) + ".xlsx"
	projectName := "  项目名称: " + project.Name
	yearAndQuarter := strconv.Itoa(year) + "第" + strconv.Itoa(quarter) + "季度  "

	xlsx := excelize.NewFile()

	leftStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"left"}}`)
	centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	rightStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"right"}}`)

	xlsx.SetColWidth("Sheet1", "A", "A", 12)
	xlsx.SetColWidth("Sheet1", "B", "B", 60)
	xlsx.SetColWidth("Sheet1", "C", "C", 12)
	xlsx.SetColWidth("Sheet1", "D", "D", 12)

	xlsx.MergeCell("Sheet1", "A1", "D1")
	xlsx.SetCellValue("Sheet1", "A1", tname)
	xlsx.SetCellStyle("Sheet1", "A1", "A1", centerStyle)

	xlsx.MergeCell("Sheet1", "A2", "B2")
	xlsx.SetCellValue("Sheet1", "A2", projectName)
	xlsx.SetCellStyle("Sheet1", "A2", "A2", leftStyle)

	xlsx.MergeCell("Sheet1", "C2", "D2")
	xlsx.SetCellValue("Sheet1", "C2", yearAndQuarter)
	xlsx.SetCellStyle("Sheet1", "C2", "C2", rightStyle)

	xlsx.SetCellValue("Sheet1", "A3", "序号")
	xlsx.SetCellValue("Sheet1", "B3", "检查要素")
	xlsx.SetCellValue("Sheet1", "C3", "分值")
	xlsx.SetCellValue("Sheet1", "D3", "得分")
	xlsx.SetCellStyle("Sheet1", "A3", "A3", centerStyle)
	xlsx.SetCellStyle("Sheet1", "B3", "B3", centerStyle)
	xlsx.SetCellStyle("Sheet1", "C3", "C3", centerStyle)
	xlsx.SetCellStyle("Sheet1", "D3", "D3", centerStyle)

	idex := 3

	for i, tr := range templateRecords {
		idex++
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), tr.Template.Name)
		xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), tr.Template.MaxScore)
		xlsx.SetCellStyle("Sheet1", "C"+strconv.Itoa(idex), "C"+strconv.Itoa(idex), centerStyle)
		if tr.Record != nil {
			if tr.Record.Score == -1 {
				xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), "/")
			} else {
				xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), tr.Record.Score)
			}
		} else {
			xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(idex), "")
		}
		xlsx.SetCellStyle("Sheet1", "D"+strconv.Itoa(idex), "D"+strconv.Itoa(idex), centerStyle)
	}
	idex++

	xlsx.MergeCell("Sheet1", "A"+strconv.Itoa(idex), "D"+strconv.Itoa(idex))
	xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), "总分: "+strconv.FormatFloat(tscore, 'f', 2, 64))
	xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), rightStyle)

	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect("https://scapi.sh2j.com/"+filePath, 302)
}
