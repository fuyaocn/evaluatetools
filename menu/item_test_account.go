package menu

// TestAccount 本地测试账户菜单
type TestAccount struct {
	SubItem
}

// InitMenu 初始化菜单
func (ths *TestAccount) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Test Account"

	testAccInfo := &TestAccountInfo{}
	ths.AddSubItem(testAccInfo.InitMenu(ths, "1"))

	testAccCrt := &TestAccountCreate{}
	ths.AddSubItem(testAccCrt.InitMenu(ths, "1"))

	changtr := &TestAccountChangeTrust{}
	ths.AddSubItem(changtr.InitMenu(ths, "1"))

	clrTestAcc := &TestAccountClear{}
	ths.AddSubItem(clrTestAcc.InitMenu(ths, "1"))

	rp := &ReturnParentMenu{}
	ths.AddSubItem(rp.InitMenu(ths, "1"))

	ea := &ExitApp{}
	ths.AddSubItem(ea.InitMenu(ths, "1"))

	return ths
}
