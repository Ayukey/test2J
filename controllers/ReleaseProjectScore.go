package controllers

import (
	"jg2j_server/models"
)

// 项目评分顶级模版
type ReleaseProjectScoreController struct {
	BaseController
}

func (c *ReleaseProjectScoreController) Release() {
	c.Data["pageTitle"] = "发布项目评分"
	projectInfos := models.SearchAllProjectInfo()
	c.Data["projectList"] = projectInfos
	c.display()
}

//存储资源
func (c *ReleaseProjectScoreController) AjaxSave() {
	pid, _ := c.GetInt("project_id")
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	filters := make([]interface{}, 0)
	filters = append(filters, "year", year)
	filters = append(filters, "quarter", quarter)
	filters = append(filters, "project_id", pid)
	recordList := models.SearchProjectSumPubInfoByFilters(filters...)
	if len(recordList) != 0 {
		c.ajaxMsg("该季度项目评分已发布", MSG_ERR)
	}

	typeReocrdIList := models.SerachScoreTypeRecordInfoIList(year, quarter, pid)

	error := "以下二级评分项还未完成评分: <br>"

	list := make([]string, 0)
	for _, t := range typeReocrdIList {
		typeReocrdIIList := models.SerachScoreTypeRecordInfoIIList(year, quarter, t.Type.ID, pid)
		for _, r := range typeReocrdIIList {
			if r.Record == nil {
				list = append(list, r.Type.Name)
				error += r.Type.Name + "<br>"
			}
		}
	}

	if len(list) == 0 {
		info := new(models.ProjectSumPubInfo)
		info.Year = year
		info.Quarter = quarter
		info.ProjectID = pid
		if _, err := models.AddProjectSumPubInfo(info); err != nil {
			c.ajaxMsg(err.Error(), MSG_ERR)
		}
		c.ajaxMsg("", MSG_OK)
	} else {
		c.ajaxMsg(error, MSG_ERR)
	}
}
