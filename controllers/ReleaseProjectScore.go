package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 项目评分顶级模版
type ReleaseProjectScoreController struct {
	BaseController
}

func (c *ReleaseProjectScoreController) Release() {
	c.Data["pageTitle"] = "发布项目评分"
	projects := models.SearchAllProjects()
	c.Data["projects"] = projects
	c.display()
}

//存储资源
func (c *ReleaseProjectScoreController) AjaxSave() {
	pid, _ := c.GetInt("project_id")
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	exist := models.ExistActiveQuarterProject(year, quarter, pid)

	if !exist {
		c.ajaxMsg(MSG_ERR, "该项目并不在该季度可发布列表中")
	}

	filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
	filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
	filter3 := models.DBFilter{Key: "pid", Value: pid}         // 项目ID
	filters := []models.DBFilter{filter1, filter2, filter3}

	records := models.SearchProjectReleaseRecordsByFilters(filters...)
	if len(records) != 0 {
		c.ajaxMsg(MSG_ERR, "该季度项目评分已发布")
	}

	template1Records := logic.SearchProjectTemplate1Records(year, quarter, pid)

	error := "以下二级评分项还未完成评分: <br>"

	list := make([]string, 0)
	for _, tr1 := range template1Records {
		template2Records := logic.SearchProjectTemplate2Records(year, quarter, tr1.Template.ID, pid)
		for _, tr2 := range template2Records {
			if tr2.Record == nil {
				list = append(list, tr2.Template.Name)
				error += tr2.Template.Name + "<br>"
			}
		}
	}

	if len(list) == 0 {
		record := new(models.ProjectReleaseRecord)
		record.Year = year
		record.Quarter = quarter
		record.PID = pid
		if err := models.AddProjectReleaseRecord(record); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
		c.ajaxMsg(MSG_OK, "")
	} else {
		c.ajaxMsg(MSG_ERR, error)
	}
}
