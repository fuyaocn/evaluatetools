package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"math"
	"sync"
)

// TestAccountCreate 测试账户创建菜单
type TestAccountCreate struct {
	SubItem
	// numOperation    int
	runTypeOut      api.SendType // 外部大循环线程运行方式
	runTypeIn       api.SendType // 内部小循环线程运行方式
	TestName        string       // 测试名称
	startBalance    string
	mainAccounts    []*MainAccountGroup
	numTestCnt      int   // 每个主账户需要负责生成多少个子账户
	mainAccCount    int64 // 数据库中主账户的总数
	mainAccCntGroup int   // 需要多少个主账户为一组
	wait            *sync.WaitGroup
}

// InitMenu 初始化
func (ths *TestAccountCreate) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Create"
	ths.wait = new(sync.WaitGroup)
	return ths
}

func (ths *TestAccountCreate) execute() {
	_L.LoggerInstance.Info(" ** Test account creatation ** \r\n")
	ths.mainAccCount = _OP.GetMainAccountCount(nil)
	_L.LoggerInstance.InfoPrint(" >> Current main account count : %d \r\n", ths.mainAccCount)
	_L.LoggerInstance.InfoPrint(" >> The program will split the master accounts into groups, with each group generating the test accounts in the following order: \r\n  >> Each main account generates N test accounts.\r\n  >> Group1 : opration count = 100/transaction\r\n  >> Group2 : opration count = 50/transaction\r\n  >> Group3 : opration count = 20/transaction\r\n  >> Group4 : opration count = 10/transaction\r\n  >> Group5 : opration count = 100/transaction\r\n  >> and then cycle ...\r\n")

	step := 0
	b := false
	for {
		switch step {
		case 0:
			if ths.mainAccCount < 1 {
				_L.LoggerInstance.InfoPrint("  >> There are not enough main accounts in the database!\r\n")
				step = -1
			} else {
				step++
			}
		case 1:
			_L.LoggerInstance.InfoPrint("  >> Please enter a name for this test(lenght <= 255): ")
			ths.TestName, _ = ths.InputString()
			step++
		case 2:
			_L.LoggerInstance.InfoPrint("  >> How many test accounts will be generated per main account? (1~100, input 0 is cancel) ? : ")
			ths.numTestCnt, b = ths.InputNumber()
			if b && ths.numTestCnt > 0 && ths.numTestCnt <= 100 {
				step++
			} else {
				step = -1
			}
		case 3:
			_L.LoggerInstance.InfoPrint("  >> How many main accounts to be a group? (1~%d, input 0 is cancel) ? : ", ths.mainAccCount)
			ths.mainAccCntGroup, b = ths.InputNumber()
			if b && ths.mainAccCntGroup > 0 && int64(ths.mainAccCntGroup) <= ths.mainAccCount {
				step++
			} else {
				step = -1
			}
		case 4:
			ths.getStartBalance(nil)
			if len(ths.startBalance) > 0 {
				step++
			} else {
				step = -1
			}
		case 5:
			_L.LoggerInstance.InfoPrint("  >> What type you want to run (1.Linear-Linear 2.Linear-Multiple 3.Multiple-Linear 4.Multiple-Multiple others:Cancel)? : ")
			selectNumber, b := ths.InputNumber()
			if b {
				step++
				switch selectNumber {
				case 1:
					ths.runTypeOut = api.Linear
					ths.runTypeIn = api.Linear
				case 2:
					ths.runTypeOut = api.Linear
					ths.runTypeIn = api.Multiple
				case 3:
					ths.runTypeOut = api.Multiple
					ths.runTypeIn = api.Linear
				case 4:
					ths.runTypeOut = api.Multiple
					ths.runTypeIn = api.Multiple
				default:
					step = -1
				}

			} else {
				step = -1
			}
		case 6:
			_L.LoggerInstance.InfoPrint("  >> Checking main account validation ... \r\n")
			// ths.wait.Init()
			ths.wait.Add(1)
			go ths.getMainAccountInfo(ths.wait)
			ths.wait.Wait()
			_L.LoggerInstance.InfoPrint("  >> Check main account complete. \r\n")
			if ths.mainAccounts != nil {
				step++
			} else {
				step = -1
			}
		case 7:
			ths.run()
			return
		default:
			_L.LoggerInstance.InfoPrint("  >> Current operation has be canceled!\r\n")
			return
		}
	}

}

func (ths *TestAccountCreate) getStartBalance(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	ths.startBalance = _AC.ConfigInstance.GetString("testaccount", "startbalance")
	if len(ths.startBalance) == 0 {
		_L.LoggerInstance.InfoPrint("  >> How many balance you want active test account(50 ~ max)? ")
		start, b := ths.InputFloat()
		if b && start >= 50.0 {
			ths.startBalance = fmt.Sprintf("%f", start)
		}
	}
}

func (ths *TestAccountCreate) getMainAccountInfo(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	ths.mainAccounts = nil
	// 取出所有主账户
	mAccounts := _OP.GetMainAccFromDB(ths.mainAccCount, nil)
	// 算出需要分多少组
	mainGroupCnt := int(math.Ceil(float64(ths.mainAccCount) / float64(ths.mainAccCntGroup)))
	_L.LoggerInstance.InfoPrint("  >> The main account will be divided into %d groups.\r\n", mainGroupCnt)
	wait := &sync.WaitGroup{}
	ths.mainAccounts = make([]*MainAccountGroup, 0)
	mainAccCntLeft := int(ths.mainAccCount)
	for i := 0; i < mainGroupCnt; i++ {
		group := &MainAccountGroup{
			TestName:       ths.TestName,
			GroupID:        i,
			runType:        ths.runTypeIn,
			TestAccountCnt: ths.numTestCnt,
			StartBalance:   ths.startBalance,
		}
		group.MainAccounts = make([]*api.AccountInfo, 0)
		index := i * ths.mainAccCntGroup
		if mainAccCntLeft < ths.mainAccCntGroup {
			ths.mainAccCntGroup = mainAccCntLeft
		}
		// wait.Init()
		for j := 0; j < ths.mainAccCntGroup; j++ {
			item := new(api.AccountInfo)
			item.Init(mAccounts[index+j].AccountID, mAccounts[index+j].SecertAddr)
			mainAccCntLeft--
			wait.Add(1)
			go item.GetInfo(wait)
			group.MainAccounts = append(group.MainAccounts, item)
		}
		switch i % 4 {
		case 0:
			group.OperationCnt = 100
		case 1:
			group.OperationCnt = 50
		case 2:
			group.OperationCnt = 20
		case 3:
			group.OperationCnt = 10
		}
		wait.Wait()
		ths.mainAccounts = append(ths.mainAccounts, group)
	}
}

func (ths *TestAccountCreate) run() {
	// wait := &sync.WaitGroup{}
	length := len(ths.mainAccounts)
	// ths.wait.Init()
	for i := 0; i < length; i++ {
		grp := ths.mainAccounts[i]
		if ths.runTypeOut == api.Linear {
			// _L.LoggerInstance.Debug(" group id = %d execute!\r\n", i)
			grp.Do(nil)
		} else {
			ths.wait.Add(1)
			go grp.Do(ths.wait)
		}
	}
	ths.wait.Wait()

	_L.LoggerInstance.InfoPrint("  >> Create test account complete!\r\n")
}

// printMainAccounts 测试代码使用的
func (ths *TestAccountCreate) printMainAccounts() {
	fmt.Printf("\r\n ===========================================\r\n")
	fmt.Printf(" mainAccCounts Info :\r\n %+v\r\n", ths.mainAccounts)
	fmt.Printf("\r\n *******************************************\r\n")
	for i, itm := range ths.mainAccounts {
		fmt.Printf(" [%d]\t Group Info :\r\n %+v\r\n", i, itm)
		for j, tm := range itm.MainAccounts {
			fmt.Printf("\t\t [%d-%d] Main account Info : %+v\r\n", i, j, tm)
		}
	}
}
