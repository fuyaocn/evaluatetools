package menu

import (
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
)

// TestAccountInfo 主账户信息菜单
type TestAccountInfo struct {
	SubItem
}

// InitMenu 初始化
func (ths *TestAccountInfo) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Information"
	return ths
}

func (ths *TestAccountInfo) execute() {
	_L.LoggerInstance.Info(" ** Test account information ** \r\n")
	countInDb := _OP.GetTestAccountCount(nil)
	tCountInDb := _OP.GetSuccessTestAccCount(nil)
	changeTrustCnt := _OP.GetNoAssetCodeTestAccCount(nil)
	assetBalanceValidCnt := int(_OP.GetAssetBalanceValidTestAccCount(nil, "0"))
	_L.LoggerInstance.InfoPrint(" >>             Total test account count : %d \r\n", countInDb)
	_L.LoggerInstance.InfoPrint(" >>            Active test account count : %d \r\n", tCountInDb)
	_L.LoggerInstance.InfoPrint(" >>          Inactive test account count : %d \r\n", countInDb-tCountInDb)
	_L.LoggerInstance.InfoPrint(" >>   Unchanged trust test account count : %d \r\n", changeTrustCnt)
	_L.LoggerInstance.InfoPrint(" >> Asset balance > 0 test account count : %d \r\n", assetBalanceValidCnt)
	ths.Wait()
}
