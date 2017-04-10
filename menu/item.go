package menu

import (
	"fmt"
	_L "jojopoper/NBi/StressTest/log"
)

const (
	NormalFlag     = 1
	BackToMenuFlag = 99999
)

// SubItem 子菜单定义
type SubItem struct {
	title      string
	Exec       func()
	subItems   []SubItemInterface
	parentItem SubItemInterface
	keyPath    string
	ExeFlag    int
	inputMemo  string
}

// SubItemInterface 子菜单接口定义
type SubItemInterface interface {
	InitMenu(parent SubItemInterface, key string) SubItemInterface
	GetTitle() string
	HasTitle() bool
	SetTitle(t string)
	GetSubItems() []SubItemInterface
	AddSubItem(itm SubItemInterface) int
	GetParentItem() SubItemInterface
	SetParentItem(p SubItemInterface)
	GetKeyPath() string
	SetKeyPath(kp string)
	GetTitlePath() string
	ExecuteFunc()
	ExecFlag() int
	PrintSubmenu()
	GetInputMemo() string
}

// InitMenu 初始化
func (ths *SubItem) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.ExeFlag = NormalFlag
	ths.parentItem = parent
	ths.Exec = ths.execute
	ths.subItems = make([]SubItemInterface, 0)
	ths.keyPath = key
	ths.inputMemo = "Select the number of items in the menu list item and press enter: "
	return ths
}

// GetTitle 获取标题
func (ths *SubItem) GetTitle() string {
	return ths.title
}

// HasTitle 获取是否含有标题
func (ths *SubItem) HasTitle() bool {
	return len(ths.title) > 0
}

// SetTitle 设置标题
func (ths *SubItem) SetTitle(t string) {
	ths.title = t
}

// GetSubItems 获取子菜单
func (ths *SubItem) GetSubItems() []SubItemInterface {
	return ths.subItems
}

// AddSubItem 添加子菜单
func (ths *SubItem) AddSubItem(itm SubItemInterface) int {
	length := len(ths.subItems)
	itm.SetParentItem(ths)
	itm.SetKeyPath(fmt.Sprintf("%s.%d", ths.keyPath, length))
	ths.subItems = append(ths.subItems, itm)
	return length
}

// GetParentItem 得到父菜单
func (ths *SubItem) GetParentItem() SubItemInterface {
	return ths.parentItem
}

// SetParentItem 设置父菜单
func (ths *SubItem) SetParentItem(p SubItemInterface) {
	ths.parentItem = p
}

// GetKeyPath 获取路径
func (ths *SubItem) GetKeyPath() string {
	return ths.keyPath
}

// SetKeyPath 设置路径
func (ths *SubItem) SetKeyPath(kp string) {
	ths.keyPath = kp
}

// GetTitlePath 获取标题路径
func (ths *SubItem) GetTitlePath() (ret string) {
	if ths.parentItem == nil {
		ret = ths.title
	} else {
		ret = ths.parentItem.GetTitlePath() + " > " + ths.title
	}
	return
}

// ExecuteFunc 执行
func (ths *SubItem) ExecuteFunc() {
	if ths.Exec != nil {
		ths.Exec()
	}
}

// ExecFlag 执行标识
func (ths *SubItem) ExecFlag() int {
	return ths.ExeFlag
}

// PrintSubmenu 打印子菜单
func (ths *SubItem) PrintSubmenu() {
	length := len(ths.subItems)
	for i := 0; i < length; i++ {
		fmt.Printf(" > %d.\t%s\r\n", i+1, ths.subItems[i].GetTitle())
	}
}

// GetInputMemo 得到说明
func (ths *SubItem) GetInputMemo() string {
	return ths.inputMemo
}

func (ths *SubItem) execute() {
	_L.LoggerInstance.Info(" *** %s ***\r\n", ths.GetTitle())
	for {
		fmt.Printf("\n\n %s\r\n\n", ths.GetTitlePath())
		ths.PrintSubmenu()
		fmt.Printf("\n  %s", ths.GetInputMemo())

		selectIndex, b := ths.InputNumber()
		if b {
			if selectIndex <= len(ths.subItems) && selectIndex >= 0 {
				ths.subItems[selectIndex-1].ExecuteFunc()
				ret := ths.subItems[selectIndex-1].ExecFlag()
				if ret == BackToMenuFlag {
					break
				}
			}
		}
	}
}

// InputString 输入字符串
func (ths *SubItem) InputString() (input string, b bool) {
	n, err := fmt.Scanf("%s\n", &input)
	if err == nil && n > 0 {
		b = true
	}
	ths.LogUserInput("string", input)
	return
}

// InputNumber 输入数字
func (ths *SubItem) InputNumber() (input int, b bool) {
	n, err := fmt.Scanf("%d\n", &input)
	if err == nil && n > 0 {
		b = true
	}
	ths.LogUserInput("number", input)
	return
}

// InputFloat 输入浮点数字
func (ths *SubItem) InputFloat() (input float64, b bool) {
	n, err := fmt.Scanf("%f\n", &input)
	if err == nil && n > 0 {
		b = true
	}
	ths.LogUserInput("float", input)
	return
}

// Wait 等待
func (ths *SubItem) Wait() {
	fmt.Printf("\r\n > Press Enter to continue...")
	input := ""
	fmt.Scanf("%s\n", &input)
}

// LogUserInput 记录用户输入的结果
func (ths *SubItem) LogUserInput(t string, v interface{}) {
	_L.LoggerInstance.Info(" $$ User input [%s] is %v\r\n", t, v)
}

// OperationCanceled 操作取消
func (ths *SubItem) OperationCanceled() {
	_L.LoggerInstance.InfoPrint(" >> Current operation has be canceled!\r\n")
}
