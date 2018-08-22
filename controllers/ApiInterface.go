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
	currentQuarter := getCurrentYearAndQuarter()
	c.apiMsg("操作成功", true, 200, currentQuarter)
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

	user, err := models.SearchUserInfoByAccount(request.UserName)

	if err != nil {
		c.apiMsg("账号不存在", false, 500, nil)
	}

	if user.Password != request.UserPwd {
		c.apiMsg("密码错误", false, 500, nil)
	}

	info := make(map[string]interface{})
	info["ID"] = user.ID
	info["DepartmentID"] = user.DepartmentID
	info["ProjectID"] = user.ProjectID
	info["Name"] = user.Name
	info["RoleID"] = user.RoleID
	info["UserName"] = user.Account
	info["UserPwd"] = user.Password
	info["UToken"] = ""
	c.apiMsg("操作成功", true, 200, info)
}

// 修改密码请求
type ModifyPwdRequest struct {
	UserName string `json:"UserName"`
	UserPwd  string `json:"UserPwd"`
	NewPwd   string `json:"NewPwd"`
}

// 用户修改密码
func (c *ApiInfterfaceController) ModifyPwd() {
	var request ModifyPwdRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	user, err := models.SearchUserInfoByAccount(request.UserName)
	if err != nil {
		c.apiMsg("账号不存在", false, 500, nil)
	}

	if user.Password != request.UserPwd {
		c.apiMsg("原始密码错误", false, 500, nil)
	}

	user.Password = request.NewPwd

	err = user.Update()
	if err != nil {
		c.apiMsg("操作失败", false, 500, nil)
	}

	c.apiMsg("操作成功", true, 200, nil)
}

// 拉取当前季度可评分的项目列表
func (c *ApiInfterfaceController) GetCanScoreProjects() {
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			projects := models.SearchAllProjectsInUse(year, quarter)
			list := make([]map[string]interface{}, len(projects))
			for i, p := range projects {
				info := make(map[string]interface{})
				info["ID"] = p.ID
				info["Name"] = p.Name
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取用户项目一级评分项权限数据
func (c *ApiInfterfaceController) GetSTMapping1() {
	UserID, _ := c.GetInt("UserID")
	maps := models.SearchSTUserMappingIByUID(UserID)
	list := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		info := make(map[string]interface{})
		info["TID"] = m.TID
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 拉取用户项目二级评分项权限数据
func (c *ApiInfterfaceController) GetSTMapping2() {
	UserID, _ := c.GetInt("UserID")
	maps := models.SearchSTUserMappingIIByUID(UserID)
	list := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		info := make(map[string]interface{})
		info["TID"] = m.TID
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 拉取项目一级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType1AndScore() {
	proid, _ := c.GetInt("proid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			infos := models.SerachScoreTypeRecordInfoIList(year, quarter, proid)
			list := make([]map[string]interface{}, len(infos))
			for i, f := range infos {
				info := make(map[string]interface{})
				info["ID"] = f.Type.ID
				info["Name"] = f.Type.Name
				fmt.Println(f.Record)
				if f.Record == nil {
					info["TotalScore"] = 0
				} else {
					info["TotalScore"] = f.Record.TotalScore //libs.Float64ToStringWithNoZero(f.Record.TotalScore)
				}
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取项目二级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType2() {
	tid, _ := c.GetInt("tid")
	proid, _ := c.GetInt("proid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			infos := models.SerachScoreTypeRecordInfoIIList(year, quarter, tid, proid)
			list := make([]map[string]interface{}, len(infos))
			for i, f := range infos {
				info := make(map[string]interface{})
				info["ID"] = f.Type.ID
				info["Name"] = f.Type.Name
				fmt.Println(f.Record)
				if f.Record == nil {
					info["TotalScore"] = 0
				} else {
					info["TotalScore"] = f.Record.TotalScore //libs.Float64ToStringWithNoZero(f.Record.TotalScore)
				}
				list[i] = info
			}
			da := make(map[string]interface{})
			da["Remark"] = ""
			da["Scores"] = list
			c.apiMsg("操作成功", true, 200, da)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取项目三级评分项数据
func (c *ApiInfterfaceController) GetProjectScoreType3() {
	tid, _ := c.GetInt("tid")
	proid, _ := c.GetInt("proid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			infos := models.SerachScoreTypeRecordInfoIIIList(year, quarter, tid, 0, proid)
			list := make([]map[string]interface{}, len(infos))
			for i, f := range infos {
				info := make(map[string]interface{})
				info["ID"] = f.Type.ID
				info["Name"] = f.Type.Name
				info["MaxScore"] = f.Type.MaxScore
				if f.Record == nil {
					info["TotalScore"] = 0
				} else {
					info["TotalScore"] = f.Record.Score //libs.Float64ToStringWithNoZero(f.Record.Score)
				}
				list[i] = info
			}
			da := make(map[string]interface{})
			typeII, _ := models.SearchScoreTypeInfoIIByID(tid)
			if typeII != nil {
				typeRecordInfoIInfos := models.SerachScoreTypeRecordInfoIIList(year, quarter, typeII.TID, proid)
				for _, typeRecord := range typeRecordInfoIInfos {
					if typeRecord.Type.ID == tid {
						if typeRecord.Record != nil {
							da["Remark"] = typeRecord.Record.Remark
						}
					}
				}
			}
			da["Scores"] = list
			c.apiMsg("操作成功", true, 200, da)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 项目三级评分项打分 同时更新相关的一、二级评分项分数
type ProjectScoreRequest struct {
	Proid  int    `json:"proid"`
	QT     string `json:"qt"`
	Scstr  string `json:"scstr"`
	Tid    int    `json:"tid"`
	Userid int    `json:"userid"`
}

func (c *ApiInfterfaceController) ProjectScore() {
	var request ProjectScoreRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.Userid
	tid := request.Tid
	pid := request.Proid
	year := 0
	quarter := 0
	yq := strings.Split(request.QT, "-")
	if len(yq) == 2 {
		year, _ = strconv.Atoi(yq[0])
		quarter, _ = strconv.Atoi(yq[1])
	}

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

	if year > 0 && quarter > 0 {
		for _, info := range scoreInfo {
			id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)
			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			filters = append(filters, "tid", tid)
			filters = append(filters, "project_id", pid)
			filters = append(filters, "scoretype_id", id)
			recordIIIList := models.SearchScoreRecordInfoIIIByFilters(filters...)

			if len(recordIIIList) == 0 {
				recordIII := new(models.ScoreRecordInfoIII)
				recordIII.UserID = uid
				recordIII.ProjectID = pid
				recordIII.ScoreTypeID = id
				recordIII.Score = score
				recordIII.Year = year
				recordIII.Quarter = quarter
				recordIII.TID = tid
				recordIII.UpdateDate = float64(time.Now().Unix())
				recordIII.Remark = ""
				models.AddScoreRecordInfoIII(recordIII)
			} else {
				recordIII := recordIIIList[0]
				recordIII.UserID = uid
				recordIII.Score = score
				recordIII.UpdateDate = float64(time.Now().Unix())
				recordIII.Update()
			}
		}

		// 更新二级评分数据

		typeII, _ := models.SearchScoreTypeInfoIIByID(tid)

		typeRecordIIIList := models.SerachScoreTypeRecordInfoIIIList(year, quarter, tid, typeII.TID, pid)
		realScoreII := 0.0
		maxScoreII := 0.0
		for _, v := range typeRecordIIIList {
			maxScoreII += v.Type.MaxScore
			if v.Record != nil {
				if v.Record.Score == -1 {
					maxScoreII -= v.Type.MaxScore
				} else {
					realScoreII += v.Record.Score
				}
			}
		}
		totalScoreII := realScoreII / maxScoreII * 100

		if typeII.Percentage == 1 {
			totalScoreII = realScoreII
		}

		filters := make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "tid", typeII.TID)
		filters = append(filters, "project_id", pid)
		filters = append(filters, "scoretype_id", tid)
		recordIIList := models.SearchScoreRecordInfoIIByFilters(filters...)

		if len(recordIIList) == 0 {
			recordII := new(models.ScoreRecordInfoII)
			recordII.ProjectID = pid
			recordII.ScoreTypeID = tid

			recordII.TotalScore = totalScoreII
			recordII.Year = year
			recordII.Quarter = quarter
			recordII.TID = typeII.TID
			recordII.UpdateDate = float64(time.Now().Unix())
			models.AddScoreRecordInfoII(recordII)
		} else {
			recordII := recordIIList[0]
			recordII.TotalScore = totalScoreII
			recordII.UpdateDate = float64(time.Now().Unix())
			recordII.Update()
		}

		// 更新一级评分数据
		typeRecordIIList := models.SerachScoreTypeRecordInfoIIList(year, quarter, typeII.TID, pid)
		totalScoreI := 0.0
		for _, v := range typeRecordIIList {
			if v.Record != nil {
				if v.Record != nil {
					totalScoreI += v.Type.Percentage * v.Record.TotalScore
				}
			}
		}

		filters = make([]interface{}, 0)
		filters = append(filters, "year", year)
		filters = append(filters, "quarter", quarter)
		filters = append(filters, "project_id", pid)
		filters = append(filters, "scoretype_id", typeII.TID)
		recordIList := models.SearchScoreRecordInfoIByFilters(filters...)
		if len(recordIList) == 0 {
			recordI := new(models.ScoreRecordInfoI)
			recordI.ProjectID = pid
			recordI.ScoreTypeID = typeII.TID
			recordI.TotalScore = totalScoreI
			recordI.Year = year
			recordI.Quarter = quarter
			recordI.UpdateDate = float64(time.Now().Unix())
			models.AddScoreRecordInfoI(recordI)
		} else {
			recordI := recordIList[0]
			recordI.TotalScore = totalScoreI
			recordI.UpdateDate = float64(time.Now().Unix())
			recordI.Update()
		}
		c.apiMsg("操作成功", true, 200, nil)
	}
	c.apiMsg("数据异常", false, 500, nil)
}

// 更新或修改二级评分的总结
type ProjectRemark2Request struct {
	Proid  int    `json:"proid"`
	QT     string `json:"qt"`
	Remark string `json:"remark"`
	Tid    int    `json:"tid"`
}

func (c *ApiInfterfaceController) ProjectRemark2() {
	var request ProjectRemark2Request
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	t2id := request.Tid      // 二级评分模版ID
	pid := request.Proid     // 项目ID
	remark := request.Remark // 总结
	year := 0
	quarter := 0
	yq := strings.Split(request.QT, "-")
	if len(yq) == 2 {
		year, _ = strconv.Atoi(yq[0])
		quarter, _ = strconv.Atoi(yq[1])
		// 查询具体二级评分模版信息
		t2, _ := models.SearchScoreTypeInfoIIByID(t2id)
		if t2 != nil {
			t1id := t2.TID
			tr2s := models.SerachScoreTypeRecordInfoIIList(year, quarter, t1id, pid)
			for _, tr2 := range tr2s {
				if tr2.Type.ID == t2id {
					if tr2.Record != nil {
						tr2.Record.Remark = remark
						tr2.Record.Update()
					} else {
						fmt.Println("木有找到记录")
						r2 := new(models.ScoreRecordInfoII)
						r2.ProjectID = pid
						r2.ScoreTypeID = t2id
						r2.Year = year
						r2.Quarter = quarter
						r2.TID = t1id
						r2.Remark = remark
						r2.UpdateDate = float64(time.Now().Unix())
						models.AddScoreRecordInfoII(r2)
					}
				}
			}
		}
	}
	c.apiMsg("操作成功", true, 200, nil)
}

// 拉取当前季度可评分的项目负责人列表(并赋上当前用户所评的分数)
func (c *ApiInfterfaceController) GetProjectors() {
	userid, _ := c.GetInt("userid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			projectors := models.SearchAllProjectLeadersInUseWithScore(year, quarter, userid)
			list := make([]map[string]interface{}, len(projectors))
			for i, p := range projectors {
				info := make(map[string]interface{})
				info["ProjectID"] = p.Project.ID
				info["ProjectName"] = p.Project.Name
				info["UserName"] = p.UserInfo.Name
				info["UserID"] = p.UserInfo.ID
				if p.Score != nil {
					info["TScore"] = libs.Float64ToStringWithNoZero(p.Score.Score)
				} else {
					info["TScore"] = 0
				}
				info["Qt"] = yearAndQuarter
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人列表发布后数据
func (c *ApiInfterfaceController) GetProjectorsBySumData() {
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			records := models.SearchProjectorSumPubInfoByOrder(filters...)

			list := make([]map[string]interface{}, len(records))

			for i, r := range records {
				info := make(map[string]interface{})
				user, _ := models.SearchUserInfoByID(r.UserID)
				project, _ := models.SearchProjectInfoByID(r.ProjectID)
				info["ID"] = i + 1
				info["ProjectID"] = project.ID
				info["ProjectName"] = project.Name
				info["UserName"] = user.Name
				info["UserID"] = user.ID
				info["TScore"] = libs.Float64ToStringWithNoZero(r.Score)
				info["Qt"] = yearAndQuarter
				list[i] = info
			}

			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人评分记录详情列表
func (c *ApiInfterfaceController) GetProjectorScoreTypeAndScorePersonal() {
	userid, _ := c.GetInt("userid")
	suserid, _ := c.GetInt("suserid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			typeRecords := logic.SearchProjectorScoreTypeRecordInfos(year, quarter, userid, suserid)
			list := make([]map[string]interface{}, len(typeRecords))
			for i, tr := range typeRecords {
				info := make(map[string]interface{})
				info["ID"] = tr.Type.ID
				info["Name"] = tr.Type.Name
				info["ScoreLimit"] = tr.Type.ScoreLimit
				if tr.Record != nil {
					info["Score"] = libs.Float64ToStringWithNoZero(tr.Record.Score)
				} else {
					info["Score"] = 0
				}
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的项目负责人评分记发布后录详情列表
func (c *ApiInfterfaceController) GetProjectorScoreTypeAndScorePersonalBySumData() {
	userid, _ := c.GetInt("userid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			typeRecords := logic.SearchProjectorScoreTypeRecordInfosBySumData(year, quarter, userid)
			list := make([]map[string]interface{}, len(typeRecords))
			for i, tr := range typeRecords {
				info := make(map[string]interface{})
				info["ID"] = tr.Type.ID
				info["Name"] = tr.Type.Name
				info["ScoreLimit"] = tr.Type.ScoreLimit
				info["Score"] = libs.Float64ToStringWithNoZero(tr.Score)
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 项目负责人评分
type ProjectChargerScoreRequest struct {
	QT      string `json:"qt"`
	Scstr   string `json:"scstr"`
	SUserid int    `json:"suserid"`
	Userid  int    `json:"userid"`
}

// 给项目负责人评分
func (c *ApiInfterfaceController) ProjectChargerScore() {
	var request ProjectChargerScoreRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.Userid
	suid := request.SUserid
	year := 0
	quarter := 0
	yq := strings.Split(request.QT, "-")
	if len(yq) == 2 {
		year, _ = strconv.Atoi(yq[0])
		quarter, _ = strconv.Atoi(yq[1])
	}

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

	if logic.CheckProjectorScoreReleaseStatus(year, quarter, uid) {
		c.apiMsg("已发布无法修改", false, 500, nil)
	}

	user, _ := models.SearchUserInfoByID(uid)

	if year > 0 && quarter > 0 {
		for _, info := range scoreInfo {
			id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)

			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			filters = append(filters, "tid", id)
			filters = append(filters, "user_id", uid)
			filters = append(filters, "scoreuser_id", suid)
			filters = append(filters, "project_id", user.ProjectID)
			recordList := models.SearchProjectorScoreRecordsByFilters(filters...)

			if len(recordList) == 0 {
				record := new(models.ProjectorScoreRecords)
				record.UserID = uid
				record.ScoreUserID = suid
				record.TID = id
				record.Score = score
				record.Year = year
				record.Quarter = quarter
				record.ProjectID = user.ProjectID
				models.AddProjectorScoreRecords(record)
			} else {
				record := recordList[0]
				record.ScoreUserID = suid
				record.Score = score
				record.Update()
			}
		}

		err := logic.SaveProjectorScoreBySingleDepartmentor(year, quarter, uid, suid)
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
	userid, _ := c.GetInt("userid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			departmentors := models.SearchAllDepartmentLeadersInUseWithScore(year, quarter, userid)
			list := make([]map[string]interface{}, len(departmentors))
			for i, p := range departmentors {
				info := make(map[string]interface{})
				info["DepartmentID"] = p.Department.ID
				info["DepartmentName"] = p.Department.Name
				info["UserName"] = p.UserInfo.Name
				info["UserID"] = p.UserInfo.ID
				if p.Score != nil {
					info["TScore"] = libs.Float64ToStringWithNoZero(p.Score.Score)
				} else {
					info["TScore"] = 0
				}
				info["Qt"] = yearAndQuarter
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人列表发布后数据
func (c *ApiInfterfaceController) GetDepartmentorsBySumData() {
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			records := models.SearchDepartmentorSumPubInfoByOrder(filters...)

			list := make([]map[string]interface{}, len(records))

			for i, r := range records {
				info := make(map[string]interface{})
				user, _ := models.SearchUserInfoByID(r.UserID)
				department, _ := models.SearchDepartmentInfoByID(r.DepartmentID)
				info["ID"] = i + 1
				info["DepartmentID"] = department.ID
				info["DepartmentName"] = department.Name
				info["UserName"] = user.Name
				info["UserID"] = user.ID
				info["TScore"] = libs.Float64ToStringWithNoZero(r.Score)
				info["Qt"] = yearAndQuarter
				list[i] = info
			}

			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人评分记录详情列表
func (c *ApiInfterfaceController) GetDepartmentScoredTypeAndScore() {
	userid, _ := c.GetInt("userid")
	suserid, _ := c.GetInt("suserid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			typeRecords := logic.SearchDepartmentorScoreTypeRecordInfos(year, quarter, userid, suserid)
			list := make([]map[string]interface{}, len(typeRecords))
			for i, tr := range typeRecords {
				info := make(map[string]interface{})
				info["ID"] = tr.Type.ID
				info["Name"] = tr.Type.Name
				info["ScoreLimit"] = tr.Type.ScoreLimit
				if tr.Record != nil {
					info["Score"] = libs.Float64ToStringWithNoZero(tr.Record.Score)
				} else {
					info["Score"] = 0
				}
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 拉取当前季度可评分的部门负责人评分记录发布后详情列表
func (c *ApiInfterfaceController) GetDepartmentScoredTypeAndScoreBySumData() {
	userid, _ := c.GetInt("userid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			typeRecords := logic.SearchDepartmentorScoreTypeRecordInfosBySumData(year, quarter, userid)
			list := make([]map[string]interface{}, len(typeRecords))
			for i, tr := range typeRecords {
				info := make(map[string]interface{})
				info["ID"] = tr.Type.ID
				info["Name"] = tr.Type.Name
				info["ScoreLimit"] = tr.Type.ScoreLimit
				info["Score"] = libs.Float64ToStringWithNoZero(tr.Score)
				list[i] = info
			}
			c.apiMsg("操作成功", true, 200, list)
		}
	}
	c.apiMsg("季度数据异常", false, 500, nil)
}

// 部门负责人评分
type DepartmentChargerScoreRequest struct {
	QT      string `json:"qt"`
	Scstr   string `json:"scstr"`
	SUserid int    `json:"suserid"`
	Userid  int    `json:"userid"`
}

// 给部门负责人评分
func (c *ApiInfterfaceController) DepartmentChargerScore() {
	var request DepartmentChargerScoreRequest
	fmt.Println(string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	uid := request.Userid
	suid := request.SUserid
	year := 0
	quarter := 0
	yq := strings.Split(request.QT, "-")
	if len(yq) == 2 {
		year, _ = strconv.Atoi(yq[0])
		quarter, _ = strconv.Atoi(yq[1])
	}

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

	user, _ := models.SearchUserInfoByID(uid)

	if logic.CheckDepartmentorScoreReleaseStatus(year, quarter, uid) {
		c.apiMsg("已发布无法修改", false, 500, nil)
	}

	fmt.Println(user)
	if year > 0 && quarter > 0 {
		for _, info := range scoreInfo {
			id, _ := strconv.Atoi(info["ID"].(string))
			score, _ := strconv.ParseFloat(info["Score"].(string), 64)
			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			filters = append(filters, "user_id", uid)
			filters = append(filters, "tid", id)
			filters = append(filters, "scoreuser_id", suid)
			recordList := models.SearchDepartmentorScoreRecordsByFilters(filters...)

			if len(recordList) == 0 {
				record := new(models.DepartmentorScoreRecords)
				record.UserID = uid
				record.ScoreUserID = suid
				record.TID = id
				record.Score = score
				record.Year = year
				record.Quarter = quarter
				models.AddDepartmentorScoreRecords(record)
			} else {
				record := recordList[0]
				record.ScoreUserID = suid
				record.Score = score
				record.Update()
			}
		}

		err := logic.SaveDepartmentorScoreBySingleProjector(year, quarter, uid, suid)
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
	typeInfo := models.SearchAllScoreTypeInfoIList()
	list := make([]map[string]interface{}, len(typeInfo))
	for i, m := range typeInfo {
		info := make(map[string]interface{})
		info["ID"] = m.ID
		info["Name"] = m.Name
		list[i] = info
	}
	c.apiMsg("操作成功", true, 200, list)
}

// 根据模版ID和季度查询项目评分排名
func (c *ApiInfterfaceController) GetProjectScoreType1() {
	tid, _ := c.GetInt("tid")
	yearAndQuarter := c.GetString("qt")
	list := strings.Split(yearAndQuarter, "-")
	if len(list) == 2 {
		year, _ := strconv.Atoi(list[0])
		quarter, _ := strconv.Atoi(list[1])
		if year > 0 && quarter > 0 {
			filters := make([]interface{}, 0)
			filters = append(filters, "year", year)
			filters = append(filters, "quarter", quarter)
			projectList := models.SearchProjectSumPubInfoByFilters(filters...)

			list := make([]map[string]string, 0)

			for _, p := range projectList {
				col := make(map[string]string)
				project, _ := models.SearchProjectInfoByID(p.ProjectID)
				col["Name"] = project.Name
				col["ProjectUserID"] = strconv.Itoa(project.UserID)

				filters := make([]interface{}, 0)
				filters = append(filters, "year", year)
				filters = append(filters, "quarter", quarter)
				filters = append(filters, "scoretype_id", tid)
				filters = append(filters, "project_id", p.ProjectID)

				recordIList := models.SearchScoreRecordInfoIByFilters(filters...)
				if len(recordIList) > 0 {
					recordI := recordIList[0]
					col["TotalScore"] = libs.Float64ToStringWithNoZero(recordI.TotalScore)
				}

				recordIIList := models.SerachScoreTypeRecordInfoIIList(year, quarter, tid, p.ProjectID)
				for i, r := range recordIIList {
					key := "score" + strconv.Itoa(i)
					if r.Record != nil {
						col[key] = libs.Float64ToStringWithNoZero(r.Record.TotalScore * r.Type.Percentage)
					}
				}

				list = append(list, col)

			}

			slice.Sort(list[:], func(i, j int) bool {
				iScore, _ := strconv.ParseFloat(list[i]["totalscore"], 64)
				jScore, _ := strconv.ParseFloat(list[j]["totalscore"], 64)
				return iScore > jScore
			})

			for i, c := range list {
				c["ID"] = fmt.Sprintf("%d", i+1)
			}

			c.apiMsg("操作成功", true, 200, list)
		}
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
