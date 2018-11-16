package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
	"log"
	"os"
	"path/filepath"
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

func (c *RecordProjectScoreController) Download1() {
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

	project, _ := models.SearchProjectByID(pid)

	template1Records := logic.SearchProjectTemplate1Records(year, quarter, pid)

	dirName := strconv.Itoa(year) + "年-第" + strconv.Itoa(quarter) + "季度" + "_" + project.Name

	path := "static/excel/" + dirName + "/" + dirName

	dirPath := "static/excel/" + dirName

	zipPath := "static/excel/" + dirName + ".zip"

	os.Mkdir(dirPath, os.ModePerm)
	os.Mkdir(path, os.ModePerm)

	for _, template1Record := range template1Records {
		filePath := path + "/" + template1Record.Template.Name + ".xlsx"

		template2Records := logic.SearchProjectTemplate2Records(year, quarter, template1Record.Template.ID, pid)

		xlsx := excelize.NewFile()
		xlsx.DeleteSheet("Sheet1")

		leftStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"left"}}`)
		centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)
		rightStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"right"}}`)

		for _, template2Record := range template2Records {
			template2Name := template2Record.Template.Name

			xlsx.NewSheet(template2Name)
			xlsx.SetColWidth(template2Name, "A", "A", 20)
			xlsx.SetColWidth(template2Name, "B", "B", 80)
			xlsx.SetColWidth(template2Name, "C", "C", 20)
			xlsx.SetColWidth(template2Name, "D", "D", 20)
			template3Records := logic.SearchProjectTemplate3Records(year, quarter, template1Record.Template.ID, template2Record.Template.ID, pid)

			xlsx.MergeCell(template2Name, "A1", "D1")
			xlsx.SetCellStyle(template2Name, "A1", "A1", centerStyle)
			xlsx.SetCellValue(template2Name, "A1", template2Name)

			xlsx.MergeCell(template2Name, "A2", "B2")
			xlsx.SetCellStyle(template2Name, "A2", "A2", leftStyle)
			xlsx.MergeCell(template2Name, "C2", "D2")
			xlsx.SetCellStyle(template2Name, "C2", "C2", rightStyle)
			xlsx.SetCellValue(template2Name, "A2", project.Name)
			xlsx.SetCellValue(template2Name, "C2", strconv.Itoa(year)+"年-第"+strconv.Itoa(quarter)+"季度")

			xlsx.SetCellValue(template2Name, "A3", "序号")
			xlsx.SetCellValue(template2Name, "B3", "检查要素")
			xlsx.SetCellValue(template2Name, "C3", "分值")
			xlsx.SetCellValue(template2Name, "D3", "得分")
			xlsx.SetCellStyle(template2Name, "A3", "A3", centerStyle)
			xlsx.SetCellStyle(template2Name, "B3", "B3", centerStyle)
			xlsx.SetCellStyle(template2Name, "C3", "C3", centerStyle)
			xlsx.SetCellStyle(template2Name, "D3", "D3", centerStyle)

			for i, template3Record := range template3Records {
				template3 := template3Record.Template
				record3 := template3Record.Record
				index := i + 4
				A := "A" + strconv.Itoa(index)
				B := "B" + strconv.Itoa(index)
				C := "C" + strconv.Itoa(index)
				D := "D" + strconv.Itoa(index)
				xlsx.SetCellStyle(template2Name, A, A, centerStyle)
				xlsx.SetCellStyle(template2Name, B, B, leftStyle)
				xlsx.SetCellStyle(template2Name, C, C, centerStyle)
				xlsx.SetCellStyle(template2Name, D, D, centerStyle)
				xlsx.SetCellValue(template2Name, A, i+1)
				xlsx.SetCellValue(template2Name, B, template3.Name)
				xlsx.SetCellValue(template2Name, C, template3.MaxScore)
				xlsx.SetCellValue(template2Name, D, record3.Score)
			}

			totalIndex := len(template3Records) + 4
			remarkIndex := len(template3Records) + 5
			xlsx.SetCellStyle(template2Name, "A"+strconv.Itoa(totalIndex), "A"+strconv.Itoa(totalIndex), centerStyle)
			xlsx.SetCellValue(template2Name, "A"+strconv.Itoa(totalIndex), "总分")
			xlsx.MergeCell(template2Name, "B"+strconv.Itoa(totalIndex), "D"+strconv.Itoa(totalIndex))
			xlsx.SetCellStyle(template2Name, "B"+strconv.Itoa(totalIndex), "B"+strconv.Itoa(totalIndex), rightStyle)
			xlsx.SetCellValue(template2Name, "B"+strconv.Itoa(totalIndex), libs.Float64ToStringWithNoZero(template2Record.Record.TotalScore))

			xlsx.SetCellStyle(template2Name, "A"+strconv.Itoa(remarkIndex), "A"+strconv.Itoa(remarkIndex), centerStyle)
			xlsx.SetCellValue(template2Name, "A"+strconv.Itoa(remarkIndex), "总评")
			xlsx.MergeCell(template2Name, "B"+strconv.Itoa(remarkIndex), "D"+strconv.Itoa(remarkIndex))
			xlsx.SetCellStyle(template2Name, "B"+strconv.Itoa(remarkIndex), "B"+strconv.Itoa(remarkIndex), rightStyle)
			xlsx.SetCellValue(template2Name, "B"+strconv.Itoa(remarkIndex), template2Record.Record.Remark)

		}

		err := xlsx.SaveAs(filePath)
		if err != nil {
			fmt.Println(err)
			c.ajaxMsg(MSG_ERR, "该季度项目评分未发布")
		}
	}

	zipDir(dirPath, zipPath)
	// https://scapi.sh2j.com/
	c.Redirect("http://localhost:8080/"+zipPath, 302)
}

func (c *RecordProjectScoreController) Download2() {
	t1id, _ := c.GetInt("t1id", 0)
	t2id, _ := c.GetInt("t2id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, pid)

	var template2Record logic.ProjectTemplate2Record
	for _, tr := range template2Records {
		if tr.Template.ID == t2id {
			template2Record = *tr
			break
		}
	}

	template2Name := template2Record.Template.Name

	template3Records := logic.SearchProjectTemplate3Records(year, quarter, t1id, t2id, pid)
	project, _ := models.SearchProjectByID(pid)

	dirName := strconv.Itoa(year) + "年-第" + strconv.Itoa(quarter) + "季度" + "_" + project.Name

	filePath := "static/excel/" + dirName + "-" + template2Name + ".xlsx"

	xlsx := excelize.NewFile()
	xlsx.DeleteSheet("Sheet1")

	leftStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"left"}}`)
	centerStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	rightStyle, _ := xlsx.NewStyle(`{"alignment":{"horizontal":"right"}}`)

	xlsx.NewSheet(template2Name)
	xlsx.SetColWidth(template2Name, "A", "A", 20)
	xlsx.SetColWidth(template2Name, "B", "B", 80)
	xlsx.SetColWidth(template2Name, "C", "C", 20)
	xlsx.SetColWidth(template2Name, "D", "D", 20)

	xlsx.MergeCell(template2Name, "A1", "D1")
	xlsx.SetCellStyle(template2Name, "A1", "A1", centerStyle)
	xlsx.SetCellValue(template2Name, "A1", template2Name)

	xlsx.MergeCell(template2Name, "A2", "B2")
	xlsx.SetCellStyle(template2Name, "A2", "A2", leftStyle)
	xlsx.MergeCell(template2Name, "C2", "D2")
	xlsx.SetCellStyle(template2Name, "C2", "C2", rightStyle)
	xlsx.SetCellValue(template2Name, "A2", project.Name)
	xlsx.SetCellValue(template2Name, "C2", strconv.Itoa(year)+"年-第"+strconv.Itoa(quarter)+"季度")

	xlsx.SetCellValue(template2Name, "A3", "序号")
	xlsx.SetCellValue(template2Name, "B3", "检查要素")
	xlsx.SetCellValue(template2Name, "C3", "分值")
	xlsx.SetCellValue(template2Name, "D3", "得分")
	xlsx.SetCellStyle(template2Name, "A3", "A3", centerStyle)
	xlsx.SetCellStyle(template2Name, "B3", "B3", centerStyle)
	xlsx.SetCellStyle(template2Name, "C3", "C3", centerStyle)
	xlsx.SetCellStyle(template2Name, "D3", "D3", centerStyle)

	for i, template3Record := range template3Records {
		template3 := template3Record.Template
		record3 := template3Record.Record
		index := i + 4
		A := "A" + strconv.Itoa(index)
		B := "B" + strconv.Itoa(index)
		C := "C" + strconv.Itoa(index)
		D := "D" + strconv.Itoa(index)
		xlsx.SetCellStyle(template2Name, A, A, centerStyle)
		xlsx.SetCellStyle(template2Name, B, B, leftStyle)
		xlsx.SetCellStyle(template2Name, C, C, centerStyle)
		xlsx.SetCellStyle(template2Name, D, D, centerStyle)
		xlsx.SetCellValue(template2Name, A, i+1)
		xlsx.SetCellValue(template2Name, B, template3.Name)
		xlsx.SetCellValue(template2Name, C, template3.MaxScore)
		xlsx.SetCellValue(template2Name, D, record3.Score)
	}

	totalIndex := len(template3Records) + 4
	remarkIndex := len(template3Records) + 5
	xlsx.SetCellStyle(template2Name, "A"+strconv.Itoa(totalIndex), "A"+strconv.Itoa(totalIndex), centerStyle)
	xlsx.SetCellValue(template2Name, "A"+strconv.Itoa(totalIndex), "总分")
	xlsx.MergeCell(template2Name, "B"+strconv.Itoa(totalIndex), "D"+strconv.Itoa(totalIndex))
	xlsx.SetCellStyle(template2Name, "B"+strconv.Itoa(totalIndex), "B"+strconv.Itoa(totalIndex), rightStyle)
	xlsx.SetCellValue(template2Name, "B"+strconv.Itoa(totalIndex), libs.Float64ToStringWithNoZero(template2Record.Record.TotalScore))

	xlsx.SetCellStyle(template2Name, "A"+strconv.Itoa(remarkIndex), "A"+strconv.Itoa(remarkIndex), centerStyle)
	xlsx.SetCellValue(template2Name, "A"+strconv.Itoa(remarkIndex), "总评")
	xlsx.MergeCell(template2Name, "B"+strconv.Itoa(remarkIndex), "D"+strconv.Itoa(remarkIndex))
	xlsx.SetCellStyle(template2Name, "B"+strconv.Itoa(remarkIndex), "B"+strconv.Itoa(remarkIndex), rightStyle)
	xlsx.SetCellValue(template2Name, "B"+strconv.Itoa(remarkIndex), template2Record.Record.Remark)

	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Println(err)
	}
	// https://scapi.sh2j.com/
	c.Redirect("http://localhost:8080/"+filePath, 302)
}

func zipDir(dir, zipFile string) {

	fz, err := os.Create(zipFile)
	if err != nil {
		log.Fatalf("Create zip file failed: %s\n", err.Error())
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path[len(dir)+1:])
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})
}
