package menu

import (
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	_OP "jojopoper/NBi/StressTest/operator"
	"math"
	"strconv"
	"sync"
)

// GroupItem 分组中的一个单体
type GroupItem struct {
	Parent         *MainAccountGroup
	TestName       string
	GroupID        int
	ItemIndex      int
	MainAccounts   *api.AccountInfo
	OperationCnt   int
	TestAccountCnt int
	testAccs       []*_DB.TTestAccount
	StartBalance   string
	signerList     []string
	signerListSize int
	static         []*_DB.TStatic
}

// DoneCreate 执行生成测试账户操作
func (ths *GroupItem) DoneCreate(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	ths.getSignerList()

	postFunc := &api.PostSend2Horizon{}
	addr := _AC.ConfigInstance.GetHorizonServer() + "/transactions"
	err := postFunc.Send(addr, api.Linear, ths.signerList...)
	if err == nil {
		bala, _ := strconv.ParseFloat(ths.StartBalance, 64)
		err = _OP.SaveTestAccToDatabase(ths.testAccs, bala, "T", nil)
	}

	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Create test account has error [grp:%d;itm:%d] : \r\n %+v\r\n", ths.GroupID, ths.ItemIndex, err)
	}
	ths.saveStatics(postFunc.GetSocket())

	err = _OP.SaveStaticsToDB(nil, ths.static)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** Save statics to database has error [grp:%d;itm:%d] : \r\n %+v\r\n", ths.GroupID, ths.ItemIndex, err)
	}
}

func (ths *GroupItem) getSignerList() {
	indexgrp := (ths.GroupID*ths.Parent.LengthOfMainAccounts + ths.ItemIndex) * ths.TestAccountCnt
	ths.testAccs = _OP.GetNewTestAccs(ths.TestAccountCnt, indexgrp, ths.GroupID, ths.ItemIndex)
	ths.signerListSize = int(math.Ceil(float64(ths.TestAccountCnt) / float64(ths.OperationCnt)))
	ths.signerList = make([]string, ths.signerListSize)
	// _L.LoggerInstance.Debug(" signerListSize = %d\r\n", ths.signerListSize)
	// _L.LoggerInstance.Debug(" testAccs size = %d\r\n", len(ths.testAccs))
	startIndex := 0
	if ths.TestAccountCnt < ths.OperationCnt {
		ths.OperationCnt = ths.TestAccountCnt
	}
	step := ths.OperationCnt
	total := ths.TestAccountCnt
	accInfo := &api.CreateAccount{}
	for i := 0; i < ths.signerListSize; i++ {
		// _L.LoggerInstance.Debug(" startIndex = %d\t endIndex = %d\t total = %d\t step = %d\r\n", startIndex, step+i*ths.OperationCnt, total, step)
		taccs := ths.testAccs[startIndex : step+i*ths.OperationCnt]
		startIndex = step + i*ths.OperationCnt
		total -= step
		if total < ths.OperationCnt {
			step = total
		}
		ths.signerList[i] = accInfo.GetSignature(ths.MainAccounts, ths.StartBalance, ths.TestName, ths.getAccountIDs(taccs)...)
	}
}

func (ths *GroupItem) getAccountIDs(src []*_DB.TTestAccount) (ret []string) {
	if src == nil {
		return
	}
	ret = make([]string, 0)
	for _, itm := range src {
		ret = append(ret, itm.AccountID)
	}
	return
}

func (ths *GroupItem) saveStatics(sockets []*net.SocketHttp) {
	ths.static = make([]*_DB.TStatic, ths.signerListSize)
	for i := 0; i < ths.signerListSize; i++ {
		ths.static[i] = &_DB.TStatic{
			StaticName:     ths.TestName,
			GroupIndex:     ths.GroupID,
			ItemIndex:      ths.ItemIndex,
			StartTime:      sockets[i].StartSend,
			CompleteTime:   sockets[i].CompleteSend,
			Success:        string(sockets[i].Result.Result),
			Action:         "active_test_account",
			Signature:      ths.signerList[i],
			MainAccID:      ths.MainAccounts.Address,
			OperationCnt:   ths.OperationCnt,
			TransactionCnt: ths.signerListSize,
		}
		if sockets[i].Result.Extras != nil {
			ths.static[i].FailureCause = sockets[i].Result.Extras.ResultXdr
		}
	}
}
