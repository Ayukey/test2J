package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 项目评分顶级模版
type ReleaseDepartmentorScoreController struct {
	BaseController
}

func (c *ReleaseDepartmentorScoreController) Release() {
	c.Data["pageTitle"] = "发布部门负责人互评"
	leaders := models.SearchAllDepartmentLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *ReleaseDepartmentorScoreController) AjaxSave() {
	uid, _ := c.GetInt("user_id")
	departmentID, _ := c.GetInt("departmentId", 0)
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	error := logic.ReleaseDepartmentLeaderScoreRecord(year, quarter, uid, departmentID)

	if error == nil {
		c.ajaxMsg(MSG_OK, "")
	} else {
		c.ajaxMsg(MSG_ERR, error.Error())
	}
}
