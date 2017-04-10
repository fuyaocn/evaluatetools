package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"math"
	"strings"
	"sync"
)

// TestAccountChangeTrust 修改信任列表菜单
type TestAccountChangeTrust struct {
	SubItem
	wait               *sync.WaitGroup
	TestName           string
	totalUnchangeCnt   int64 // 所有未changtrust的账户个数
	willTochangeCnt    int   // 将要进行changetrust的账户个数
	isGroup            bool  // 是否要分组进行
	wantNumberList     []int // 分组时按照理论数据值分组， 2000,3000,5000,100000
	wantNumberListSize int
	groups             []*TestAccChangeTrustGroup
	groupsSize         int
	runType            api.SendType
	runTypeSub         api.SendType
	testAccs           []*api.AccountInfo
	accInDb            []*_DB.TTestAccount
	signerList         []string
	signerListSize     int
	static             []*_DB.TStatic
	issuer             string
	code               string
}

// InitMenu 初始化
func (ths *TestAccountChangeTrust) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Change trust"

	ths.issuer = _AC.ConfigInstance.GetString("issueraccountinfo", "id")
	ths.code = _AC.ConfigInstance.GetString("issueraccountinfo", "code")

	ths.wantNumberListSize = _AC.ConfigInstance.GetNumber("changetrustgroup", "group_level")
	ths.wantNumberList = make([]int, ths.wantNumberListSize)
	for i := 0; i < ths.wantNumberListSize; i++ {
		ths.wantNumberList[i] = _AC.ConfigInstance.GetNumber("changetrustgroup", fmt.Sprintf("group%d", (i+1)))
	}
	return ths
}

func (ths *TestAccountChangeTrust) execute() {
	_L.LoggerInstance.Info(" ** Test account change trust ** \r\n")
	step := 0
	b := false
	for {
		switch step {
		case 0:
			ths.totalUnchangeCnt = _OP.GetNoAssetCodeTestAccCount(nil)
			countInDb := _OP.GetTestAccountCount(nil)
			tCountInDb := _OP.GetSuccessTestAccCount(nil)
			_L.LoggerInstance.InfoPrint(" >> Current   total   test account count : %d \r\n", countInDb)
			_L.LoggerInstance.InfoPrint(" >> Current  active   test account count : %d \r\n", tCountInDb)
			_L.LoggerInstance.InfoPrint(" >> Current inactive  test account count : %d \r\n", countInDb-tCountInDb)
			_L.LoggerInstance.InfoPrint(" >> Current not asset test account count : %d \r\n", ths.totalUnchangeCnt)
			step++
		case 1:
			_L.LoggerInstance.InfoPrint("  >> Please enter a name for this test(lenght <= 255): ")
			ths.TestName, _ = ths.InputString()
			step++
		case 2:
			_L.LoggerInstance.InfoPrint(" >> Are you want to create group (y/n others:Cancel)? ")
			input, b := ths.InputString()
			if b && strings.ToLower(input) == "n" { // 不分组进行测试
				ths.isGroup = false
				step++
			} else if b && strings.ToLower(input) == "y" { // 分组进行测试
				ths.isGroup = true
				step += 3
			} else {
				step = -1
			}
		case 3:
			_L.LoggerInstance.InfoPrint(" >> How many test account to be changed trust (1~%d, 0 is cancel)? ", ths.totalUnchangeCnt)
			ths.willTochangeCnt, b = ths.InputNumber()
			if b && ths.willTochangeCnt > 0 && ths.willTochangeCnt <= int(ths.totalUnchangeCnt) {
				step++
			} else {
				step = -1
			}
		case 4:
			_L.LoggerInstance.InfoPrint(" >> What type you want to run (1.Linear 2.Multiple others:Cancel)? : ")
			selectNumber, b := ths.InputNumber()
			if b {
				step += 2
				switch selectNumber {
				case 1:
					ths.runType = api.Linear
				case 2:
					ths.runType = api.Multiple
				default:
					step = -1
				}

			} else {
				step = -1
			}
		case 5:
			output := " >> The program will split the test accounts into groups, with each group changing-trust the test accounts in the following order: \r\n"
			for i, itm := range ths.wantNumberList {
				output += fmt.Sprintf("  >> Group%d : Account = %d\r\n", (i + 1), itm)
			}
			output += fmt.Sprintf("  >> Group%d : Account = %d\r\n  >> and then cycle ...\r\n", (ths.wantNumberListSize + 1),
				ths.wantNumberList[0])
			_L.LoggerInstance.InfoPrint(output)
			// for debug
			// _L.LoggerInstance.DebugPrint(" group = \r\n%+v\r\n", ths.groups)
			// for i, itm := range ths.groups {
			// 	_L.LoggerInstance.DebugPrint(" group [%d] = \r\n%+v\r\n", i+1, itm)
			// }
			_L.LoggerInstance.InfoPrint(" >> What type you want to run (1.Linear-Linear 2.Linear-Multiple 3.Multiple-Linear 4.Multiple-Multiple others:Cancel)? : ")
			selectNumber, b := ths.InputNumber()
			if b {
				step++
				switch selectNumber {
				case 1:
					ths.runType = api.Linear
					ths.runTypeSub = api.Linear
				case 2:
					ths.runType = api.Linear
					ths.runTypeSub = api.Multiple
				case 3:
					ths.runType = api.Multiple
					ths.runTypeSub = api.Linear
				case 4:
					ths.runType = api.Multiple
					ths.runTypeSub = api.Multiple
				default:
					step = -1
				}
				ths.getGroup()
				ths.groupsSize = len(ths.groups)
				_L.LoggerInstance.InfoPrint(" >> Program create %d groups.\r\n", ths.groupsSize)

				// _L.LoggerInstance.DebugPrint(" >>>>>> runType = %d\trunTypeSub = %d\r\n", ths.runType, ths.runTypeSub)

			} else {
				step = -1
			}
		case 6:
			info := " >> Please confirm fellow informations:\r\n"
			info += fmt.Sprintf(" >> Test name: %s\r\n", ths.TestName)
			if ths.isGroup {
				info += fmt.Sprintf(" >> Group test: %d groups\r\n", ths.groupsSize)
			} else {
				info += fmt.Sprintf(" >> Change trust test account: %d\r\n", ths.willTochangeCnt)
			}
			switch ths.runType {
			case api.Linear:
				info += " >> runType: Linear\r\n"
			case api.Multiple:
				info += " >> runType: Multiple\r\n"
			}
			if ths.isGroup {
				switch ths.runTypeSub {
				case api.Linear:
					info += " >> group runType: Linear\r\n"
				case api.Multiple:
					info += " >> group runType: Multiple\r\n"
				}
			}
			info += " >> Are you confirm(y/n)? "
			_L.LoggerInstance.InfoPrint(info)
			input, b := ths.InputString()
			if b && strings.ToLower(input) == "y" {
				step++
			} else {
				step = -1
			}
		case 7:
			ths.run()
			return
		default:
			ths.OperationCanceled()
			return
		}
	}
}

func (ths *TestAccountChangeTrust) run() {
	num := 0
	if ths.isGroup {
		num = int(ths.totalUnchangeCnt)
	} else {
		num = ths.willTochangeCnt
	}

	if !ths.getTestAccount(num) || !ths.checkTestAccInfo() {
		ths.OperationCanceled()
		return
	}

	if ths.isGroup {
		ths.groupRun()
	} else {
		ths.liteRun()
	}
}

func (ths *TestAccountChangeTrust) getGroup() {
	ths.getTestAccount(int(ths.totalUnchangeCnt))
	ths.groups = make([]*TestAccChangeTrustGroup, 0)
	grpIdx := 0
	total := int(ths.totalUnchangeCnt)
	startIndex := 0
	endIndex := 0
	left := total
	for left > 0 {
		wantIndex := grpIdx % ths.wantNumberListSize
		if left < ths.wantNumberList[wantIndex] {
			endIndex += left
		} else {
			endIndex += ths.wantNumberList[wantIndex]
		}
		size := endIndex - startIndex
		grp := &TestAccChangeTrustGroup{
			ItemGroupBase: ItemGroupBase{
				TestName:             ths.TestName,
				GroupID:              grpIdx,
				RunType:              ths.runTypeSub,
				LengthOfTestAccounts: size,
				Issuer:               ths.issuer,
				Code:                 ths.code,
			},
		}
		// _L.LoggerInstance.DebugPrint(" >>>>>> group info = %+v\r\n", grp)
		grp.TestAccounts = make([]*api.AccountInfo, size)
		grp.AccInDb = make([]*_DB.TTestAccount, size)
		copy(grp.TestAccounts, ths.testAccs[startIndex:endIndex])
		copy(grp.AccInDb, ths.accInDb[startIndex:endIndex])
		startIndex = endIndex
		left = total - endIndex
		grpIdx++
		ths.groups = append(ths.groups, grp)
	}
}

func (ths *TestAccountChangeTrust) getTestAccount(num int) bool {
	_L.LoggerInstance.InfoPrint(" >> Get test account from database ...\r\n")
	ths.accInDb = _OP.GetNumberofTestAccForChangTrust(num)
	if ths.accInDb == nil {
		return false
	}
	cntAcc := len(ths.accInDb)
	_L.LoggerInstance.InfoPrint(" >> Read %d recoeds from database\r\n", cntAcc)
	ths.testAccs = make([]*api.AccountInfo, cntAcc)
	for i, itm := range ths.accInDb {
		ths.testAccs[i] = new(api.AccountInfo)
		ths.testAccs[i].Init(itm.AccountID, itm.SecertAddr)
	}
	return true
}

func (ths *TestAccountChangeTrust) checkTestAccInfo() bool {
	_L.LoggerInstance.InfoPrint(" >> Get account information from horizon ...\r\n")
	cntAcc := len(ths.testAccs)
	_L.LoggerInstance.DebugPrint(" ths.testAccs length = %d\r\n", cntAcc)
	groupCnt := 500
	if groupCnt > cntAcc {
		groupCnt = cntAcc
	}
	ths.wait = new(sync.WaitGroup)
	cycCnt := int(math.Ceil(float64(cntAcc) / float64(groupCnt)))
	_L.LoggerInstance.DebugPrint(" cycCnt = %d\r\n", cycCnt)
	leftCnt := cntAcc % groupCnt
	_L.LoggerInstance.DebugPrint(" leftCnt = %d\r\n", leftCnt)
	for idx := 0; idx < cycCnt; idx++ {
		index := idx * groupCnt
		_L.LoggerInstance.DebugPrint(" index = %d\r\n", index)
		if idx == cycCnt-1 && leftCnt > 0 {
			groupCnt = leftCnt
		}
		_L.LoggerInstance.DebugPrint(" groupCnt = %d\r\n", groupCnt)

		for i := 0; i < groupCnt; i++ {
			ths.wait.Add(1)
			go func(acc *api.AccountInfo) {
				err := acc.GetInfo(ths.wait)
				if err != nil {
					_L.LoggerInstance.ErrorPrint("  **** get account[%s] has error:\r\n%+v\r\n", acc.Address)
				}
			}(ths.testAccs[i+index])
		}
		ths.wait.Wait()
	}
	ths.wait.Wait()

	for _, itm := range ths.testAccs {
		if itm.Status != 0 {
			return false
		}
	}
	_L.LoggerInstance.InfoPrint(" >> Get account information from horizon complete!\r\n")
	return true
}

func (ths *TestAccountChangeTrust) groupRun() {
	ths.wait = new(sync.WaitGroup)
	for i := 0; i < ths.groupsSize; i++ {
		grp := ths.groups[i]
		if ths.runType == api.Linear {
			// _L.LoggerInstance.Debug(" group id = %d execute!\r\n", i)
			grp.Do(nil)
		} else {
			ths.wait.Add(1)
			go grp.Do(ths.wait)
		}
	}
	ths.wait.Wait()

	_L.LoggerInstance.InfoPrint("  >> Test account change trust complete!\r\n")
}

func (ths *TestAccountChangeTrust) liteRun() {
	_L.LoggerInstance.InfoPrint(" >> Issuer = %s\r\n >> Code = %s\r\n", ths.issuer, ths.code)
	_L.LoggerInstance.InfoPrint(" >> Get signatures and post to horizon...\r\n")

	grp := &TestAccChangeTrustGroup{
		ItemGroupBase: ItemGroupBase{
			TestName:             ths.TestName,
			GroupID:              -1,
			RunType:              ths.runType,
			LengthOfTestAccounts: ths.willTochangeCnt,
			Issuer:               ths.issuer,
			Code:                 ths.code,
		},
	}
	// _L.LoggerInstance.DebugPrint(" >>>>>> group info = %+v\r\n", grp)
	grp.TestAccounts = make([]*api.AccountInfo, ths.willTochangeCnt)
	grp.AccInDb = make([]*_DB.TTestAccount, ths.willTochangeCnt)
	copy(grp.TestAccounts, ths.testAccs)
	copy(grp.AccInDb, ths.accInDb)
	grp.Do(nil)
}
