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
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	error := logic.ReleaseDepartmentorScore(year, quarter, uid)

	if error == nil {
		c.ajaxMsg("", MSG_OK)
	} else {
		c.ajaxMsg(error.Error(), MSG_ERR)
	}
}
