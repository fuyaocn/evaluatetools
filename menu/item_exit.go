package menu

import (
	_L "jojopoper/NBi/StressTest/log"
	"os"
)

// ExitApp 退出钱包程序
type ExitApp struct {
	SubItem
}

// InitMenu 初始化
func (ths *ExitApp) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Exit"
	return ths
}

func (ths *ExitApp) execute() {
	_L.LoggerInstance.Info(" ** Exit app ** \r\n")
	os.Exit(0)
}
