package controllers

import (
	"jg2j_server/models"
	"strconv"
	"strings"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Data["pageTitle"] = "系统首页"
	c.TplName = "public/main.html"
}

func (c *HomeController) Start() {
	c.Data["pageTitle"] = "当前考核季度"

	quarterActive, err := models.SearchQuarterInActive()

	currentYear, currentQuarter := getCurrentYearAndQuarter()

	row := make(map[string]interface{})
	year := currentYear
	quarter := currentQuarter
	inActive := false
	if err == nil {
		year = quarterActive.Year
		quarter = quarterActive.Quarter
		inActive = true
	}
	row["year"] = year
	row["quarter"] = quarter
	row["inActive"] = inActive

	projects := models.SearchAllProjects()
	activeProjects := models.SearchAllActiveQuarterProjects(year, quarter)

	projectLeaders := models.SearchAllProjectLeadersInProject()
	activeProjectLeaders := models.SearchAllActiveQuarterProjectLeaders(year, quarter)

	departmentLeaders := models.SearchAllDepartmentLeadersInDepartment()
	activeDepartmentLeaders := models.SearchAllActiveQuarterDepartmentLeaders(year, quarter)

	c.Data["Source"] = row
	c.Data["projects"] = len(projects)
	c.Data["activeProjects"] = len(activeProjects)
	c.Data["projectLeaders"] = len(projectLeaders)
	c.Data["activeProjectLeaders"] = len(activeProjectLeaders)
	c.Data["departmentLeaders"] = len(departmentLeaders)
	c.Data["activeDepartmentLeaders"] = len(activeDepartmentLeaders)

	c.display()
}

//存储资源
func (c *HomeController) Active() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	inActive, _ := c.GetInt("inActive", 0)

	if inActive == 1 {
		err := models.ActiveQuarter(year, quarter)
		err = models.ActiveQuarterAllProjects(year, quarter)
		err = models.ActiveQuarterAllProjectLeaders(year, quarter)
		err = models.ActiveQuarterAllDepartmentLeaders(year, quarter)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		err := models.UnActiveQuarter(year, quarter)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "success")
}

func (c *HomeController) ActiveProjects() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	c.Data["pageTitle"] = strconv.Itoa(year) + "年第" + strconv.Itoa(quarter) + "季度--考核项目列表"
	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.display()
}

func (c *HomeController) ListActiveProjects() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	allActiveProjects := models.SearchAllActiveQuarterProjects(year, quarter)
	allProjects := models.SearchAllProjects()

	list := make([]map[string]interface{}, len(allProjects))
	for k, v := range allProjects {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		u, err := models.SearchUserByID(v.Leader)
		if err == nil {
			row["leaderName"] = u.Name
		}
		isActive := false

		for _, p := range allActiveProjects {
			if v.ID == p.PID {
				isActive = true
				break
			}
		}
		row["isActive"] = isActive

		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *HomeController) HandleActiveProjects() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	pid, _ := c.GetInt("pid", 0)
	inActive, _ := c.GetInt("inActive", 0)

	if inActive == 1 {
		err := models.AddQuarterActiveProject(year, quarter, pid)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		err := models.DeleteQuarterActiveProject(year, quarter, pid)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "success")
}

func (c *HomeController) ActiveProjectLeaders() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	c.Data["pageTitle"] = strconv.Itoa(year) + "年第" + strconv.Itoa(quarter) + "季度--项目负责人列表"
	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.display()
}

func (c *HomeController) ListActiveProjectLeaders() {

	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	allActiveProjectLeaders := models.SearchAllActiveQuarterProjectLeaders(year, quarter)
	allProjectLeaders := models.SearchAllProjectLeadersInProject()

	list := make([]map[string]interface{}, len(allProjectLeaders))
	for k, v := range allProjectLeaders {
		row := make(map[string]interface{})
		row["id"] = v.User.ID
		row["pid"] = v.Project.ID
		row["name"] = v.User.Name
		row["pname"] = v.Project.Name
		isActive := false

		for _, p := range allActiveProjectLeaders {
			if v.User.ID == p.UID && v.Project.ID == p.ProjectID {
				isActive = true
				break
			}
		}
		row["isActive"] = isActive

		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *HomeController) HandleActiveProjectLeaders() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	inActive, _ := c.GetInt("inActive", 0)

	projectLeader := strings.Split(c.GetString("projectLeader"), "|")
	projectLeaderID := 0
	projectID := 0
	if len(projectLeader) == 2 {
		projectLeaderID, _ = strconv.Atoi(projectLeader[0])
		projectID, _ = strconv.Atoi(projectLeader[1])
	}

	if inActive == 1 {
		err := models.AddQuarterActiveProjectLeader(year, quarter, projectLeaderID, projectID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		err := models.DeleteQuarterActiveProjectLeader(year, quarter, projectLeaderID, projectID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "success")
}

func (c *HomeController) ActiveDepartmentLeaders() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	c.Data["pageTitle"] = strconv.Itoa(year) + "年第" + strconv.Itoa(quarter) + "季度--部门负责人列表"
	row := make(map[string]interface{})
	row["year"] = year
	row["quarter"] = quarter
	c.Data["Source"] = row
	c.display()
}

func (c *HomeController) ListActiveDepartmentLeaders() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	allActiveDepartmentLeaders := models.SearchAllActiveQuarterDepartmentLeaders(year, quarter)
	allDepartmentLeaders := models.SearchAllDepartmentLeadersInDepartment()

	list := make([]map[string]interface{}, len(allDepartmentLeaders))
	for k, v := range allDepartmentLeaders {
		row := make(map[string]interface{})
		row["id"] = v.User.ID
		row["did"] = v.Department.ID
		row["name"] = v.User.Name
		row["dname"] = v.Department.Name
		isActive := false

		for _, p := range allActiveDepartmentLeaders {
			if v.User.ID == p.UID && v.Department.ID == p.DepartmentID {
				isActive = true
				break
			}
		}
		row["isActive"] = isActive

		list[k] = row
	}
	c.ajaxList(MSG_OK, "成功", list)
}

func (c *HomeController) HandleActiveDepartmentLeaders() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	inActive, _ := c.GetInt("inActive", 0)

	departmentLeader := strings.Split(c.GetString("departmentLeader"), "|") // 部门负责人（被打分）
	departmentLeaderID := 0
	departmentID := 0

	if len(departmentLeader) == 2 {
		departmentLeaderID, _ = strconv.Atoi(departmentLeader[0])
		departmentID, _ = strconv.Atoi(departmentLeader[1])
	}

	if inActive == 1 {
		err := models.AddQuarterActiveDepartmentLeader(year, quarter, departmentLeaderID, departmentID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	} else {
		err := models.DeleteQuarterActiveDepartmentLeader(year, quarter, departmentLeaderID, departmentID)
		if err != nil {
			c.ajaxMsg(MSG_ERR, err.Error())
		}
	}

	c.ajaxMsg(MSG_OK, "success")
}
