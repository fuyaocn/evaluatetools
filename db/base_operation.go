package db

import (
	"sync"

	"fmt"

	"github.com/go-xorm/xorm"
)

// BaseOperation 基础操作
type BaseOperation struct {
	engine    *xorm.Engine
	locker    *sync.Mutex
	currKey   string
	tableName string
}

// Init 初始化
func (ths *BaseOperation) Init(e *xorm.Engine) {
	ths.locker = &sync.Mutex{}
	ths.engine = e
}

// GetKey get key string
func (ths *BaseOperation) GetKey() string {
	return ths.currKey
}

// Query query exeute
func (ths *BaseOperation) Query(qtype int, v ...interface{}) (ret interface{}, err error) {
	switch qtype {
	case QtAddRecord: // 添加一条记录
		ths.locker.Lock()
		ret, err = ths.engine.InsertOne(v[0])
		ths.locker.Unlock()
	case QtAddRecords: // 添加很多数据
		ths.locker.Lock()
		ret, err = ths.engine.Insert(v...)
		ths.locker.Unlock()
	case QtClearAllRecord: // 清空数据表
		ths.locker.Lock()
		ret, err = ths.engine.Exec(fmt.Sprintf("TRUNCATE TABLE %s", ths.tableName))
		ths.locker.Unlock()
	case QtGetCount: // 得到指定条件的数据个数
		ths.locker.Lock()
		ret, err = ths.engine.Count(v[0])
		ths.locker.Unlock()
	}
	return
}

// GetEngine 获取数据库引擎，如果需要操作必须注意使用线程锁
func (ths *BaseOperation) GetEngine() *xorm.Engine {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	return ths.engine
}
