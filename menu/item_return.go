package menu

import (
	_L "jojopoper/NBi/StressTest/log"
)

// ReturnParentMenu 返回上一级
type ReturnParentMenu struct {
	SubItem
}

// InitMenu 初始化
func (ths *ReturnParentMenu) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Go back"
    ths.ExeFlag = BackToMenuFlag
	return ths
}

func (ths *ReturnParentMenu) execute() {
	_L.LoggerInstance.Info(" ** go back ** \r\n")
}
