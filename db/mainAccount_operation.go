package db

import (
	"github.com/go-xorm/xorm"
)

// MainAccOprtion 主账户数据库操作
type MainAccOprtion struct {
	BaseOperation
}

// Init 初始化
func (ths *MainAccOprtion) Init(e *xorm.Engine) {
	ths.BaseOperation.Init(e)
	ths.currKey = KeyMainAccount
	ths.tableName = "t_main_account"
}

// Query query exeute
func (ths *MainAccOprtion) Query(qtype int, v ...interface{}) (ret interface{}, err error) {
	return ths.BaseOperation.Query(qtype, v...)
}
