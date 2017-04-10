package db

import "github.com/go-xorm/xorm"

const (
	KeyTestAccount = "TTestAccount"
	KeyStatics     = "TStatics"
	KeyMainAccount = "TMainAccount"

	QtGetRecord = iota + 1
	QtCheckRecord
	QtQuaryAllRecord
	QtAddRecord
	QtAddRecords
	QtUpdateRecord
	QtUpdateRecords
	QtDeleteRecord
	QtSeachRecord
	QtClearAllRecord
	QtGetCount
	QtGetCountRecords
	QtCommand

	MySqlDriver    DatabaseType = "mysql"
	SqliteDriver   DatabaseType = "sqlite3"
	PostgresDriver DatabaseType = "postgres"
)

// DatabaseType 数据库类型
type DatabaseType string

// DatabaseInfo 数据库基本信息定义
type DatabaseInfo struct {
	DbType    DatabaseType
	AliasName string
	Host      string
	Port      string
	UserName  string
	Password  string
	IsDebug   bool
}

// OperationInterface 数据库接口定义
type OperationInterface interface {
	Init(e *xorm.Engine)
	GetKey() string
	Query(qtype int, v ...interface{}) (interface{}, error)
	GetEngine() *xorm.Engine
}

// TMainAccount 主账户定义
type TMainAccount struct {
	ID         uint64 `xorm:"'id' pk autoincr"`
	Index      int
	AccountID  string `xorm:"'account_id' notnull unique"`
	SecertAddr string `xorm:"notnull unique"`
	Balance    float64
	Success    string `xorm:"varchar(1) notnull"`
}

// TTestAccount 测试账户定义
type TTestAccount struct {
	ID             uint64 `xorm:"'id' pk autoincr"`
	Index          int
	GroupIndex     int
	GroupItemIndex int
	AccountID      string `xorm:"'account_id' notnull unique"`
	SecertAddr     string `xorm:"notnull unique"`
	Balance        float64
	AssetBalance   float64
	AssetCode      string `xorm:"varchar(64)"`
	Success        string `xorm:"varchar(1) notnull"`
	InUse          string `xorm:"varchar(10)"`
}

// TStatic 统计信息定义
type TStatic struct {
	ID             uint64 `xorm:"'id' pk autoincr"`
	StaticName     string `xorm:"varchar(255)"`
	GroupIndex     int
	ItemIndex      int
	StartTime      int64
	CompleteTime   int64
	Success        string `xorm:"varchar(64) notnull"`
	FailureCause   string `xorm:"varchar(2048)"`
	Action         string `xorm:"varchar(64)"`
	Signature      string `xorm:"varchar(65535) notnull"`
	MainAccID      string `xorm:"'main_acc_id'"`
	OperationCnt   int
	TransactionCnt int
}
