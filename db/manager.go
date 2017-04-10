package db

import (
	"fmt"
	_L "jojopoper/NBi/StressTest/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

//DataBaseInstance database instance
var DataBaseInstance *Manager

// Manager 数据库管理
type Manager struct {
	dbInfo     *DatabaseInfo
	dbEngine   *xorm.Engine
	operations map[string]OperationInterface
}

// CreateDBInstance 创建实例
func CreateDBInstance(dbConfig *DatabaseInfo) *Manager {
	ret := &Manager{}
	ret.Init(dbConfig)
	return ret
}

// Init 初始化数据库
func (ths *Manager) Init(dbConfig *DatabaseInfo) {
	_L.LoggerInstance.InfoPrint("[Manager:InitDB] Init database begin\r\n")
	ths.dbInfo = dbConfig

	ths.initEngine()

	// 注册Orm的数据库表
	ths.ormRegModels()

	ths.initOperation()
	_L.LoggerInstance.InfoPrint("[Manager:InitDB] Init database success\r\n")
}

func (ths *Manager) initEngine() {
	ths.dbEngine = nil
	switch ths.dbInfo.DbType {
	case MySqlDriver:
		ths.dbEngine = ths.getMySQLEngine()
	case PostgresDriver:
		ths.dbEngine = ths.getPostgresEngine()
	}
	if ths.dbEngine == nil {
		_L.LoggerInstance.ErrorPrint("[Manager:initEngine] Undefined db type = %s\r\n", ths.dbInfo.DbType)
		panic(1)
	}
	ths.dbEngine.ShowErr = true
	ths.dbEngine.ShowWarn = true

	if ths.dbInfo.IsDebug {
		ths.dbEngine.ShowDebug = true
		ths.dbEngine.ShowInfo = true
		ths.dbEngine.ShowSQL = true
	}
}

func (ths *Manager) getMySQLEngine() *xorm.Engine {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", //Asia%2FShanghai
		ths.dbInfo.UserName, ths.dbInfo.Password, ths.dbInfo.Host, ths.dbInfo.Port, ths.dbInfo.AliasName)
	ret, err := xorm.NewEngine(string(MySqlDriver), dataSourceName)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[Manager:getMySqlEngine] Create MySql has error! \r\n\t%v\r\n", err)
		return nil
	}
	err = ret.Ping()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[Manager:getMySqlEngine] Create MySql Ping error! \r\n\t%v\r\n", err)
		return nil
	}
	return ret
}

func (ths *Manager) getPostgresEngine() *xorm.Engine {
	dataSourceName := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		ths.dbInfo.AliasName, ths.dbInfo.UserName, ths.dbInfo.Password, ths.dbInfo.Host, ths.dbInfo.Port)
	ret, err := xorm.NewEngine(string(PostgresDriver), dataSourceName)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[Manager:getPostgresEngine] Create Postgres has error! \r\n\t%v\r\n", err)
		return nil
	}
	return ret
}

// ormRegModels 初始化数据库表
func (ths *Manager) ormRegModels() {
	err := ths.dbEngine.Sync(new(TMainAccount), new(TTestAccount), new(TStatic))
	if err != nil {
		_L.LoggerInstance.InfoPrint("[Manager:ormRegModels] XORM Engine Sync is err %v\r\n", err)
		panic(1)
	}
}

func (ths *Manager) initOperation() {
	if ths.operations == nil {
		ths.operations = make(map[string]OperationInterface)
	}

	mainAccOpt := &MainAccOprtion{}
	mainAccOpt.Init(ths.dbEngine)
	ths.operations[mainAccOpt.GetKey()] = mainAccOpt

	testAccOpt := &TestAccOprtion{}
	testAccOpt.Init(ths.dbEngine)
	ths.operations[testAccOpt.GetKey()] = testAccOpt

	stat := &StaticsOprtion{}
	stat.Init(ths.dbEngine)
	ths.operations[stat.GetKey()] = stat
}

// GetOperation 得到对应的操作数据库控制器
func (ths *Manager) GetOperation(key string) OperationInterface {
	return ths.operations[key]
}
