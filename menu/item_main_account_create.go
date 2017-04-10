package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"strconv"
	"sync"
)

// MainAccountCreator 创建主账户菜单
type MainAccountCreator struct {
	SubItem
	maxIndex     int
	startBalance string
	rootInfo     *api.AccountInfo
}

// InitMenu 初始化
func (ths *MainAccountCreator) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "New"
	ths.rootInfo = new(api.AccountInfo)
	ths.rootInfo.Init(_AC.ConfigInstance.GetString("rootaccountinfo", "id"), _AC.ConfigInstance.GetString("rootaccountinfo", "key"))
	return ths
}

func (ths *MainAccountCreator) execute() {
	_L.LoggerInstance.Info(" ** New main account ** \r\n")
	countInDb := _OP.GetMainAccountCount(nil)
	ths.maxIndex = _OP.GetMainAccountMaxIndex(nil)
	_L.LoggerInstance.InfoPrint(" >> Current main account count : %d \r\n", countInDb)
	_L.LoggerInstance.InfoPrint("  >> How many main account you want to creation(input 0 is cancel)? : ")
	numAccount, b := ths.InputNumber()
	if b && numAccount > 0 {
		_L.LoggerInstance.InfoPrint("  >> How many operations(pre transaction) you want to signed(1~100, input 0 is cancel)? : ")
		numOperation, b := ths.InputNumber()
		if b && numOperation > 0 && numOperation <= 100 {
			ths.run(numAccount, numOperation)
			return
		}
	}
	_L.LoggerInstance.InfoPrint("  >> Current operation has be canceled!\r\n")
}

func (ths *MainAccountCreator) run(accCnt, optCnt int) {
	wait := &sync.WaitGroup{}
	// wait.Init()
	wait.Add(2)
	go ths.getRootAccountInfo(wait)
	go ths.getStartBalance(wait)
	wait.Wait()
	if len(ths.startBalance) == 0 {
		return
	}

	group := accCnt / optCnt
	left := accCnt % optCnt
	current := 0
	gindex := ths.maxIndex + 1
	// wait.Init()

	for {
		if left > 0 {
			current = left
			left = 0
		} else if group > 0 {
			current = optCnt
			group--
		} else {
			break
		}
		wait.Add(1)
		go ths.createGroup(current, gindex, wait)
		wait.Wait()

		gindex += current
	}
}

func (ths *MainAccountCreator) getStartBalance(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	ths.startBalance = _AC.ConfigInstance.GetString("mainaccount", "startbalance")
	if len(ths.startBalance) == 0 {
		_L.LoggerInstance.InfoPrint("  >> How many balance you want active account(30 ~ max)? ")
		start, b := ths.InputFloat()
		if b && start >= 30.0 {
			ths.startBalance = fmt.Sprintf("%f", start)
		}
	}
}

func (ths *MainAccountCreator) getRootAccountInfo(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	err := ths.rootInfo.GetInfo(nil)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  *** Get root account information has error :\r\n %+v\r\n", err)
	}
}

func (ths *MainAccountCreator) createGroup(cnt int, gindex int, wg *sync.WaitGroup) {
	defer wg.Done()

	accs := _OP.GetNewMainAccs(cnt, gindex)
	if accs == nil {
		return
	}

	createAcc := &api.CreateAccount{}
	accIDs := ths.getAccountIDs(accs)
	signature := createAcc.GetSignature(ths.rootInfo, ths.startBalance, "create_test_acc", accIDs...)
	addr := _AC.ConfigInstance.GetHorizonServer() + "/transactions"
	err := createAcc.Send(addr, api.Linear, signature)
	if err == nil {
		bala, _ := strconv.ParseFloat(ths.startBalance, 64)
		err = _OP.SaveMainAccToDatabase(accs, bala, "T", nil)
	}

	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Create main account has error [grp:%d] : \r\n %+v\r\n", gindex, err)
	}
}

func (ths *MainAccountCreator) getAccountIDs(src []*_DB.TMainAccount) (ret []string) {
	if src == nil {
		return
	}
	ret = make([]string, 0)
	for _, itm := range src {
		ret = append(ret, itm.AccountID)
	}
	return
}
