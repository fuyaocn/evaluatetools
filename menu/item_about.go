package menu

import (
	"fmt"
	_L "jojopoper/NBi/StressTest/log"
)

// SoftwareAbout 关于
type SoftwareAbout struct {
	SubItem
}

// InitMenu 初始化
func (ths *SoftwareAbout) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "About"
	return ths
}

func (ths *SoftwareAbout) execute() {
	_L.LoggerInstance.Info(" ** About show ** \r\n")
	fmt.Println("")
	fmt.Printf("\tSoftware Version : 1.1.0.0\r\n")
	fmt.Println("")
	ths.Wait()
}
