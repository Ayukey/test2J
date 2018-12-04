package controllers

import (
	"encoding/json"
	"fmt"
	"jg2j_server/libs"
	"jg2j_server/logic"
	"jg2j_server/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/bradfitz/slice"
)

type ApiInfterfaceController struct {
	beego.Controller
}

// 用户登录请求
func (c *ApiInfterfaceController) GetCurrentQuarter() {
	active, err := models.SearchQuarterInActive()
	year, quarter := getCurrentYearAndQuarter()
	inActive := false
	if err == nil {
		inActive = true
		year = active.Year
		quarter = active.Quarter
	}

	info := make(map[string]interface{})
	info["quarter"] = strconv.Itoa(year) + "-" + strconv.Itoa(quarter)
	info["inActive"] = inActive

	c.apiMsg("操作成功", true, 200, info)
}

// 用户登录请求
type UserLoginRequest struct {
	UserName string `json:"UserName"`
	UserPwd  string `json:"UserPwd"`
}

// 用户登录
func (c *ApiInfterfaceController) UserLogin() {
	var request UserLoginRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	user, err := models.SearchUserByAccount(request.UserName)

	if err != nil {
		c.apiMsg("账号不存在", false, 500, nil)
	}

	if user.Password != request.UserPwd {
		c.apiMsg("密码错误", false, 500, nil)
	}

	projects := models.SearchProjectsByLeader(user.ID)
	projectIDs := make([]int, len(projects))
	for i, p := range projects {
		projectIDs[i] = p.ID
	}

	departments := models.SearchDepartmentsByLeader(user.ID)
	departmentIDs := make([]int, len(departments))
	for i, d := range departments {
		departmentIDs[i] = d.ID
	}

	info := make(map[string]interface{})
	info["uid"] = user.ID
	info["departmentIds"] = departmentIDs
	info["projectIds"] = projectIDs
	info["name"] = user.Name
	info["roleId"] = user.PositionID
	info["account"] = user.Account
	info["password"] = user.Password
	c.apiMsg("操作成功", true, 200, info)
}

// 修改密码请求
type ModifyPwdRequest struct {
	Account     string `json:"account"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// 用户修改密码(over)
func (c *ApiInfterfaceController) ModifyPwd() {
	var request ModifyPwdRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	user, err := models.SearchUserByAccount(request.Account)
	if err != nil {
		c.apiMsg("账号不存在", false, 500, nil)
	}

	if user.Password != request.OldPassword {
		c.apiMsg("原始密码错误", false, 500, nil)
	}

	user.Password = request.NewPassword

	err = user.Update()
	if err != nil {
		c.apiMsg("操作失败", false, 500, nil)
	}

	c.apiMsg("操作成功", true, 200, nil)
}

// 拉取当前季度可评分的项目列表
func (c *ApiInfterfaceController) GetCanScoreProjects() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 {
		projects := models.SearchAllActiveQuarterProjects(year, quarter)
		list := make([]map[string]interface{}, 0)
		for _, p := range projects {
			info := make(map[string]interface{})
			info["ID"] = p.PID
			project, err := models.SearchProjectByID(p.PID)
			if err == nil {
				info["Name"] = project.Name
				list = append(list, info)
			}

		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取用户项目一级评分项权限数据
func (c *ApiInfterfaceController) GetSTMapping1() {
	UserID, _ := c.GetInt("UserID")
	permissions := models.SearchPt1PermissionsByUID(UserID)
	list := make([]map[string]interface{}, len(permissions))
	for i, permission := range permissions {
		info := make(map[string]interface{})
		info["TID"] = permission.TID
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 拉取用户项目二级评分项权限数据
func (c *ApiInfterfaceController) GetSTMapping2() {
	UserID, _ := c.GetInt("UserID")
	permissions := models.SearchPt2PermissionsByUID(UserID)
	list := make([]map[string]interface{}, len(permissions))
	for i, permission := range permissions {
		info := make(map[string]interface{})
		info["TID"] = permission.TID
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 拉取项目一级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType1AndScore() {
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 && pid > 0 {
		template1Records := logic.SearchProjectTemplate1Records(year, quarter, pid)
		list := make([]map[string]interface{}, len(template1Records))
		for i, tr := range template1Records {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			if tr.Record == nil {
				info["TotalScore"] = 0
			} else {
				info["TotalScore"] = tr.Record.TotalScore //libs.Float64ToStringWithNoZero(f.Record.TotalScore)
			}
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取项目二级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType2() {
	t1id, _ := c.GetInt("t1id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 && t1id > 0 && pid > 0 {
		template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, pid)
		list := make([]map[string]interface{}, len(template2Records))
		for i, tr := range template2Records {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			if tr.Record == nil {
				info["TotalScore"] = 0
				info["Remark"] = ""
			} else {
				info["TotalScore"] = tr.Record.TotalScore //libs.Float64ToStringWithNoZero(f.Record.TotalScore)
				info["Remark"] = tr.Record.Remark
			}
			list[i] = info
		}
		da := make(map[string]interface{})
		da["Scores"] = list
		c.apiMsg("操作成功", true, 200, da)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取项目三级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType3() {
	t1id, _ := c.GetInt("t1id", 0)
	t2id, _ := c.GetInt("t2id", 0)
	pid, _ := c.GetInt("pid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 && t1id > 0 && t2id > 0 && pid > 0 {
		template3Records := logic.SearchProjectTemplate3Records(year, quarter, t1id, t2id, pid)
		list := make([]map[string]interface{}, len(template3Records))
		for i, tr := range template3Records {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			info["MaxScore"] = tr.Template.MaxScore
			if tr.Record == nil {
				info["TotalScore"] = 0
			} else {
				info["TotalScore"] = tr.Record.Score //libs.Float64ToStringWithNoZero(f.Record.Score)
			}
			list[i] = info
		}
		da := make(map[string]interface{})
		da["Scores"] = list
		c.apiMsg("操作成功", true, 200, da)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 项目三级评分项打分 同时更新相关的一、二级评分项分数
type ProjectScoreRequest struct {
	Pid     int    `json:"pid"`
	Year    int    `json:"year"`
	Quarter int    `json:"quarter"`
	Scstr   string `json:"scstr"`
	T1id    int    `json:"t1id"`
	T2id    int    `json:"t2id"`
	Userid  int    `json:"userid"`
	Remark  string `json:"remark"`
}

func (c *ApiInfterfaceController) ProjectScore() {
	var request ProjectScoreRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.Userid
	t1id := request.T1id
	t2id := request.T2id
	pid := request.Pid
	year := request.Year
	quarter := request.Quarter
	remark := request.Remark // 总结

	si := strings.Split(request.Scstr, ",")
	scoreInfo := make([]map[string]interface{}, 0)
	for _, s := range si {
		sii := strings.Split(s, "|")
		if len(sii) == 2 {
			info := make(map[string]interface{})
			info["ID"] = sii[0]
			info["Score"] = sii[1]
			scoreInfo = append(scoreInfo, info)
		}
	}

	if year > 0 && quarter > 0 && pid > 0 && t1id > 0 && t2id > 0 {
		for _, info := range scoreInfo {
			t3id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)

			filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
			filter3 := models.DBFilter{Key: "t1id", Value: t1id}       // 对应一级模版ID
			filter4 := models.DBFilter{Key: "t2id", Value: t2id}       // 对应二级模版ID
			filter5 := models.DBFilter{Key: "t3id", Value: t3id}       // 对应二级模版ID
			filter6 := models.DBFilter{Key: "pid", Value: pid}         // 项目ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

			record3List := models.SearchProjectScoreRecord3sByFilters(filters...)

			if len(record3List) == 0 {
				record3 := new(models.ProjectScoreRecord3)
				record3.UID = uid
				record3.PID = pid
				record3.Score = score
				record3.Year = year
				record3.Quarter = quarter
				record3.T1ID = t1id
				record3.T2ID = t2id
				record3.T3ID = t3id
				record3.UpdateDate = float64(time.Now().Unix())
				models.AddProjectScoreRecord3(record3)
			} else {
				record3 := record3List[0]
				record3.UID = uid
				record3.Score = score
				record3.UpdateDate = float64(time.Now().Unix())
				record3.Update()
			}
		}

		// 更新二级评分数据
		template2, _ := models.SearchProjectTemplate2ByID(t2id)

		template3Records := logic.SearchProjectTemplate3Records(year, quarter, t1id, t2id, pid)
		realScoreII := 0.0
		maxScoreII := 0.0
		for _, v := range template3Records {
			maxScoreII += v.Template.MaxScore
			if v.Record != nil {
				if v.Record.Score == -1 {
					maxScoreII -= v.Template.MaxScore
				} else {
					realScoreII += v.Record.Score
				}
			}
		}
		totalScoreII := realScoreII / maxScoreII * 100

		if template2.Percentage == 1 {
			totalScoreII = realScoreII
		}

		filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filter3 := models.DBFilter{Key: "t1id", Value: t1id}       // 一级模版ID
		filter4 := models.DBFilter{Key: "t2id", Value: t2id}       // 二级模版ID
		filter5 := models.DBFilter{Key: "pid", Value: pid}         // 项目ID
		filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}
		record2s := models.SearchProjectScoreRecord2sByFilters(filters...)

		if len(record2s) == 0 {
			record2 := new(models.ProjectScoreRecord2)
			record2.PID = pid
			record2.T1ID = t1id
			record2.T2ID = t2id
			record2.TotalScore = totalScoreII
			record2.Year = year
			record2.Remark = remark
			record2.Quarter = quarter
			record2.UpdateDate = float64(time.Now().Unix())
			models.AddProjectScoreRecord2(record2)
		} else {
			record2 := record2s[0]
			record2.Remark = remark
			record2.TotalScore = totalScoreII
			record2.UpdateDate = float64(time.Now().Unix())
			record2.Update()
		}

		// 更新一级评分数据
		template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, pid)
		totalScoreI := 0.0
		for _, v := range template2Records {
			if v.Record != nil {
				if v.Record != nil {
					totalScoreI += v.Template.Percentage * v.Record.TotalScore
				}
			}
		}

		filter1 = models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 = models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filter3 = models.DBFilter{Key: "t1id", Value: t1id}       // 一级模版ID
		filter4 = models.DBFilter{Key: "pid", Value: pid}         // 项目ID
		filters = []models.DBFilter{filter1, filter2, filter3, filter4}

		record1s := models.SearchProjectScoreRecord1sByFilters(filters...)
		if len(record1s) == 0 {
			record1 := new(models.ProjectScoreRecord1)
			record1.PID = pid
			record1.T1ID = t1id
			record1.TotalScore = totalScoreI
			record1.Year = year
			record1.Quarter = quarter
			record1.UpdateDate = float64(time.Now().Unix())
			models.AddProjectScoreRecord1(record1)
		} else {
			record1 := record1s[0]
			record1.TotalScore = totalScoreI
			record1.UpdateDate = float64(time.Now().Unix())
			record1.Update()
		}
		c.apiMsg("操作成功", true, 200, nil)
	}
	c.apiMsg("数据异常", false, 500, nil)
}

// 更新或修改二级评分的总结
type ProjectRemark2Request struct {
	Pid     int    `json:"pid"`
	Year    int    `json:"year"`
	Quarter int    `json:"quarter"`
	Remark  string `json:"remark"`
	T2id    int    `json:"t2id"`
	T1id    int    `json:"t1id"`
}

func (c *ApiInfterfaceController) ProjectRemark2() {
	var request ProjectRemark2Request
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	t1id := request.T1id     // 二级评分模版ID
	t2id := request.T2id     // 二级评分模版ID
	pid := request.Pid       // 项目ID
	remark := request.Remark // 总结
	year := request.Year
	quarter := request.Quarter

	template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, pid)
	for _, tr := range template2Records {
		if tr.Template.ID == t2id {
			if tr.Record != nil {
				tr.Record.Remark = remark
				tr.Record.Update()
			} else {
				fmt.Println("木有找到记录")
				record2 := new(models.ProjectScoreRecord2)
				record2.PID = pid
				record2.T2ID = t2id
				record2.Year = year
				record2.Quarter = quarter
				record2.T1ID = t1id
				record2.Remark = remark
				record2.UpdateDate = float64(time.Now().Unix())
				models.AddProjectScoreRecord2(record2)
			}
		}
	}
	c.apiMsg("操作成功", true, 200, nil)
}

// 拉取当前季度可评分的项目负责人列表(并赋上当前用户所评的分数)
func (c *ApiInfterfaceController) GetProjectors() {
	uid, _ := c.GetInt("uid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 {
		projectLeaders := models.SearchAllActiveQuarterProjectLeaders(year, quarter)
		list := make([]map[string]interface{}, len(projectLeaders))
		for i, p := range projectLeaders {
			project, _ := models.SearchProjectByID(p.ProjectID)
			leader, _ := models.SearchUserByID(p.UID)
			info := make(map[string]interface{})
			info["ProjectID"] = project.ID
			info["ProjectName"] = project.Name
			info["UserName"] = leader.Name
			info["UserID"] = leader.ID

			filter1 := models.DBFilter{Key: "year", Value: year}             // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter}       // 季度
			filter3 := models.DBFilter{Key: "uid", Value: leader.ID}         // 被打分人ID
			filter4 := models.DBFilter{Key: "suid", Value: uid}              // 打分人ID
			filter5 := models.DBFilter{Key: "project_id", Value: project.ID} // 项目ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

			records := models.SearchProjectLeaderSumScoreRecordsByFilters(filters...)

			if len(records) == 1 {
				info["TScore"] = libs.Float64ToStringWithNoZero(records[0].Score)
			} else {
				info["TScore"] = 0
			}
			info["year"] = year
			info["year"] = quarter
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人列表发布后数据
func (c *ApiInfterfaceController) GetProjectorsBySumData() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 {
		filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filters := []models.DBFilter{filter1, filter2}

		records := models.SearchProjectLeaderReleaseRecordsByOrder(filters...)

		list := make([]map[string]interface{}, len(records))

		for i, r := range records {
			info := make(map[string]interface{})
			user, _ := models.SearchUserByID(r.UID)
			project, _ := models.SearchProjectByID(r.ProjectID)
			info["ID"] = i + 1
			info["ProjectID"] = project.ID
			info["ProjectName"] = project.Name
			info["UserName"] = user.Name
			info["UserID"] = user.ID
			info["TScore"] = libs.Float64ToStringWithNoZero(r.Score)
			info["Year"] = year
			info["Quarter"] = quarter
			list[i] = info
		}

		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人评分记录详情列表
func (c *ApiInfterfaceController) GetProjectorScoreTypeAndScorePersonal() {
	uid, _ := c.GetInt("uid")
	suid, _ := c.GetInt("suid")
	pid, _ := c.GetInt("pid")
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 {
		templateRecords := logic.SearchProjectLeaderTemplateRecords(year, quarter, uid, suid, pid)
		list := make([]map[string]interface{}, len(templateRecords))
		for i, tr := range templateRecords {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			info["ScoreLimit"] = tr.Template.ScoreLimit
			if tr.Record != nil {
				info["Score"] = libs.Float64ToStringWithNoZero(tr.Record.Score)
			} else {
				info["Score"] = 0
			}
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人评分记发布后录详情列表
func (c *ApiInfterfaceController) GetProjectorScoreTypeAndScorePersonalBySumData() {
	uid, _ := c.GetInt("uid", 0)
	pid, _ := c.GetInt("pid")
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 {
		templateRecords := logic.SearchProjectLeaderTemplateAverageRecords(year, quarter, uid, pid)
		list := make([]map[string]interface{}, len(templateRecords))
		for i, tr := range templateRecords {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			info["ScoreLimit"] = tr.Template.ScoreLimit
			info["Score"] = libs.Float64ToStringWithNoZero(tr.Score)
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 项目负责人评分
type ProjectChargerScoreRequest struct {
	Year    int    `json:"year"`
	PID     int    `json:"pid"`
	Quarter int    `json:"quarter"`
	Scstr   string `json:"scstr"`
	SUID    int    `json:"suid"`
	UID     int    `json:"uid"`
}

// 给项目负责人评分
func (c *ApiInfterfaceController) ProjectChargerScore() {
	var request ProjectChargerScoreRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.UID
	suid := request.SUID
	year := request.Year
	quarter := request.Quarter
	pid := request.PID

	si := strings.Split(request.Scstr, ",")
	scoreInfo := make([]map[string]interface{}, 0)
	for _, s := range si {
		sii := strings.Split(s, "|")
		if len(sii) == 2 {
			info := make(map[string]interface{})
			info["ID"] = sii[0]
			info["Score"] = sii[1]
			scoreInfo = append(scoreInfo, info)
		}
	}

	if logic.CheckProjectLeaderReleaseReacord(year, quarter, uid, pid) {
		c.apiMsg("已发布无法修改", false, 500, nil)
	}

	if year > 0 && quarter > 0 {
		for _, info := range scoreInfo {
			id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)

			filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
			filter3 := models.DBFilter{Key: "tid", Value: id}          // 被打分人ID
			filter4 := models.DBFilter{Key: "uid", Value: uid}         // 被打分人ID
			filter5 := models.DBFilter{Key: "suid", Value: suid}       // 打分人ID
			filter6 := models.DBFilter{Key: "project_id", Value: pid}  // 项目ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

			records := models.SearchProjectLeaderScoreRecordsByFilters(filters...)

			if len(records) == 0 {
				record := new(models.ProjectLeaderScoreRecord)
				record.UID = uid
				record.SUID = suid
				record.TID = id
				record.Score = score
				record.Year = year
				record.Quarter = quarter
				record.ProjectID = pid
				models.AddProjectLeaderScoreRecord(record)
			} else {
				record := records[0]
				record.SUID = suid
				record.Score = score
				record.Update()
			}
		}

		err := logic.SaveProjectLeaderSumScoreRecord(year, quarter, uid, suid, pid)
		if err != nil {
			c.apiMsg(err.Error(), false, 500, nil)
		} else {
			c.apiMsg("操作成功", true, 200, nil)
		}

	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人列表(并赋上当前用户所评的分数)
func (c *ApiInfterfaceController) GetDepartmentors() {
	uid, _ := c.GetInt("uid", 0)
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 {
		departmentLeaders := models.SearchAllActiveQuarterDepartmentLeaders(year, quarter)

		list := make([]map[string]interface{}, len(departmentLeaders))
		for i, d := range departmentLeaders {
			department, _ := models.SearchDepartmentByID(d.DepartmentID)
			leader, _ := models.SearchUserByID(d.UID)
			info := make(map[string]interface{})
			info["DepartmentID"] = department.ID
			info["DepartmentName"] = department.Name
			info["UserName"] = leader.Name
			info["UserID"] = leader.ID

			filter1 := models.DBFilter{Key: "year", Value: year}                   // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter}             // 季度
			filter3 := models.DBFilter{Key: "uid", Value: leader.ID}               // 被打分人ID
			filter4 := models.DBFilter{Key: "suid", Value: uid}                    // 打分人ID
			filter5 := models.DBFilter{Key: "department_id", Value: department.ID} // 部门ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5}

			records := models.SearchDepartmentLeaderSumScoreRecordsByFilters(filters...)

			if len(records) == 1 {
				info["TScore"] = libs.Float64ToStringWithNoZero(records[0].Score)
			} else {
				info["TScore"] = 0
			}
			info["year"] = year
			info["year"] = quarter
			list[i] = info
		}

		c.apiMsg("操作成功", true, 200, list)
	}

	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人列表发布后数据
func (c *ApiInfterfaceController) GetDepartmentorsBySumData() {
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 {
		filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filters := []models.DBFilter{filter1, filter2}

		records := models.SearchDepartmentLeaderReleaseRecordsByOrder(filters...)

		list := make([]map[string]interface{}, len(records))

		for i, r := range records {
			info := make(map[string]interface{})
			user, _ := models.SearchUserByID(r.UID)
			department, _ := models.SearchDepartmentByID(r.DepartmentID)
			info["ID"] = i + 1
			info["DepartmentID"] = department.ID
			info["DepartmentName"] = department.Name
			info["UserName"] = user.Name
			info["UserID"] = user.ID
			info["TScore"] = libs.Float64ToStringWithNoZero(r.Score)
			info["Year"] = year
			info["Quarter"] = quarter
			list[i] = info
		}

		c.apiMsg("操作成功", true, 200, list)
	}

	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人评分记录详情列表
func (c *ApiInfterfaceController) GetDepartmentScoredTypeAndScore() {
	uid, _ := c.GetInt("uid")
	suid, _ := c.GetInt("suid")
	did, _ := c.GetInt("did")
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)
	if year > 0 && quarter > 0 {
		templateRecords := logic.SearchDepartmentLeaderScoreTemplateRecords(year, quarter, uid, suid, did)
		list := make([]map[string]interface{}, len(templateRecords))
		for i, tr := range templateRecords {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			info["ScoreLimit"] = tr.Template.ScoreLimit
			if tr.Record != nil {
				info["Score"] = libs.Float64ToStringWithNoZero(tr.Record.Score)
			} else {
				info["Score"] = 0
			}
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)

}

// 拉取当前季度可评分的部门负责人评分记录发布后详情列表
func (c *ApiInfterfaceController) GetDepartmentScoredTypeAndScoreBySumData() {
	uid, _ := c.GetInt("uid", 0)
	did, _ := c.GetInt("did")
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 {
		templateRecords := logic.SearchDepartmentLeaderTemplateAverageRecords(year, quarter, uid, did)
		list := make([]map[string]interface{}, len(templateRecords))
		for i, tr := range templateRecords {
			info := make(map[string]interface{})
			info["ID"] = tr.Template.ID
			info["Name"] = tr.Template.Name
			info["ScoreLimit"] = tr.Template.ScoreLimit
			info["Score"] = libs.Float64ToStringWithNoZero(tr.Score)
			list[i] = info
		}
		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 部门负责人评分
type DepartmentChargerScoreRequest struct {
	Year    int    `json:"year"`
	DID     int    `json:"did"`
	Quarter int    `json:"quarter"`
	Scstr   string `json:"scstr"`
	SUID    int    `json:"suid"`
	UID     int    `json:"uid"`
}

// 给部门负责人评分
func (c *ApiInfterfaceController) DepartmentChargerScore() {
	var request DepartmentChargerScoreRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.UID
	suid := request.SUID
	year := request.Year
	quarter := request.Quarter
	did := request.DID

	si := strings.Split(request.Scstr, ",")
	scoreInfo := make([]map[string]interface{}, 0)
	for _, s := range si {
		sii := strings.Split(s, "|")
		if len(sii) == 2 {
			info := make(map[string]interface{})
			info["ID"] = sii[0]
			info["Score"] = sii[1]
			scoreInfo = append(scoreInfo, info)
		}
	}

	if logic.CheckDepartmentLeaderReleaseReacord(year, quarter, uid, did) {
		c.apiMsg("已发布无法修改", false, 500, nil)
	}

	if year > 0 && quarter > 0 {
		for _, info := range scoreInfo {
			id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)

			filter1 := models.DBFilter{Key: "year", Value: year}         // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter}   // 季度
			filter3 := models.DBFilter{Key: "tid", Value: id}            // 被打分人ID
			filter4 := models.DBFilter{Key: "uid", Value: uid}           // 被打分人ID
			filter5 := models.DBFilter{Key: "suid", Value: suid}         // 打分人ID
			filter6 := models.DBFilter{Key: "department_id", Value: did} // 部门ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4, filter5, filter6}

			records := models.SearchDepartmentLeaderScoreRecordsByFilters(filters...)

			if len(records) == 0 {
				record := new(models.DepartmentLeaderScoreRecord)
				record.UID = uid
				record.SUID = suid
				record.TID = id
				record.Score = score
				record.Year = year
				record.Quarter = quarter
				record.DepartmentID = did
				models.AddDepartmentLeaderScoreRecord(record)
			} else {
				record := records[0]
				record.SUID = suid
				record.Score = score
				record.Update()
			}
		}

		err := logic.SaveDepartmentLeaderSumScoreRecord(year, quarter, uid, suid, did)
		if err != nil {
			c.apiMsg(err.Error(), false, 500, nil)
		} else {
			c.apiMsg("操作成功", true, 200, nil)
		}

	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 获取项目评分一级目录列表
func (c *ApiInfterfaceController) GetScoreType1() {
	template1 := models.SearchAllProjectTemplate1s()
	list := make([]map[string]interface{}, len(template1))
	for i, t1 := range template1 {
		info := make(map[string]interface{})
		info["ID"] = t1.ID
		info["Name"] = t1.Name
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 根据模版ID和季度查询项目评分排名
func (c *ApiInfterfaceController) GetProjectScoreType1() {
	t1id, _ := c.GetInt("t1id")
	year, _ := c.GetInt("year", 0)
	quarter, _ := c.GetInt("quarter", 0)

	if year > 0 && quarter > 0 {
		filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
		filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
		filters := []models.DBFilter{filter1, filter2}

		records := models.SearchProjectReleaseRecordsByFilters(filters...)
		fmt.Println(records)

		list := make([]map[string]string, 0)

		for _, r := range records {
			fmt.Println(r)
			col := make(map[string]string)
			project, _ := models.SearchProjectByID(r.PID)
			col["Name"] = project.Name
			col["ProjectUserID"] = strconv.Itoa(project.Leader)
			col["ProjectID"] = strconv.Itoa(r.PID)

			filter1 := models.DBFilter{Key: "year", Value: year}       // 年度
			filter2 := models.DBFilter{Key: "quarter", Value: quarter} // 季度
			filter3 := models.DBFilter{Key: "t1id", Value: t1id}       // 一级模版
			filter4 := models.DBFilter{Key: "pid", Value: r.PID}       // 项目ID
			filters := []models.DBFilter{filter1, filter2, filter3, filter4}

			recordIs := models.SearchProjectScoreRecord1sByFilters(filters...)
			if len(recordIs) > 0 {
				recordI := recordIs[0]
				col["TotalScore"] = libs.Float64ToStringWithNoZero(recordI.TotalScore)
			}

			template2Records := logic.SearchProjectTemplate2Records(year, quarter, t1id, r.PID)
			for i, tr := range template2Records {
				key := "score" + strconv.Itoa(i)
				if tr.Record != nil {
					col[key] = libs.Float64ToStringWithNoZero(tr.Record.TotalScore * tr.Template.Percentage)
				}
			}

			list = append(list, col)

		}

		slice.Sort(list[:], func(i, j int) bool {
			iScore, _ := strconv.ParseFloat(list[i]["TotalScore"], 64)
			jScore, _ := strconv.ParseFloat(list[j]["TotalScore"], 64)
			return iScore > jScore
		})

		for i, c := range list {
			c["ID"] = fmt.Sprintf("%d", i+1)
		}

		c.apiMsg("操作成功", true, 200, list)
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// Private Function

// api返回
func (c *ApiInfterfaceController) apiMsg(msg interface{}, status bool, code int, data interface{}) {
	out := make(map[string]interface{})
	out["Code"] = code
	out["Msg"] = msg
	out["Data"] = data
	out["Success"] = status
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}
