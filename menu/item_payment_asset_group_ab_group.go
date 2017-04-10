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

/*

	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	_OP "jojopoper/NBi/StressTest/operator"
	"math"
	"strconv"
	"sync"
*/

// PayAssetGroupABGroup 分组账户发送资产分组定义
type PayAssetGroupABGroup struct {
	ItemGroupBase
	Sender            []*api.AccountInfo
	SenderDB          []*_DB.TTestAccount
	LengthOfSenderAcc int
	Amount            string
}

// Do 执行
func (ths *PayAssetGroupABGroup) Do(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	if !ths.checkTestAccInfo() || !ths.checkSenderAccInfo() {
		return
	}
	ths.ActionName = "group_a_b_pay_asset"
	pay := &api.Payment{}
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ]\r\n >> Issuer = %s\r\n >> Code = %s\r\n >> Amount = %s\r\n", ths.GroupID, ths.Issuer, ths.Code, ths.Amount)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Get signatures and send to core...\r\n", ths.GroupID)
	ths.signerList = pay.GetSignature(ths.Issuer, ths.Code, ths.Amount, ths.TestName, ths.Sender, ths.TestAccounts)
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

	ths.checkTestAccInfo()
	ths.checkSenderAccInfo()

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Save static datas to database...\r\n", ths.GroupID)
	ths.saveStatics(pay.GetCoreSocket())
	err = _OP.SaveStaticsToDB(nil, ths.static)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("  **** [ GroupID = %d ] Save statics to database has error : \r\n %+v\r\n", ths.GroupID, err)
	}

	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Update test account informations to database...\r\n", ths.GroupID)
	ths.saveTestAccInfo()
	_OP.UpdateTestAccInfoToDB(ths.AccInDb)
	ths.saveSenderAccInfo()
	_OP.UpdateTestAccInfoToDB(ths.SenderDB)
	_L.LoggerInstance.InfoPrint(" >> [ GroupID = %d ] Root asset payment operation had completed!\r\n", ths.GroupID)
}

func (ths *PayAssetGroupABGroup) checkSenderAccInfo() bool {
	_L.LoggerInstance.InfoPrint(" >> [ GroupIndex = %d ] Get account information from horizon ...\r\n", ths.GroupID)
	groupCnt := 500
	if groupCnt > ths.LengthOfSenderAcc {
		groupCnt = ths.LengthOfSenderAcc
	}
	ths.wait = new(sync.WaitGroup)
	cycCnt := int(math.Ceil(float64(groupCnt) / float64(ths.LengthOfSenderAcc)))
	// _L.LoggerInstance.DebugPrint(" cycCnt = %d\r\n", cycCnt)
	leftCnt := groupCnt % ths.LengthOfSenderAcc
	// _L.LoggerInstance.DebugPrint(" leftCnt = %d\r\n", leftCnt)
	for idx := 0; idx < cycCnt; idx++ {
		index := idx * groupCnt
		// _L.LoggerInstance.DebugPrint(" index = %d\r\n", index)
		if idx == cycCnt-1 && leftCnt > 0 {
			groupCnt = leftCnt
		}
		// _L.LoggerInstance.DebugPrint(" groupCnt = %d\r\n", groupCnt)

		for i := 0; i < groupCnt; i++ {
			ths.wait.Add(1)
			go func(acc *api.AccountInfo) {
				err := acc.GetInfo(ths.wait)
				if err != nil {
					_L.LoggerInstance.ErrorPrint("  **** get account[%s] has error:\r\n%+v\r\n", acc.Address)
				}
			}(ths.Sender[i+index])
		}
		ths.wait.Wait()
	}

	for i := 0; i < ths.LengthOfSenderAcc; i++ {
		if ths.Sender[i].Status != 0 {
			return false
		}
	}
	_L.LoggerInstance.InfoPrint(" >> [ GroupIndex = %d ] Get account information from horizon complete!\r\n", ths.GroupID)
	return true
}

func (ths *PayAssetGroupABGroup) saveStatics(sockets []*net.SocketHttp) {
	ths.ItemGroupBase.saveStatics(sockets)
	for i := 0; i < ths.signerListSize; i++ {
		ths.static[i].MainAccID = ths.Sender[i].Address
		ths.static[i].Success = string(sockets[i].CoreResult.Result)
		if !sockets[i].CoreResult.IsCorrect() {
			ths.static[i].FailureCause = sockets[i].CoreResult.Error
		}
	}
}

func (ths *PayAssetGroupABGroup) saveSenderAccInfo() {
	for idx, itm := range ths.Sender {
		if len(itm.Assets) > 0 {
			for _, asset := range itm.Assets {
				if asset.Issuer == ths.Issuer {
					ths.SenderDB[idx].AssetCode = asset.Code
					ths.SenderDB[idx].AssetBalance, _ = strconv.ParseFloat(asset.Balance, 64)
					break
				}
			}
		}
		ths.SenderDB[idx].Balance = itm.Balance
	}
}
