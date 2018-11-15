package controllers

import (
	"jg2j_server/libs"
	"jg2j_server/models"
	"strings"
)

// 用户信息
type UserInfoController struct {
	BaseController
}

// 跳转用户信息模块
func (c *UserInfoController) List() {
	c.Data["pageTitle"] = "用户列表"
	positions := models.SearchAllPositions()
	c.Data["positions"] = positions
	c.display()
}

// 查询用户信息数据
func (c *UserInfoController) Table() {
	positionID, _ := c.GetInt("position_id", 0)

	positions := models.SearchAllPositions()

	filters := make([]models.DBFilter, 0)

	if positionID != 0 {
		filter1 := models.DBFilter{Key: "position_id", Value: positionID}
		filters = append(filters, filter1)
	}

	result := models.SearchAllUsers(filters...)

	users := make([]map[string]interface{}, len(result))

	for i, v := range result {
		user := make(map[string]interface{})
		user["id"] = v.ID
		user["account"] = v.Account
		user["name"] = v.Name
		user["role"] = ""
		position := getPostionRoleInfo(positions, v.PositionID)
		if position != nil {
			user["role"] = position.Name
		}
		users[i] = user
	}

	c.ajaxList(MSG_OK, "成功", users)
}

// 跳转用户信息新增模块
func (c *UserInfoController) Add() {
	c.Data["pageTitle"] = "添加用户"

	positions := models.SearchAllPositions()
	c.Data["positions"] = positions

	c.display()
}

// 跳转用户信息编辑模块
func (c *UserInfoController) Edit() {
	c.Data["pageTitle"] = "修改用户信息"

	id, _ := c.GetInt("id", 0)
	u, err := models.SearchUserByID(id)
	if err != nil {
		c.Ctx.WriteString("用户不存在")
		return
	}

	positions := models.SearchAllPositions()

	user := make(map[string]interface{})
	user["id"] = u.ID
	user["account"] = u.Account
	user["name"] = u.Name
	user["position_id"] = u.PositionID

	c.Data["user"] = user
	c.Data["positions"] = positions
	c.display()
}

// 跳转用户信息修改密码模块
func (c *UserInfoController) Change() {
	c.Data["pageTitle"] = "修改密码"

	id, _ := c.GetInt("id", 0)
	u, err := models.SearchUserByID(id)
	if err != nil {
		c.Ctx.WriteString("用户不存在")
		return
	}

	user := make(map[string]interface{})
	user["id"] = u.ID
	user["account"] = u.Account
	user["name"] = u.Name

	c.Data["user"] = user
	c.display()
}

// 保存用户信息
func (c *UserInfoController) AjaxSave() {
	id, _ := c.GetInt("id")

	if id == 0 {
		// 修改用户信息逻辑
		password := strings.TrimSpace(c.GetString("password"))
		rpassword := strings.TrimSpace(c.GetString("r_password"))

		if password != rpassword {
			c.ajaxMsg(MSG_ERR, "两遍密码不一致")
		}

		user := new(models.User)
		user.Account = strings.TrimSpace(c.GetString("account"))
		user.Name = strings.TrimSpace(c.GetString("name"))
		user.Password = libs.Md5([]byte(password))
		user.PositionID, _ = c.GetInt("position_id", 0)
		user.Status = 1

		if err := models.AddUser(user); err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}

		c.ajaxMsg(MSG_OK, "")
	}

	// 修改用户信息逻辑
	user, _ := models.SearchUserByID(id)
	user.Name = strings.TrimSpace(c.GetString("name"))
	user.PositionID, _ = c.GetInt("position_id", 0)

	if err := user.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "")
}

// 保存用户密码信息
func (c *UserInfoController) AjaxSavePassword() {
	password := strings.TrimSpace(c.GetString("password"))
	rpassword := strings.TrimSpace(c.GetString("r_password"))

	if password != rpassword {
		c.ajaxMsg(MSG_ERR, "两遍密码不一致")
	}

	id, _ := c.GetInt("id")
	user, _ := models.SearchUserByID(id)
	user.Password = libs.Md5([]byte(password))

	if err := user.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "")
}

// 删除用户
func (c *UserInfoController) DeleteUser() {
	id, err := c.GetInt("id")
	if err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}

	user, err := models.SearchUserByID(id)
	if err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	user.Status = 0

	if err = user.Update(); err != nil {
		c.ajaxMsg(MSG_ERR, err.Error())
	}
	c.ajaxMsg(MSG_OK, "")
}
