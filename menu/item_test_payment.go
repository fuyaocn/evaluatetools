package menu

// PaymentAsset 本地测试账户菜单
type PaymentAsset struct {
	SubItem
}

// InitMenu 初始化菜单
func (ths *PaymentAsset) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Payment Asset"

	testAccInfo := &TestAccountInfo{}
	ths.AddSubItem(testAccInfo.InitMenu(ths, "1"))

	rootPay := &RootPaymentAsset{}
	ths.AddSubItem(rootPay.InitMenu(ths, "1"))

	msPay := &PayAssetGroupAB{}
	ths.AddSubItem(msPay.InitMenu(ths, "1"))

	rp := &ReturnParentMenu{}
	ths.AddSubItem(rp.InitMenu(ths, "1"))

	ea := &ExitApp{}
	ths.AddSubItem(ea.InitMenu(ths, "1"))

	return ths
}
