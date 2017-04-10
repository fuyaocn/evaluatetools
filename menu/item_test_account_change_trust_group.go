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

// TestAccChangeTrustGroup 测试账户分组建立信任
type TestAccChangeTrustGroup struct {
	ItemGroupBase
}

// Do 执行
func (ths *TestAccChangeTrustGroup) Do(wt *sync.WaitGroup) {
	isHorizonExecute := true
	if wt != nil {
		defer wt.Done()
	}
	if !ths.checkTestAccInfo() {
		return
	}
	ths.ActionName = "change_trust"
	ths.OperationCnt = 1

	ct := &api.ChangeTrust{}
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ]\r\n >> Issuer = %s\r\n >> Code = %s\r\n", ths.GroupID, ths.Issuer, ths.Code)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Get signatures and post to horizon...\r\n", ths.GroupID)
	ths.signerList = ct.GetSignature(ths.Issuer, ths.Code, ths.TestName, ths.TestAccounts...)
	if ths.signerList == nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Get signatures failure!\r\n", ths.GroupID)
		return
	}
	ths.signerListSize = len(ths.signerList)
	if isHorizonExecute {
		addr := _AC.ConfigInstance.GetHorizonServer() + "/transactions"
		err := ct.Send(addr, ths.RunType, ths.signerList...)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Post signatures to horizon has error :\r\n%+v\r\n", ths.GroupID, err)
		}
	} else {
		addr := _AC.ConfigInstance.GetCoreServer()
		err := ct.SendCore(addr, ths.RunType, ths.signerList...)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Post signatures to core has error :\r\n%+v\r\n", ths.GroupID, err)
		}

		time.Sleep(time.Duration(5 * time.Second))
	}

	ths.checkTestAccInfo()

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Save static datas to database...\r\n", ths.GroupID)
	if isHorizonExecute {
		ths.saveStatics(ct.GetSocket(), false)
	} else {
		ths.saveStatics(ct.GetCoreSocket(), true)
	}
	err := _OP.SaveStaticsToDB(nil, ths.static)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Save statics to database has error : \r\n %+v\r\n", ths.GroupID, err)
	}

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Update test account informations to database...\r\n", ths.GroupID)
	ths.saveTestAccInfo()
	_OP.UpdateTestAccInfoToDB(ths.AccInDb)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Change trust operation had completed!\r\n", ths.GroupID)
}

func (ths *TestAccChangeTrustGroup) saveStatics(sockets []*net.SocketHttp, isCore bool) {
	ths.ItemGroupBase.saveStatics(sockets)
	for i := 0; i < ths.signerListSize; i++ {
		ths.static[i].MainAccID = ths.TestAccounts[i].Address
		// // for core post
		if isCore {
			ths.static[i].Success = string(sockets[i].CoreResult.Result)
			if !sockets[i].CoreResult.IsCorrect() {
				ths.static[i].FailureCause = sockets[i].CoreResult.Error
			}
		}
	}
}
