package menu

import (
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"strings"
)

// MainAccountClear 主账户清空数据库菜单
type MainAccountClear struct {
	SubItem
}

// InitMenu 初始化
func (ths *MainAccountClear) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Clear DB"
	return ths
}

func (ths *MainAccountClear) execute() {
	_L.LoggerInstance.Info(" ** Main account CLEAR ** \r\n")
	countInDb := _OP.GetMainAccountCount(nil)
	_L.LoggerInstance.InfoPrint(" >> Current main account count : %d \r\n", countInDb)
	_L.LoggerInstance.InfoPrint(" >> Are you confirm 'CLEAR' main account database? (yes/no) : ")
	input, b := ths.InputString()
	if b && strings.ToLower(input) == "yes" {
		err := _OP.ClearMainAccount(nil)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Clear main account database has error :\r\n %+v\r\n", err)
		}
		return
	}
	_L.LoggerInstance.InfoPrint(" >> 'CLEAR' main account database already be canceled!\r\n")
}
