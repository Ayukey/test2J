package controllers

import (
	"jg2j_server/logic"
	"jg2j_server/models"
)

// 发布项目负责人评分记录
type ReleaseProjectorScoreController struct {
	BaseController
}

func (c *ReleaseProjectorScoreController) Release() {
	c.Data["pageTitle"] = "发布项目负责人互评"
	leaders := models.SearchAllProjectLeaders()
	c.Data["leaders"] = leaders
	c.display()
}

//存储资源
func (c *ReleaseProjectorScoreController) AjaxSave() {
	uid, _ := c.GetInt("user_id")
	year, _ := c.GetInt("year")
	quarter, _ := c.GetInt("quarter")

	error := logic.ReleaseProjectorScore(year, quarter, uid)

	if error == nil {
		c.ajaxMsg("", MSG_OK)
	} else {
		c.ajaxMsg(error.Error(), MSG_ERR)
	}
}
