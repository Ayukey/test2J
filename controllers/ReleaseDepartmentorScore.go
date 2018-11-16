package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
	"strings"
)

// 项目评分顶级模版
type ReleaseDepartmentorScoreController struct {
	BaseController
}

func (c *ReleaseDepartmentorScoreController) Release() {
	c.Data["pageTitle"] = "发布部门负责人互评"
	departmentLeaders := models.SearchAllDepartmentLeadersInDepartment()
	c.Data["departmentLeaders"] = departmentLeaders
	c.display()
}

//存储资源
func (c *ReleaseDepartmentorScoreController) AjaxSave() {
	departmentLeader := strings.Split(c.GetString("departmentLeader"), "|")
	departmentLeaderID := 0
	departmentID := 0
	if len(departmentLeader) == 2 {
		departmentLeaderID, _ = strconv.Atoi(departmentLeader[0])
		departmentID, _ = strconv.Atoi(departmentLeader[1])
	}

	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	error := logic.ReleaseDepartmentLeaderScoreRecord(year, quarter, departmentLeaderID, departmentID)

	if error == nil {
		c.ajaxMsg(MSG_OK, "")
	} else {
		c.ajaxMsg(MSG_ERR, error.Error())
	}
}
