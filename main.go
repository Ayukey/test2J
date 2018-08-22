package main

import (
	_ "jg2j_server/models"
	_ "jg2j_server/routers"

	"github.com/astaxie/beego"

	log "github.com/cihub/seelog"
)

func main() {
	logger, err := log.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
	beego.Run()

}
