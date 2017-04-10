package menu

import (
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"strings"
)

// TestAccountClear 测试账户清空数据库菜单
type TestAccountClear struct {
	SubItem
}

// InitMenu 初始化
func (ths *TestAccountClear) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Clear DB"
	return ths
}

func (ths *TestAccountClear) execute() {
	_L.LoggerInstance.Info(" ** Test account CLEAR ** \r\n")
	countInDb := _OP.GetTestAccountCount(nil)
	_L.LoggerInstance.InfoPrint(" >> Current test account count : %d \r\n", countInDb)
	_L.LoggerInstance.InfoPrint(" >> Are you confirm 'CLEAR' test account database? (yes/no) : ")
	input, b := ths.InputString()
	if b && strings.ToLower(input) == "yes" {
		err := _OP.ClearTestAccount(nil)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Clear test account database has error :\r\n %+v\r\n", err)
		}
		return
	}
	_L.LoggerInstance.InfoPrint(" >> 'CLEAR' test account database already be canceled!\r\n")
}
