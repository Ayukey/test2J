package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
	"strings"
)

// 发布项目负责人评分记录
type ReleaseProjectorScoreController struct {
	BaseController
}

func (c *ReleaseProjectorScoreController) Release() {
	c.Data["pageTitle"] = "发布项目负责人互评"
	projectLeaders := models.SearchAllProjectLeadersInProject()
	c.Data["projectLeaders"] = projectLeaders
	c.display()
}

//存储资源
func (c *ReleaseProjectorScoreController) AjaxSave() {
	projectLeader := strings.Split(c.GetString("projectLeader"), "|")
	projectLeaderID := 0
	projectID := 0
	if len(projectLeader) == 2 {
		projectLeaderID, _ = strconv.Atoi(projectLeader[0])
		projectID, _ = strconv.Atoi(projectLeader[1])
	}
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	error := logic.ReleaseProjectLeaderScoreRecord(year, quarter, projectLeaderID, projectID)

	if error == nil {
		c.ajaxMsg(MSG_OK, "")
	} else {
		c.ajaxMsg(MSG_ERR, error.Error())
	}
}
