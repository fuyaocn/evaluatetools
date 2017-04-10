package menu

import (
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
)

// MainAccountInfo 主账户信息菜单
type MainAccountInfo struct {
	SubItem
	CountInDb int64
}

// InitMenu 初始化
func (ths *MainAccountInfo) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Information"
	return ths
}

func (ths *MainAccountInfo) execute() {
	_L.LoggerInstance.Info(" ** Main account information ** \r\n")
	ths.CountInDb = _OP.GetMainAccountCount(nil)
	maxIndex := _OP.GetMainAccountMaxIndex(nil)
	_L.LoggerInstance.InfoPrint(" >> Current main account count : %d \r\n", ths.CountInDb)
	_L.LoggerInstance.InfoPrint(" >> Current main account last index : %d \r\n", maxIndex)
	ths.Wait()
}
