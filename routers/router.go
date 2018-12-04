package routers

import (
	"jg2j_server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.AuthController{}, "*:Login")
	beego.Router("/logout", &controllers.AuthController{}, "*:Logout")

	beego.Router("/", &controllers.HomeController{}, "*:Index")
	beego.Router("/home", &controllers.HomeController{}, "*:Index")
	beego.AutoRouter(&controllers.HomeController{})

	beego.AutoRouter(&controllers.UserInfoController{})
	beego.AutoRouter(&controllers.DepartmentInfoController{})
	beego.AutoRouter(&controllers.ProjectInfoController{})

	beego.AutoRouter(&controllers.ScoreTypeInfoIController{})
	beego.AutoRouter(&controllers.ScoreTypeInfoIIController{})
	beego.AutoRouter(&controllers.ScoreTypeInfoIIIController{})
	beego.AutoRouter(&controllers.ProjectorScoreTypeInfoController{})
	beego.AutoRouter(&controllers.DepartmentorScoreTypeInfoController{})

	beego.AutoRouter(&controllers.RecordProjectScoreController{})
	beego.AutoRouter(&controllers.RecordProjectorScoreController{})
	beego.AutoRouter(&controllers.RecordDepartmentorScoreController{})

	beego.AutoRouter(&controllers.ProjectScoreRankController{})
	beego.AutoRouter(&controllers.ProjectorScoreRankController{})
	beego.AutoRouter(&controllers.DepartmentorScoreRankController{})

	beego.AutoRouter(&controllers.ReleaseProjectScoreController{})
	beego.AutoRouter(&controllers.ReleaseProjectorScoreController{})
	beego.AutoRouter(&controllers.ReleaseDepartmentorScoreController{})

	beego.AutoRouter(&controllers.ScoreRecordInfoIController{})
	beego.AutoRouter(&controllers.ScoreRecordInfoIIController{})
	beego.AutoRouter(&controllers.ScoreRecordInfoIIIController{})
	beego.AutoRouter(&controllers.ScoreDepartmentorRecordInfoController{})
	beego.AutoRouter(&controllers.ScoreProjectorRecordInfoController{})

	beego.AutoRouter(&controllers.STUserMappingIController{})
	beego.AutoRouter(&controllers.STUserMappingIIController{})

	// API

	beego.Router("/AppCommon/GetCurrentQuarter", &controllers.ApiInfterfaceController{}, "*:GetCurrentQuarter")
	beego.Router("/UserInfo/UserLogin", &controllers.ApiInfterfaceController{}, "*:UserLogin")
	beego.Router("/UserInfo/ModifyPwd", &controllers.ApiInfterfaceController{}, "*:ModifyPwd")
	beego.Router("/Projects/GetCanScoreProjects", &controllers.ApiInfterfaceController{}, "*:GetCanScoreProjects")
	beego.Router("/Projects/GetSTMapping1", &controllers.ApiInfterfaceController{}, "*:GetSTMapping1")
	beego.Router("/Projects/GetSTMapping2", &controllers.ApiInfterfaceController{}, "*:GetSTMapping2")
	beego.Router("/Projects/GetProjectScoreType1AndScore", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType1AndScore")
	beego.Router("/Projects/GetProjectScoreType2", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType2")
	beego.Router("/Projects/GetProjectScoreType3", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType3")
	beego.Router("/Projects/ProjectScore", &controllers.ApiInfterfaceController{}, "*:ProjectScore")
	beego.Router("/Projects/ProjectRemark2", &controllers.ApiInfterfaceController{}, "*:ProjectRemark2")

	beego.Router("/Projects/GetProjectors", &controllers.ApiInfterfaceController{}, "*:GetProjectors")
	beego.Router("/Projects/GetProjectorScoreTypeAndScorePersonal", &controllers.ApiInfterfaceController{}, "*:GetProjectorScoreTypeAndScorePersonal")
	beego.Router("/Projects/ProjectChargerScore", &controllers.ApiInfterfaceController{}, "*:ProjectChargerScore")

	beego.Router("/Departments/GetDepartmentors", &controllers.ApiInfterfaceController{}, "*:GetDepartmentors")
	beego.Router("/Departments/GetDepartmentScoredTypeAndScore", &controllers.ApiInfterfaceController{}, "*:GetDepartmentScoredTypeAndScore")
	beego.Router("/Departments/DepartmentChargerScore", &controllers.ApiInfterfaceController{}, "*:DepartmentChargerScore")

	beego.Router("/Projects/GetScoreType1", &controllers.ApiInfterfaceController{}, "*:GetScoreType1")
	beego.Router("/ProjectSumData/GetProjectScoreType1", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType1")
	beego.Router("/ProjectSumData/GetProjectScoreType2", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType2")
	beego.Router("/ProjectSumData/GetProjectScoreType3", &controllers.ApiInfterfaceController{}, "*:GetProjectScoreType3")

	beego.Router("/ProjectSumData/GetProjectors", &controllers.ApiInfterfaceController{}, "*:GetProjectorsBySumData")
	beego.Router("/ProjectSumData/GetProjectorScoreTypeAndScorePersonal", &controllers.ApiInfterfaceController{}, "*:GetProjectorScoreTypeAndScorePersonalBySumData")

	beego.Router("/DepartmentSumData/GetDepartmentors", &controllers.ApiInfterfaceController{}, "*:GetDepartmentorsBySumData")
	beego.Router("/DepartmentSumData/GetDepartmentScoredTypeAndScore", &controllers.ApiInfterfaceController{}, "*:GetDepartmentScoredTypeAndScoreBySumData")
}
