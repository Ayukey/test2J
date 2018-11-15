package controllers

import (
	"fmt"
	"jg2j_server/models"
	"strconv"
)

// 用户信息
type STUserMappingIIController struct {
	BaseController
}

func (c *STUserMappingIIController) List() {
	c.Data["pageTitle"] = "二级评分模版用户权限列表"
	tid, _ := c.GetInt("id", 0)

	result := models.SearchAllUsers()
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["account"] = v.Account
		row["name"] = v.Name
		_, err := models.SearchPt2PermissionsByTUID(tid, v.ID)
		if err != nil {
			row["on"] = 0
		} else {
			row["on"] = 1
		}
		list[k] = row
	}
	c.Data["tid"] = tid
	c.Data["users"] = list
	c.display()
}

//存储资源
func (c *STUserMappingIIController) AjaxSave() {
	println("===========D========")
	tid, _ := c.GetInt("tid")

	users := make([]string, 0)
	for key, _ := range c.Ctx.Input.Context.Request.Form {
		if key != "tid" {
			users = append(users, key)
		}
	}

	fmt.Println(users)

	models.ClearPt2PermissionsByTID(tid)

	for _, u := range users {
		permission := new(models.Pt2Permission)
		permission.TID = tid
		permission.UID, _ = strconv.Atoi(u)
		models.AddPt2Permission(permission)
	}
	c.ajaxMsg(MSG_OK, "")
}
