package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	_OP "jojopoper/NBi/StressTest/operator"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// PayAssetGroupAB 分组多对一,或一对多，或一对一发送资产
type PayAssetGroupAB struct {
	SubItem
	TestName             string
	groupAdb             []*_DB.TTestAccount // 发送资产的账户组A
	groupASize           int                 // 发送资产的账户组A长度
	groupBdb             []*_DB.TTestAccount // 发送资产的账户组B
	groupBSize           int                 // 发送资产的账户组B长度
	changeTrustCnt       int                 // 数据库中已经信任资产的账户
	assetBalanceValidCnt int                 // 数据库中有资产余额的账户
	sendAmount           string              // 发送的金额
	sendCode             string              // 发送的code
	sendIssuer           string              // 发送的网关地址
	runType              api.SendType
	runTypeSub           api.SendType
	groups               []*PayAssetGroupABGroup
	groupSize            int
	wait                 *sync.WaitGroup
}

// InitMenu 初始化
func (ths *PayAssetGroupAB) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "GroupA:B Payment"
	ths.sendCode = _AC.ConfigInstance.GetString("issueraccountinfo", "code")
	ths.sendIssuer = _AC.ConfigInstance.GetString("issueraccountinfo", "id")
	return ths
}

func (ths *PayAssetGroupAB) execute() {
	_L.LoggerInstance.Info(" ** Test account multi-single payment asset ** \r\n")
	step := 0
	for {
		switch step {
		case 0:
			step = ths.getTestName(step)
		case 1:
			step = ths.getSendAmount(step)
		case 2:
			step = ths.getGroupACnt(step)
		case 3:
			step = ths.getGroupBCnt(step)
		case 4:
			step = ths.getRunType(step)
			if step != -1 {
				step = ths.getGroup(step)
			}
			// if step != -1 {
			// 	_L.LoggerInstance.DebugPrint(" >> groups siz = %d\r\n", ths.groupSize)
			// 	for i, itm := range ths.groups {
			// 		_L.LoggerInstance.DebugPrint(" >> [%d]\tgroupsInfo = %+v\r\n", i, itm)
			// 		_L.LoggerInstance.DebugPrint(" >> [%d]\tSenderDB = %+v\r\n", i, itm.SenderDB)
			// 		_L.LoggerInstance.DebugPrint(" >> [%d]\t  Sender = %+v\r\n", i, itm.Sender)
			// 		_L.LoggerInstance.DebugPrint(" >> [%d]\t AccInDb = %+v\r\n", i, itm.AccInDb)
			// 		_L.LoggerInstance.DebugPrint(" >> [%d]\t TestAccounts = %+v\r\n\r\n", i, itm.TestAccounts)
			// 	}

			// }
		case 5:
			step = ths.getConfirm(step)
		case 6:
			ths.run()
			return
		default:
			ths.OperationCanceled()
			return
		}
	}
}

func (ths *PayAssetGroupAB) getTestName(step int) int {
	_L.LoggerInstance.InfoPrint(" >> Please enter a name for this test(lenght <= 255): ")
	ths.TestName, _ = ths.InputString()
	step++
	return step
}

func (ths *PayAssetGroupAB) getGroupACnt(step int) int {
	step++
	tCountInDb := _OP.GetSuccessTestAccCount(nil)
	ths.changeTrustCnt = int(tCountInDb - _OP.GetNoAssetCodeTestAccCount(nil))
	ths.assetBalanceValidCnt = int(_OP.GetAssetBalanceValidTestAccCount(nil, ths.sendAmount))
	_L.LoggerInstance.InfoPrint(" >>            Active test account : %d\r\n", tCountInDb)
	_L.LoggerInstance.InfoPrint(" >>     Changed trust test account : %d\r\n", ths.changeTrustCnt)
	_L.LoggerInstance.InfoPrint(" >> Asset balance > %s test account : %d\r\n", ths.sendAmount, ths.assetBalanceValidCnt)
	if ths.changeTrustCnt <= 0 {
		step = -1
	}
	_L.LoggerInstance.InfoPrint(" >> How many test account in send-group A (1~%d, 0 is cancel) ? ", ths.assetBalanceValidCnt)
	num, b := ths.InputNumber()
	if b && num > 0 && num <= ths.assetBalanceValidCnt {
		ths.groupASize = num
	} else {
		step = -1
	}
	return step
}

func (ths *PayAssetGroupAB) getGroupBCnt(step int) int {
	step++
	_L.LoggerInstance.InfoPrint(" >> Please enter the group ratio (A:B) , group B is receive asset test account group : ")
	num, b := ths.InputFloat()
	if b && num > 0 {
		if int(math.Ceil(float64(ths.groupASize)/num)) == int(float64(ths.groupASize)/num) {
			ths.groupBSize = int(float64(ths.groupASize) / num)
			if ths.changeTrustCnt >= ths.groupASize+ths.groupBSize {
				_L.LoggerInstance.InfoPrint(" >> Group B size : %d\r\n", ths.groupBSize)
			} else {
				_L.LoggerInstance.ErrorPrint(" >> Ratio is error! Group A size + group B size is out of total changed trust account number!\r\n")
				step = -1
			}
		} else {
			_L.LoggerInstance.ErrorPrint(" >> Ratio is error! Group B size is not integer!\r\n")
			step = -1
		}
	}
	return step
}

func (ths *PayAssetGroupAB) getSendAmount(step int) int {
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

func (ths *PayAssetGroupAB) getRunType(step int) int {
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
	return step
}

func (ths *PayAssetGroupAB) getConfirm(step int) int {
	info := " >> Please confirm fellow informations:\r\n"
	info += fmt.Sprintf(" >>                         Test name : %s\r\n", ths.TestName)
	info += fmt.Sprintf(" >>    Send asset group(group A) size : %d\r\n", ths.groupASize)
	info += fmt.Sprintf(" >> receive asset group(group B) size : %d\r\n", ths.groupBSize)
	info += fmt.Sprintf(" >>                        %s amount : %s\r\n", ths.sendCode, ths.sendAmount)
	switch ths.runType {
	case api.Linear:
		info += " >>                           runType : Linear\r\n"
	case api.Multiple:
		info += " >>                           runType : Multiple\r\n"
	}
	switch ths.runTypeSub {
	case api.Linear:
		info += " >>                     group runType : Linear\r\n"
	case api.Multiple:
		info += " >>                     group runType : Multiple\r\n"
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

func (ths *PayAssetGroupAB) getGroup(step int) int {
	if ths.getGroupFromDB() {
		splitANum := 0
		splitBNum := 0
		if ths.groupASize >= ths.groupBSize {
			splitANum = ths.groupASize / ths.groupBSize
			splitBNum = 1
		} else {
			splitANum = 1
			splitBNum = ths.groupBSize / ths.groupASize
		}

		ths.groupSize = ths.groupASize / splitANum
		ths.groups = make([]*PayAssetGroupABGroup, 0)
		startAindex := 0
		endAindex := splitANum
		startBindex := 0
		endBindex := splitBNum
		for i := 0; i < ths.groupSize; i++ {
			grp := &PayAssetGroupABGroup{
				ItemGroupBase: ItemGroupBase{
					TestName:             ths.TestName,
					GroupID:              i,
					RunType:              ths.runTypeSub,
					LengthOfTestAccounts: splitBNum,
					OperationCnt:         splitBNum,
					Issuer:               ths.sendIssuer,
					Code:                 ths.sendCode,
				},
				Amount:            ths.sendAmount,
				LengthOfSenderAcc: splitANum,
			}
			grp.Sender = make([]*api.AccountInfo, 0)
			grp.SenderDB = make([]*_DB.TTestAccount, 0)
			grp.TestAccounts = make([]*api.AccountInfo, 0)
			grp.AccInDb = make([]*_DB.TTestAccount, 0)

			for j := startAindex; j < endAindex; j++ {
				accInfo := &api.AccountInfo{}
				accInfo.Init(ths.groupAdb[j].AccountID, ths.groupAdb[j].SecertAddr)
				grp.Sender = append(grp.Sender, accInfo)
				grp.SenderDB = append(grp.SenderDB, ths.groupAdb[j])
			}
			startAindex = endAindex
			endAindex += splitANum

			for j := startBindex; j < endBindex; j++ {
				accInfo := &api.AccountInfo{}
				accInfo.Init(ths.groupBdb[j].AccountID, ths.groupBdb[j].SecertAddr)
				grp.TestAccounts = append(grp.TestAccounts, accInfo)
				grp.AccInDb = append(grp.AccInDb, ths.groupBdb[j])
			}
			startBindex = endBindex
			endBindex += splitBNum

			ths.groups = append(ths.groups, grp)
		}
	} else {
		step = -1
	}
	return step
}

func (ths *PayAssetGroupAB) getGroupFromDB() bool {
	balanceUp0 := ths.randSlice(_OP.GetNumberofTestAccForAssetBalance(ths.assetBalanceValidCnt, ths.sendAmount, ths.sendCode))
	balanceEqu0 := _OP.GetNumberofTestAccForIs0AssetBalance(ths.sendCode)

	lenBalanUp0 := len(balanceUp0)
	if lenBalanUp0 < ths.groupASize {
		_L.LoggerInstance.ErrorPrint(" >> Number of 'asset balance > 0' is not greater than group A size(%d)!\r\n", ths.groupASize)
		return false
	}

	ths.groupAdb = make([]*_DB.TTestAccount, ths.groupASize)
	copy(ths.groupAdb, balanceUp0[0:ths.groupASize])
	lenBalanUp0 = lenBalanUp0 - ths.groupASize
	if lenBalanUp0 > 0 {
		balanceEqu0 = append(balanceEqu0, balanceUp0[ths.groupASize:]...)
	}
	balanceEqu0 = ths.randSlice(balanceEqu0)
	lenBalanEqu0 := len(balanceEqu0)
	if lenBalanEqu0 < ths.groupBSize {
		_L.LoggerInstance.ErrorPrint(" >> Number of 'changed trust' is not greater than group B size(%d)!\r\n", ths.groupBSize)
		return false
	}
	ths.groupBdb = make([]*_DB.TTestAccount, ths.groupBSize)
	copy(ths.groupBdb, balanceEqu0[0:ths.groupBSize])
	return true
}

func (ths *PayAssetGroupAB) randSlice(src []*_DB.TTestAccount) []*_DB.TTestAccount {
	srcLen := len(src)
	if srcLen != 0 {
		rr := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < srcLen; i++ {
			index := rr.Intn(srcLen)
			tmp := src[i]
			src[i] = src[index]
			src[index] = tmp
		}
	}
	return src
}

func (ths *PayAssetGroupAB) run() {
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

	_L.LoggerInstance.InfoPrint("  >> Group send asset %s to test account complete!\r\n", ths.sendCode)
}
