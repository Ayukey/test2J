package models

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type DBFilter struct {
	Key   string
	Value interface{}
}

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
	orm.RegisterModelWithPrefix(
		beego.AppConfig.String("db.prefix"),
		new(Admin),
		new(Role),
		new(Auth),

		new(User),
		new(Department),
		new(Position),
		new(Project),

		new(ProjectTemplate1),
		new(ProjectTemplate2),
		new(ProjectTemplate3),

		new(ProjectScoreRecord1),
		new(ProjectScoreRecord2),
		new(ProjectScoreRecord3),

		new(Pt1Permission),
		new(Pt2Permission),

		new(ProjectReleaseRecord),

		new(DepartmentLeaderTemplate),
		new(DepartmentLeaderScoreRecord),
		new(DepartmentLeaderSumScoreRecord),
		new(DepartmentLeaderReleaseRecord),

		new(ProjectLeaderTemplate),
		new(ProjectLeaderScoreRecord),
		new(ProjectLeaderSumScoreRecord),
		new(ProjectLeaderReleaseRecord),

		new(QuarterActive),
		new(QuarterActiveDepartmentLeader),
		new(QuarterActiveProject),
		new(QuarterActiveProjectLeader),
	)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	orm.RunSyncdb("default", false, true)
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
