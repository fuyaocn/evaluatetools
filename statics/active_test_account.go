package statics

// ActiveTestAccount 激活测试账户统计信息定义
type ActiveTestAccount struct {
	BaseDefine
	Signature      string
	MainAccID      string
	OperationCnt   int
	TransactionCnt int
}

// NewActivcTestAccStatic 创建一个激活测试账户的统计实体
func NewActivcTestAccStatic() (ret *ActiveTestAccount) {
	ret = new(ActiveTestAccount)
	return ret.Init()
}

// Init 初始化
func (ths *ActiveTestAccount) Init() *ActiveTestAccount {
	ths.Action = "active_test_acc"
	return ths
}

// SetSign 设置签名字符串
func (ths *ActiveTestAccount) SetSign(s string) {
	ths.Signature = s
}

// SetCount 设置个数
// oCnt operation的个数; tCnt transaction的个数
func (ths *ActiveTestAccount) SetCount(oCnt, tCnt int) {
	ths.OperationCnt = oCnt
	ths.TransactionCnt = tCnt
}

// SetMainAccAddr 设置主账户的ID
func (ths *ActiveTestAccount) SetMainAccAddr(aid string) {
	ths.MainAccID = aid
}
