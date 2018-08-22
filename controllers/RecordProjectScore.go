package controllers

import (
	"fmt"
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
	types := models.SearchAllScoreTypeInfoIList()
	projects := models.SearchAllProjectInfo()
	c.Data["projects"] = projects
	c.Data["types"] = types
	c.display()
}

func (c *RecordProjectScoreController) Search1() {
	tid, _ := c.GetInt("tid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "project_id", pid)
	recordList := models.SearchProjectSumPubInfoByFilters(filters...)
	if len(recordList) == 0 {
		c.ajaxMsg("该季度项目评分未发布", MSG_ERR)
	}

	typeRecordList := models.SerachScoreTypeRecordInfoIList(year, quarter, pid)
	list := make([]map[string]interface{}, 1)

	for _, v := range typeRecordList {
		row := make(map[string]interface{})
		if v.Type.ID == tid {
			row["id"] = v.Type.ID
			row["name"] = v.Type.Name

			if v.Record == nil {
				row["score"] = "暂无评分"
			} else {
				row["score"] = v.Record.TotalScore
			}
			list[0] = row
		}
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *RecordProjectScoreController) Score2() {
	tid, _ := c.GetInt("tid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	t, _ := models.SearchScoreTypeInfoIByID(tid)
	p, _ := models.SearchProjectInfoByID(pid)
	row := make(map[string]interface{})
	row["tid"] = tid
	row["pid"] = pid
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.Data["pageTitle"] = strconv.Itoa(year) + "第" + strconv.Itoa(quarter) + "季度" + " (" + p.Name + ")" + "--" + t.Name
	c.display()
}

func (c *RecordProjectScoreController) Search2() {
	tid, _ := c.GetInt("tid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordList := models.SerachScoreTypeRecordInfoIIList(year, quarter, tid, pid)
	list := make([]map[string]interface{}, len(typeRecordList))

	for k, v := range typeRecordList {
		row := make(map[string]interface{})
		row["id"] = v.Type.ID
		row["name"] = v.Type.Name

		if v.Record == nil {
			row["score"] = "暂无评分"
		} else {
			row["score"] = v.Record.TotalScore
		}

		list[k] = row
	}

	count := int64(len(list))
	c.ajaxList("成功", MSG_OK, count, list)
}

func (c *RecordProjectScoreController) Score3() {
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

	c.Data["pageTitle"] = strconv.Itoa(year) + "第" + strconv.Itoa(quarter) + "季度" + " (" + project.Name + ")" + "--" + scoreTypeII.Name
	c.display()
}

func (c *RecordProjectScoreController) Search3() {
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

func (c *RecordProjectScoreController) Download() {
	tid, _ := c.GetInt("tid", 0)
	tname := c.GetString("tname", "")
	tscore, _ := c.GetFloat("tscore", 0)
	ttid, _ := c.GetInt("ttid", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	typeRecordIIIList := models.SerachScoreTypeRecordInfoIIIList(year, quarter, tid, ttid, pid)
	project, _ := models.SearchProjectInfoByID(pid)

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

	for i, tr := range typeRecordIIIList {
		idex++
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(idex), i+1)
		xlsx.SetCellStyle("Sheet1", "A"+strconv.Itoa(idex), "A"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(idex), tr.Type.Name)
		xlsx.SetCellStyle("Sheet1", "B"+strconv.Itoa(idex), "B"+strconv.Itoa(idex), centerStyle)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(idex), tr.Type.MaxScore)
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
