package menu

// MainAccount 本地主账户菜单
type MainAccount struct {
	SubItem
}

// InitMenu 初始化菜单
func (ths *MainAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Main Account"

	mainAccInfo := &MainAccountInfo{}
	ths.AddSubItem(mainAccInfo.InitMenu(ths, "1"))

	mainAccCrt := &MainAccountCreator{}
	ths.AddSubItem(mainAccCrt.InitMenu(ths, "1"))

	mainAccClr := &MainAccountClear{}
	ths.AddSubItem(mainAccClr.InitMenu(ths, "1"))

	rp := &ReturnParentMenu{}
	ths.AddSubItem(rp.InitMenu(ths, "1"))

	ea := &ExitApp{}
	ths.AddSubItem(ea.InitMenu(ths, "1"))

	return ths
}
