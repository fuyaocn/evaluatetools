package menu

import (
	"jojopoper/NBi/StressTest/api"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	"math"
	"strconv"
	"sync"
)

// ItemGroupBase 分组基础定义
type ItemGroupBase struct {
	TestName             string
	GroupID              int
	TestAccounts         []*api.AccountInfo
	AccInDb              []*_DB.TTestAccount
	LengthOfTestAccounts int
	OperationCnt         int
	RunType              api.SendType
	wait                 *sync.WaitGroup
	signerList           []string
	signerListSize       int
	static               []*_DB.TStatic
	Issuer               string
	Code                 string
	ActionName           string
}

func (ths *ItemGroupBase) checkTestAccInfo() bool {
	_L.LoggerInstance.InfoPrint(" >> [ GroupIndex = %d ] Get account information from horizon ...\r\n", ths.GroupID)
	groupCnt := 500
	if groupCnt > ths.LengthOfTestAccounts {
		groupCnt = ths.LengthOfTestAccounts
	}
	ths.wait = new(sync.WaitGroup)
	cycCnt := int(math.Ceil(float64(ths.LengthOfTestAccounts) / float64(groupCnt)))
	// _L.LoggerInstance.DebugPrint(" cycCnt = %d\r\n", cycCnt)
	leftCnt := ths.LengthOfTestAccounts % groupCnt
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
					_L.LoggerInstance.ErrorPrint("  **** get account[%s] has error:\r\n%+v\r\n", acc.Address, err)
				}
			}(ths.TestAccounts[i+index])
		}
		ths.wait.Wait()
	}

	for i := 0; i < ths.LengthOfTestAccounts; i++ {
		if ths.TestAccounts[i].Status != 0 {
			return false
		}
	}
	_L.LoggerInstance.InfoPrint(" >> [ GroupIndex = %d ] Get account information from horizon complete!\r\n", ths.GroupID)
	return true
}

func (ths *ItemGroupBase) saveStatics(sockets []*net.SocketHttp) {
	if sockets == nil {
		_L.LoggerInstance.ErrorPrint("Socket is nil! Can not save Static to database\n")
		return
	}
	ths.static = make([]*_DB.TStatic, ths.signerListSize)
	for i := 0; i < ths.signerListSize; i++ {
		sk := sockets[i]
		if sk != nil {
			ths.static[i] = &_DB.TStatic{
				StaticName:   ths.TestName,
				GroupIndex:   ths.GroupID,
				ItemIndex:    i,
				StartTime:    sk.StartSend,
				CompleteTime: sk.CompleteSend,
				Success:      string(sk.Result.Result),
				Action:       ths.ActionName,
				Signature:    ths.signerList[i],
				// MainAccID:      ths.TestAccounts[i].Address,
				OperationCnt:   ths.OperationCnt,
				TransactionCnt: ths.LengthOfTestAccounts,
			}
			if sk.Result != nil && sk.Result.Extras != nil {
				ths.static[i].FailureCause = sk.Result.Extras.ResultXdr
			}
		} else {
			_L.LoggerInstance.ErrorPrint("Socket [%d] is nil! Can not save Static to database\n", i)
		}
	}
}

func (ths *ItemGroupBase) saveTestAccInfo() {
	for idx, itm := range ths.TestAccounts {
		if len(itm.Assets) > 0 {
			for _, asset := range itm.Assets {
				if asset.Issuer == ths.Issuer {
					ths.AccInDb[idx].AssetCode = asset.Code
					ths.AccInDb[idx].AssetBalance, _ = strconv.ParseFloat(asset.Balance, 64)
					break
				}
			}
		}
		ths.AccInDb[idx].Balance = itm.Balance
	}
}
