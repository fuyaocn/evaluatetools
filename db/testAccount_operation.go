package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

// TestAccOprtion 主账户数据库操作
type TestAccOprtion struct {
	BaseOperation
	sqlcmd string
}

// Init 初始化
func (ths *TestAccOprtion) Init(e *xorm.Engine) {
	ths.BaseOperation.Init(e)
	ths.currKey = KeyTestAccount
	ths.tableName = "t_test_account"
	ths.sqlcmd = "UPDATE `t_test_account` SET `index`=?, `group_index`=?, `group_item_index`=?, `balance`=?, `asset_balance`=?, `asset_code`=?, `in_use`=?"
	ths.sqlcmd += " WHERE `account_id`=?"
}

// Query query exeute
func (ths *TestAccOprtion) Query(qtype int, v ...interface{}) (ret interface{}, err error) {
	switch qtype {
	case QtUpdateRecord:
		ths.locker.Lock()
		defer ths.locker.Unlock()
		testAcc := v[0].(*TTestAccount)
		return ths.updateData(testAcc)
	case QtUpdateRecords:
		ths.locker.Lock()
		defer ths.locker.Unlock()
		testAcc := v[0].([]*TTestAccount)
		return ths.updateDatas(testAcc)
	}
	return ths.BaseOperation.Query(qtype, v...)
}

func (ths *TestAccOprtion) updateDatas(src []*TTestAccount) (ret interface{}, err error) {
	for idx, itm := range src {
		ret, err = ths.updateData(itm)
		if err != nil {
			return ret, fmt.Errorf(" Update [%d] test account data to database has error : \r\n%+v", idx, err)
		}
	}
	return
}

func (ths *TestAccOprtion) updateData(itm *TTestAccount) (interface{}, error) {
	return ths.engine.Exec(ths.sqlcmd, itm.Index, itm.GroupIndex, itm.GroupItemIndex, itm.Balance, itm.AssetBalance, itm.AssetCode, itm.InUse, itm.AccountID)
}
