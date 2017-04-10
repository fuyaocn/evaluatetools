package db

import (
	"github.com/go-xorm/xorm"
)

// StaticsOprtion 统计结果数据库操作
type StaticsOprtion struct {
	BaseOperation
}

// Init 初始化
func (ths *StaticsOprtion) Init(e *xorm.Engine) {
	ths.BaseOperation.Init(e)
	ths.currKey = KeyStatics
	ths.tableName = "t_statics"
}

// Query query exeute
func (ths *StaticsOprtion) Query(qtype int, v ...interface{}) (ret interface{}, err error) {
	return ths.BaseOperation.Query(qtype, v...)
}
