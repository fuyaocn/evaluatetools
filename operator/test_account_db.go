package operator

import (
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"sync"

	"fmt"

	_kp "github.com/stellar/go/keypair"
)

// ClearTestAccount 清空测试账户数据库所有内容
func ClearTestAccount(w *sync.WaitGroup) (err error) {
	if w != nil {
		defer w.Done()
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount)
	_, err = operat.Query(_DB.QtClearAllRecord)
	return
}

// GetTestAccountCount 获取数据库中所有测试账户的个数
func GetTestAccountCount(w *sync.WaitGroup) int64 {
	if w != nil {
		defer w.Done()
	}
	tmp := &_DB.TTestAccount{}
	return getTestAccCount(tmp)
}

// GetSuccessTestAccCount 获取数据库中有效的测试账户的个数
func GetSuccessTestAccCount(w *sync.WaitGroup) int64 {
	if w != nil {
		defer w.Done()
	}
	tmp := &_DB.TTestAccount{
		Success: "T",
	}
	return getTestAccCount(tmp)
}

// GetNoAssetCodeTestAccCount 获取数据库中无AssetCode的测试账户的个数
func GetNoAssetCodeTestAccCount(w *sync.WaitGroup) int64 {
	if w != nil {
		defer w.Done()
	}
	tmp := &_DB.TTestAccount{
		AssetCode: "-",
		Success:   "T",
	}
	return getTestAccCount(tmp)
}

// GetAssetBalanceValidTestAccCount 获取数据库中AssetBalance有值的测试账户的个数
func GetAssetBalanceValidTestAccCount(w *sync.WaitGroup, validBalance string) int64 {
	if w != nil {
		defer w.Done()
	}
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	tmp := &_DB.TTestAccount{}
	ret, err := eng.Where(fmt.Sprintf("asset_balance >= %s", validBalance)).Count(tmp)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" >>>> Get test account count[asset_balance > %s] from database has error :\r\n %+v\r\n", validBalance, err)
		return -1
	}
	return ret
}

func getTestAccCount(src *_DB.TTestAccount) int64 {
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount)
	ret, err := operat.Query(_DB.QtGetCount, src)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" >>>> Get test account count from database has error :\r\n %+v\r\n", err)
		return -1
	}
	return ret.(int64)
}

var testAccLocker *sync.Mutex = new(sync.Mutex)

// GetNewTestAccs 得到一组测试账户信息
// cnt 需要创建的个数；index 第几组第几个的起始编号，grpIndex 所在的是第几组，grpItmIndex所在的是第几个
func GetNewTestAccs(cnt int, index, grpIndex, grpItmIndex int) (ret []*_DB.TTestAccount) {
	testAccLocker.Lock()
	defer testAccLocker.Unlock()
	if cnt <= 0 {
		return nil
	}
	ret = make([]*_DB.TTestAccount, cnt)
	for i := 0; i < cnt; i++ {
		ret[i] = CreateTestAcc(i+index, grpIndex, grpItmIndex)
		_L.LoggerInstance.DebugPrint(" > Addr = %s\r\n > Skey = %s\r\n", ret[i].AccountID, ret[i].SecertAddr)
	}
	return
}

// CreateTestAcc 创建一个测试账户信息
func CreateTestAcc(idx, grpIndex, grpItmIndex int) *_DB.TTestAccount {
	full, err := _kp.Random()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Create keypair has error \r\n%+v\r\n", err)
		return nil
	}
	return &_DB.TTestAccount{
		Index:          idx,
		GroupIndex:     grpIndex,
		GroupItemIndex: grpItmIndex,
		AccountID:      full.Address(),
		SecertAddr:     full.Seed(),
		Balance:        0,
		AssetBalance:   0,
		AssetCode:      "-",
	}
}

// SaveTestAccToDatabase 保存结果到数据库
func SaveTestAccToDatabase(src []*_DB.TTestAccount, balance float64, result string, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}
	for _, itm := range src {
		itm.Success = result
		itm.Balance = balance
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount)
	_, err := operat.Query(_DB.QtAddRecords, src)
	return err
}

// GetNumberofTestAccForChangTrust 从数据库中获取指定数量的未change-trust测试账户
func GetNumberofTestAccForChangTrust(num int) (ret []*_DB.TTestAccount) {
	if num <= 0 {
		return
	}
	ret = make([]*_DB.TTestAccount, 0)
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	err := eng.Where("success = 'T'").Where("asset_code='-'").Limit(num).Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get unchanged trust test account from database has error :\r\n %+v\r\n", err)
	}
	return
}

// GetNumberofTestAccForSendAsset 从数据库中获取指定数量的change-trust测试账户
func GetNumberofTestAccForSendAsset(num int, assetCode string) (ret []*_DB.TTestAccount) {
	if num <= 0 {
		return
	}
	ret = make([]*_DB.TTestAccount, 0)
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	err := eng.Where("success = 'T'").Where(fmt.Sprintf("asset_code='%s'", assetCode)).Limit(num).Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get changed trust test account from database has error :\r\n %+v\r\n", err)
	}
	return
}

// GetNumberofTestAccForAssetBalance 从数据库中获取指定数量的asset balance > 0测试账户
func GetNumberofTestAccForAssetBalance(num int, amount, assetCode string) (ret []*_DB.TTestAccount) {
	if num <= 0 {
		return
	}
	ret = make([]*_DB.TTestAccount, 0)
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	err := eng.Where(fmt.Sprintf("asset_code='%s'", assetCode)).Where(fmt.Sprintf("asset_balance > %s", amount)).Limit(num).Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get asset balance > 0 test account from database has error :\r\n %+v\r\n", err)
	}
	return
}

// GetNumberofTestAccForAssetBalanceFromTo 从数据库中获取指定数量的asset balance 在一个区间内的测试账户
func GetNumberofTestAccForAssetBalanceFromTo(amountFrom, amountTo, assetCode string) (ret []*_DB.TTestAccount) {
	ret = make([]*_DB.TTestAccount, 0)
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	err := eng.Where(fmt.Sprintf("asset_code='%s'", assetCode)).Where(fmt.Sprintf("asset_balance >= %s", amountTo)).Where(fmt.Sprintf("asset_balance < %s", amountFrom)).Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get asset balance > 0 test account from database has error :\r\n %+v\r\n", err)
	}
	return
}

// GetNumberofTestAccForIs0AssetBalance 从数据库中获取指定数量的asset balance == 0测试账户
func GetNumberofTestAccForIs0AssetBalance(assetCode string) (ret []*_DB.TTestAccount) {
	ret = make([]*_DB.TTestAccount, 0)
	eng := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount).GetEngine()
	err := eng.Where(fmt.Sprintf("asset_code='%s'", assetCode)).Where("asset_balance = 0").Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get asset balance == 0 test account from database has error :\r\n %+v\r\n", err)
	}
	return
}

// UpdateTestAccInfoToDB 更新数据库中测试账户信息
func UpdateTestAccInfoToDB(src []*_DB.TTestAccount) {
	opera := _DB.DataBaseInstance.GetOperation(_DB.KeyTestAccount)
	_, err := opera.Query(_DB.QtUpdateRecords, src)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Update test account informations has error :\r\n%+v\r\n", err)
	}
}
