package main

import (
	"fmt"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	_M "jojopoper/NBi/StressTest/menu"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_L.LoggerInstance = _L.NewLoggerInstance(fmt.Sprintf("NBi.StressTest.%s", time.Now().Format("2006-01-02_15.04.05.000")))
	_L.LoggerInstance.OpenDebug = true
	_L.LoggerInstance.SetLogFunCallDepth(4)
	_L.LoggerInstance.InfoPrint(" > Init Appconfig instance...\r\n")
	_AC.ConfigInstance = _AC.NewConfigController()
	dbConf := _AC.ConfigInstance.GetDateBaseConf()
	if dbConf == nil {
		panic(0)
	}
	_L.LoggerInstance.InfoPrint("Init database ...\r\n")
	_DB.DataBaseInstance = _DB.CreateDBInstance(dbConf)

	_M.MainMenuInstace = _M.CreateMenuInstance()
	_M.MainMenuInstace.ExecuteFunc()
}
