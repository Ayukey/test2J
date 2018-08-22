package models

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(Admin),
		new(Role),
		new(RoleAuth),
		new(Auth),
		new(PositionRoleInfo),
		new(UserInfo),
		new(DepartmentInfo),
		new(ProjectInfo),
		new(ProjectorInfo),
		new(ScoreTypeInfoI),
		new(ScoreTypeInfoII),
		new(ScoreTypeInfoIII),
		new(ScoreRecordInfoI),
		new(ScoreRecordInfoII),
		new(ScoreRecordInfoIII),
		new(STUserMappingI),
		new(STUserMappingII),
		new(STUserMappingIII),
		new(ProjectorScoreRecords),
		new(ProjectorScoreTypeInfo),
		new(ProjectorSumPubInfo),
		new(ProjectorScoreSumInfo),
		new(DepartmentorScoreRecords),
		new(DepartmentorScoreTypeInfo),
		new(DepartmentorSumPubInfo),
		new(DepartmentorScoreSumInfo),
		new(ProjectSumPubInfo))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
