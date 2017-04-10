package operator

import (
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"strconv"
	"sync"

	_kp "github.com/stellar/go/keypair"
)

// ClearMainAccount 清空主账户数据库所有内容
func ClearMainAccount(w *sync.WaitGroup) (err error) {
	if w != nil {
		defer w.Done()
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyMainAccount)
	_, err = operat.Query(_DB.QtClearAllRecord)
	return
}

// GetMainAccountCount 获取数据库中主账户所有个数
func GetMainAccountCount(w *sync.WaitGroup) int64 {
	if w != nil {
		defer w.Done()
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyMainAccount)
	tmp := &_DB.TMainAccount{}
	ret, err := operat.Query(_DB.QtGetCount, tmp)
	if err != nil {
		_L.LoggerInstance.ErrorPrint(" >>>> Get main account count from database has error :\r\n %+v\r\n", err)
		return -1
	}
	return ret.(int64)
}

// GetMainAccountMaxIndex 获取数据库中主账户Index最大值
func GetMainAccountMaxIndex(w *sync.WaitGroup) int {
	if w != nil {
		defer w.Done()
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyMainAccount)
	eng := operat.GetEngine()
	cmd := "select max(index) from t_main_account"
	ret, err := eng.Query(cmd)
	if err == nil {
		val, ok := ret[0]["max"]
		if ok {
			ret, err := strconv.Atoi(string(val))
			if err == nil {
				return ret
			}
		} else {
			return 0
		}
	}
	_L.LoggerInstance.ErrorPrint(" >>>> Get main account Max(Index) from database has error :\r\n %+v\r\n", err)
	return 0
}

// GetMainAccFromDB 从数据库中获取指定数量的主账户
func GetMainAccFromDB(num int64, w *sync.WaitGroup) (ret []*_DB.TMainAccount) {
	if w != nil {
		defer w.Done()
	}
	ret = make([]*_DB.TMainAccount, 0)
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyMainAccount)
	eng := operat.GetEngine()
	err := eng.Where("`index`>=? and `index`<=? and `success`=?", 1, num, "T").Find(&ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Get main account [index = 1~%d] has error :\r\n %+v\r\n", num, err)
	}
	return
}

// GetNewMainAccs 得到一组主账户信息
// cnt 需要创建的个数；grpIndex 第几组
func GetNewMainAccs(cnt int, grpIndex int) (ret []*_DB.TMainAccount) {
	if cnt <= 0 {
		return nil
	}
	ret = make([]*_DB.TMainAccount, cnt)
	for i := 0; i < cnt; i++ {
		ret[i] = CreateMainAcc(i + grpIndex)
		_L.LoggerInstance.DebugPrint(" > Addr = %s\r\n > Skey = %s\r\n", ret[i].AccountID, ret[i].SecertAddr)
	}
	return
}

// CreateMainAcc 创建一个主账户信息
func CreateMainAcc(idx int) *_DB.TMainAccount {
	full, err := _kp.Random()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Create keypair has error \r\n%+v\r\n", err)
		return nil
	}
	return &_DB.TMainAccount{
		Index:      idx,
		AccountID:  full.Address(),
		SecertAddr: full.Seed(),
		Balance:    0,
	}
}

// SaveMainAccToDatabase 保存结果到数据库
func SaveMainAccToDatabase(src []*_DB.TMainAccount, balance float64, result string, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}
	for _, itm := range src {
		itm.Success = result
		itm.Balance = balance
	}
	operat := _DB.DataBaseInstance.GetOperation(_DB.KeyMainAccount)
	_, err := operat.Query(_DB.QtAddRecords, src)
	return err
}
