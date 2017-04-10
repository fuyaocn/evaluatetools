package menu

import (
	"fmt"
	_AC "jojopoper/NBi/StressTest/appconf"
)

// MainMenu 菜单
type MainMenu struct {
	SubItem
}

// MainMenuInstace 菜单唯一实例
var MainMenuInstace SubItemInterface

// CreateMenuInstance 创建主菜单
func CreateMenuInstance() SubItemInterface {
	ret := &MainMenu{}
	return ret.Init()
}

// Init 初始化
func (ths *MainMenu) Init() SubItemInterface {
	ths.InitMenu(nil, "0")
	ths.title = "Main menu"
	ths.Exec = ths.execute

	mainAcc := &MainAccount{}
	ths.AddSubItem(mainAcc.InitMenu(ths, "0"))

	testAcc := &TestAccount{}
	ths.AddSubItem(testAcc.InitMenu(ths, "0"))

	payAsset := &PaymentAsset{}
	ths.AddSubItem(payAsset.InitMenu(ths, "0"))

	makeorder := &MakeOrderAB{}
	ths.AddSubItem(makeorder.InitMenu(ths, "0"))

	about := &SoftwareAbout{}
	ths.AddSubItem(about.InitMenu(ths, "0"))

	exit := &ExitApp{}
	ths.AddSubItem(exit.InitMenu(ths, "0"))
	return ths
}

// Execute 执行函数
func (ths *MainMenu) execute() {
	fmt.Println("\r\n*****************************************************")
	fmt.Printf("##               %s               ##\r\n", _AC.ConfigInstance.GetAppName())
	ths.SubItem.execute()
}
