package menu

import (
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	_OP "jojopoper/NBi/StressTest/operator"
	"sync"
	"time"
)

/*

	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	_OP "jojopoper/NBi/StressTest/operator"
	"sync"
*/

// RootPaymentAssetGroup 根账户发送资产分组定义
type RootPaymentAssetGroup struct {
	ItemGroupBase
	RootAcc *api.AccountInfo
	Amount  string
}

// Do 执行
func (ths *RootPaymentAssetGroup) Do(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	if !ths.checkTestAccInfo() {
		return
	}
	ths.ActionName = "root_pay_asset"
	pay := &api.Payment{}
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ]\r\n >> Issuer = %s\r\n >> Code = %s\r\n >> Amount = %s\r\n", ths.GroupID, ths.Issuer, ths.Code, ths.Amount)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Get signatures and send to core...\r\n", ths.GroupID)
	// ths.signerList = pay.GetSignature(ths.Issuer, ths.Code, ths.Amount, []*api.AccountInfo{ths.RootAcc}, ths.TestAccounts)
	ths.getSignatures(pay)
	if ths.signerList == nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Get signatures failure!\r\n", ths.GroupID)
		return
	}
	ths.signerListSize = len(ths.signerList)

	_L.LoggerInstance.DebugPrint(" >>>> sign list size = %d\r\n", ths.signerListSize)

	addr := _AC.ConfigInstance.GetCoreServer()
	err := pay.SendCore(addr, ths.RunType, ths.signerList...)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Send signatures to core has error :\r\n%+v\r\n", ths.GroupID, err)
	}

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Waiting 10s for core executable...\r\n", ths.GroupID)
	time.Sleep(time.Duration(10 * time.Second))
	ths.checkTestAccInfo()

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Save static datas to database...\r\n", ths.GroupID)
	ths.saveStatics(pay.GetCoreSocket())
	err = _OP.SaveStaticsToDB(nil, ths.static)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Save statics to database has error : \r\n %+v\r\n", ths.GroupID, err)
	}

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Update test account informations to database...\r\n", ths.GroupID)
	ths.saveTestAccInfo()
	_OP.UpdateTestAccInfoToDB(ths.AccInDb)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Root asset payment operation had completed!\r\n", ths.GroupID)
}

func (ths *RootPaymentAssetGroup) saveStatics(sockets []*net.SocketHttp) {
	ths.ItemGroupBase.saveStatics(sockets)
	for i := 0; i < ths.signerListSize; i++ {
		ths.static[i].MainAccID = ths.RootAcc.Address
		ths.static[i].Success = string(sockets[i].CoreResult.Result)
		if !sockets[i].CoreResult.IsCorrect() {
			ths.static[i].FailureCause = sockets[i].CoreResult.Error
		}
	}
}

func (ths *RootPaymentAssetGroup) getSignatures(pay *api.Payment) {
	total := ths.LengthOfTestAccounts
	step := ths.OperationCnt
	if ths.OperationCnt > total {
		ths.OperationCnt = total
	}
	start := 0
	end := 0
	left := total
	ths.updateRootAccInfo()
	src := []*api.AccountInfo{ths.RootAcc}
	ths.signerList = make([]string, 0)
	index := 0
	for left > 0 {
		if left < step {
			step = left
		}
		end += step
		// _L.LoggerInstance.DebugPrint("Start = %d\tEnd = %d\tStep = %d\tLeft = %d\r\n", start, end, step, left)
		accInfos := make([]*api.AccountInfo, step)
		copy(accInfos, ths.TestAccounts[start:end])
		b64s := pay.GetSignature(ths.Issuer, ths.Code, ths.Amount, ths.TestName, src, accInfos)
		// _L.LoggerInstance.DebugPrint("[%d] b64s(length = %d) = %+v\r\n", index, len(b64s), b64s)
		ths.signerList = append(ths.signerList, b64s...)
		start = end
		left = total - end
		index++
	}
}

func (ths *RootPaymentAssetGroup) updateRootAccInfo() {
	tmp := new(api.AccountInfo)
	tmp.Init(ths.RootAcc.Address, ths.RootAcc.Secret)
	err := tmp.GetInfo(nil)
	if err == nil {
		ths.RootAcc = tmp
		_L.LoggerInstance.InfoPrint(" ## Root Acc Info = [\n%+v\n]\n", ths.RootAcc)
		for idx, itm := range ths.RootAcc.Assets {
			_L.LoggerInstance.InfoPrint(" ## Root Acc Asset Info [%d]= [\n%+v\n]\n", idx, itm)
		}
	} else {
		_L.LoggerInstance.ErrorPrint(" **** Update Root account information has error :\n%+v\n", err)
	}
}
