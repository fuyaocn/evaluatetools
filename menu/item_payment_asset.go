package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"strings"
	"sync"
)

// RootPayAssetGroup root发送资产到测试账户预定义分组格式
type RootPayAssetGroup struct {
	AccountCnt   int
	OperationCnt int
}

// RootPaymentAsset 根账户发资产给测试账户菜单
type RootPaymentAsset struct {
	SubItem
	TestName           string
	rootAccInfo        *api.AccountInfo // 发送资产的账户
	sendTestAccCnt     int              // 非分组状态实际接收资产测试账户个数
	operationCnt       int              // 非分组状态一个transaction中承载的operation个数
	changeTrustCnt     int              // 数据库中已经信任资产的账户
	sendAmount         string           // 发送的金额
	sendCode           string           // 发送的code
	sendIssuer         string           // 发送的网关地址
	isGroup            bool
	groups             []*RootPaymentAssetGroup
	groupSize          int
	wantNumberList     []*RootPayAssetGroup
	wantNumberListSize int
	runType            api.SendType
	runTypeSub         api.SendType
	testAccs           []*api.AccountInfo
	accInDb            []*_DB.TTestAccount
	wait               *sync.WaitGroup
}

// InitMenu 初始化
func (ths *RootPaymentAsset) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Root payment"
	ths.rootAccInfo = new(api.AccountInfo)
	ths.rootAccInfo.Init(_AC.ConfigInstance.GetString("rootaccountinfo", "id"), _AC.ConfigInstance.GetString("rootaccountinfo", "key"))
	ths.sendCode = _AC.ConfigInstance.GetString("issueraccountinfo", "code")
	ths.sendIssuer = _AC.ConfigInstance.GetString("issueraccountinfo", "id")
	ths.wantNumberListSize = _AC.ConfigInstance.GetNumber("rootpayassetgroup", "group_level")
	ths.wantNumberList = make([]*RootPayAssetGroup, ths.wantNumberListSize)
	for i := 0; i < ths.wantNumberListSize; i++ {
		keyval := fmt.Sprintf("group%d", i+1)
		ths.wantNumberList[i] = &RootPayAssetGroup{
			AccountCnt:   _AC.ConfigInstance.GetNumber("rootpayassetgroup", keyval+"_acc_cnt"),
			OperationCnt: _AC.ConfigInstance.GetNumber("rootpayassetgroup", keyval+"_opt_cnt"),
		}
	}
	return ths
}

func (ths *RootPaymentAsset) execute() {
	_L.LoggerInstance.Info(" ** Root payment to test account ** \r\n")
	step := 0
	for {
		switch step {
		case 0: // 获取根账户中存在的资产
			step = ths.getRootAccInfo(step)
		case 1:
			step = ths.getTestName(step)
		case 2:
			step = ths.getSendAmount(step)
		case 3:
			step = ths.getIsGroup(step)
		case 4:
			step = ths.getNumberOfTestAcc(step)
		case 5:
			step = ths.getRunType(step)
		case 6:
			output := fmt.Sprintf(" >> The program will split the test accounts into groups, with each group test accounts receive asset[%s] in the following order: \r\n", ths.sendCode)
			for i, itm := range ths.wantNumberList {
				output += fmt.Sprintf("  >> Group%d : Account = %d\toperation = %d\r\n", (i + 1), itm.AccountCnt, itm.OperationCnt)
			}
			output += fmt.Sprintf("  >> Group%d : Account = %d\toperation = %d\r\n  >> and then cycle ...\r\n", (ths.wantNumberListSize + 1),
				ths.wantNumberList[0].AccountCnt, ths.wantNumberList[0].OperationCnt)
			_L.LoggerInstance.InfoPrint(output)
			step = ths.getRunType(step)
			if step != -1 {
				ths.getGroup()
				ths.groupSize = len(ths.groups)
				_L.LoggerInstance.InfoPrint(" >> Program create %d groups.\r\n", ths.groupSize)
			}
		case 7:
			step = ths.getConfirm(step)
		case 8:
			ths.run()
			return
		default:
			ths.OperationCanceled()
			return
		}
	}
}

func (ths *RootPaymentAsset) getRootAccInfo(step int) int {
	_L.LoggerInstance.InfoPrint(" >> Get root account information ...\r\n")
	err := ths.rootAccInfo.GetInfo(nil)
	if err == nil {
		noAsset := true
		step++
		for _, itm := range ths.rootAccInfo.Assets {
			if len(itm.Issuer) != 0 {
				noAsset = false
				_L.LoggerInstance.InfoPrint(" >> Asset Info: Code[%s]\tBalance[%s]\r\n", itm.Code, itm.Balance)
			}
		}
		if noAsset {
			_L.LoggerInstance.ErrorPrint(" >> Root account[%s] is NOT change trust!", ths.rootAccInfo.Address)
			step = -1
		} else {
			tCountInDb := _OP.GetSuccessTestAccCount(nil)
			ths.changeTrustCnt = int(tCountInDb - _OP.GetNoAssetCodeTestAccCount(nil))
			_L.LoggerInstance.InfoPrint(" >> Current    active     test account : %d\r\n", tCountInDb)
			_L.LoggerInstance.InfoPrint(" >> Current changed trust test account : %d\r\n", ths.changeTrustCnt)
			if ths.changeTrustCnt <= 0 {
				step = -1
			}
		}
	} else {
		step = -1
	}
	return step
}

func (ths *RootPaymentAsset) getTestName(step int) int {
	_L.LoggerInstance.InfoPrint(" >> Please enter a name for this test(lenght <= 255): ")
	ths.TestName, _ = ths.InputString()
	step++
	return step
}

func (ths *RootPaymentAsset) getSendAmount(step int) int {
	_L.LoggerInstance.InfoPrint(" >> How many asset[%s] AMOUNT you will sent to test account(1~max, 0 is cancel)? ", ths.sendCode)
	amount, b := ths.InputFloat()
	step++
	if b && amount > 0.0 {
		ths.sendAmount = fmt.Sprintf("%f", amount)
	} else {
		step = -1
	}
	return step
}

func (ths *RootPaymentAsset) getIsGroup(step int) int {
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
	return step
}

func (ths *RootPaymentAsset) getNumberOfTestAcc(step int) int {
	step++
	_L.LoggerInstance.InfoPrint(" >> How many test accounts will receive asset[%s] (1~%d, others is cancel)? ", ths.sendCode, ths.changeTrustCnt)
	num, b := ths.InputNumber()
	if b && num > 0 && num <= ths.changeTrustCnt {
		ths.sendTestAccCnt = num
	} else {
		step = -1
	}

	if step != -1 {
		_L.LoggerInstance.InfoPrint(" >> How many operations per transaction (1~100, others is cancel)? ")
		num, b = ths.InputNumber()
		if b && num > 0 && num <= 100 {
			ths.operationCnt = num
		} else {
			step = -1
		}
	}
	return step
}

func (ths *RootPaymentAsset) getRunType(step int) int {
	if ths.isGroup {
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
		} else {
			step = -1
		}
	} else {
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
	}
	return step
}

func (ths *RootPaymentAsset) getConfirm(step int) int {
	info := " >> Please confirm fellow informations:\r\n"
	info += fmt.Sprintf(" >> Test name: %s\r\n", ths.TestName)
	if ths.isGroup {
		info += fmt.Sprintf(" >> Group test: %d groups\r\n", ths.groupSize)
	} else {
		info += fmt.Sprintf(" >> Receive asset test account: %d\r\n", ths.sendTestAccCnt)
	}
	info += fmt.Sprintf(" >> %s amount: %s\r\n", ths.sendCode, ths.sendAmount)
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
	return step
}

func (ths *RootPaymentAsset) getGroup() {
	ths.getTestAccount(ths.changeTrustCnt)
	ths.groups = make([]*RootPaymentAssetGroup, 0)
	grpIdx := 0
	total := ths.changeTrustCnt
	startIndex := 0
	endIndex := 0
	left := total
	for left > 0 {
		wantIndex := grpIdx % ths.wantNumberListSize
		if left < ths.wantNumberList[wantIndex].AccountCnt {
			endIndex += left
		} else {
			endIndex += ths.wantNumberList[wantIndex].AccountCnt
		}
		size := endIndex - startIndex
		grp := &RootPaymentAssetGroup{
			ItemGroupBase: ItemGroupBase{
				TestName:             ths.TestName,
				GroupID:              grpIdx,
				RunType:              ths.runTypeSub,
				LengthOfTestAccounts: size,
				OperationCnt:         ths.wantNumberList[wantIndex].OperationCnt,
				Issuer:               ths.sendIssuer,
				Code:                 ths.sendCode,
			},
			RootAcc: ths.rootAccInfo,
			Amount:  ths.sendAmount,
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

func (ths *RootPaymentAsset) getTestAccount(num int) bool {
	_L.LoggerInstance.InfoPrint(" >> Get test account from database ...\r\n")
	ths.accInDb = _OP.GetNumberofTestAccForSendAsset(num, ths.sendCode)
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

func (ths *RootPaymentAsset) run() {
	num := 0
	if ths.isGroup {
		num = ths.changeTrustCnt
	} else {
		num = ths.sendTestAccCnt
	}

	if !ths.getTestAccount(num) {
		ths.OperationCanceled()
		return
	}

	if ths.isGroup {
		ths.groupRun()
	} else {
		ths.liteRun()
	}
}

func (ths *RootPaymentAsset) groupRun() {
	ths.wait = new(sync.WaitGroup)
	for i := 0; i < ths.groupSize; i++ {
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

	_L.LoggerInstance.InfoPrint("  >> Send asset %s to test account complete!\r\n", ths.sendCode)
}

func (ths *RootPaymentAsset) liteRun() {
	_L.LoggerInstance.InfoPrint(" >> Issuer = %s\r\n >> Code = %s\r\n >> Amount = %s\r\n", ths.sendIssuer, ths.sendCode, ths.sendAmount)
	_L.LoggerInstance.InfoPrint(" >> Get signatures and post to horizon...\r\n")

	grp := &RootPaymentAssetGroup{
		ItemGroupBase: ItemGroupBase{
			TestName:             ths.TestName,
			GroupID:              -1,
			RunType:              ths.runType,
			LengthOfTestAccounts: ths.sendTestAccCnt,
			OperationCnt:         ths.operationCnt,
			Issuer:               ths.sendIssuer,
			Code:                 ths.sendCode,
		},
		RootAcc: ths.rootAccInfo,
		Amount:  ths.sendAmount,
	}
	grp.TestAccounts = make([]*api.AccountInfo, ths.sendTestAccCnt)
	grp.AccInDb = make([]*_DB.TTestAccount, ths.sendTestAccCnt)
	copy(grp.TestAccounts, ths.testAccs)
	copy(grp.AccInDb, ths.accInDb)
	grp.Do(nil)
}
