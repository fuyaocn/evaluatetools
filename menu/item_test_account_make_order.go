package menu

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
	_AC "jojopoper/NBi/StressTest/appconf"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	_OP "jojopoper/NBi/StressTest/operator"
	"math/rand"
	"time"
)

// MakeOrderAB 创建挂单
type MakeOrderAB struct {
	SubItem
	ItemGroupBase
	tCountInDb           int                 // 数据库中已经激活的账户
	changeTrustCnt       int                 // 数据库中已经信任资产的账户
	assetBalanceValidCnt int                 // 数据库中有资产余额的账户
	buydb                []*_DB.TTestAccount // 发送资产的账户组A
	buySize              int                 // 发送资产的账户组A长度
	selldb               []*_DB.TTestAccount // 发送资产的账户组B
	sellSize             int                 // 发送资产的账户组B长度
	priceHigh            float64
	priceLow             float64
	priceSpan            float64
	amountHigh           float64
	amountLow            float64
	amountSpan           float64
	randCalc             *rand.Rand
}

// InitMenu 初始化
func (ths *MakeOrderAB) InitMenu(parent SubItemInterface, key string) SubItemInterface {
	ths.SubItem.InitMenu(parent, key)
	ths.Exec = ths.execute
	ths.title = "Make Order"
	ths.Code = _AC.ConfigInstance.GetString("issueraccountinfo", "code")
	ths.Issuer = _AC.ConfigInstance.GetString("issueraccountinfo", "id")
	ths.priceHigh = _AC.ConfigInstance.GetFloat("ordertest", "price_high")
	ths.priceLow = _AC.ConfigInstance.GetFloat("ordertest", "price_low")
	ths.priceSpan = ths.priceHigh - ths.priceLow
	ths.amountHigh = _AC.ConfigInstance.GetFloat("ordertest", "amount_high")
	ths.amountLow = _AC.ConfigInstance.GetFloat("ordertest", "amount_low")
	ths.amountSpan = ths.amountHigh - ths.amountLow
	ths.RunType = api.Linear
	ths.randCalc = rand.New(rand.NewSource(time.Now().UnixNano()))
	ths.ActionName = "make_order_random"
	ths.GroupID = 0
	return ths
}

func (ths *MakeOrderAB) execute() {
	_L.LoggerInstance.Info(" ** Test account change trust ** \r\n")
	step := 0
	for {
		switch step {
		case 0:
			step = ths.getInformation(step)
		case 1:
			step = ths.getTestName(step)
		case 2:
			step = ths.getAccBuy(step)
		case 3:
			step = ths.getAccSell(step)
		case 4:
			step = ths.getGroup(step)
		case 5:
			step = ths.typeModeSelect(step)
		case 6:
			ths.run()
			return
		default:
			ths.OperationCanceled()
			return
		}
	}
}

func (ths *MakeOrderAB) getInformation(step int) int {
	step++
	countInDb := int(_OP.GetTestAccountCount(nil))
	ths.tCountInDb = int(_OP.GetSuccessTestAccCount(nil))
	ths.changeTrustCnt = ths.tCountInDb - int(_OP.GetNoAssetCodeTestAccCount(nil))
	ths.assetBalanceValidCnt = int(_OP.GetAssetBalanceValidTestAccCount(nil, fmt.Sprintf("%f", ths.amountHigh)))
	_L.LoggerInstance.InfoPrint(" >>             Total test account count : %d \r\n", countInDb)
	_L.LoggerInstance.InfoPrint(" >>            Active test account count : %d \r\n", ths.tCountInDb)
	_L.LoggerInstance.InfoPrint(" >>          Inactive test account count : %d \r\n", countInDb-ths.tCountInDb)
	_L.LoggerInstance.InfoPrint(" >>     Changed trust test account count : %d \r\n", ths.changeTrustCnt)
	_L.LoggerInstance.InfoPrint(" >> Asset balance > %f test account count : %d \r\n", ths.amountHigh, ths.assetBalanceValidCnt)
	return step
}

func (ths *MakeOrderAB) getTestName(step int) int {
	_L.LoggerInstance.InfoPrint(" >> Please enter a name for this test(lenght <= 255): ")
	ths.TestName, _ = ths.InputString()
	step++
	return step
}

func (ths *MakeOrderAB) getAccBuy(step int) int {
	step++
	_L.LoggerInstance.InfoPrint(" >> How many account to making buy order(1~%d, 0 is cancel)? ", ths.tCountInDb)
	num, b := ths.InputNumber()
	if b && num > 0 && num <= ths.tCountInDb {
		ths.buySize = num
	} else {
		step = -1
	}
	return step
}

func (ths *MakeOrderAB) getAccSell(step int) int {
	step++
	_L.LoggerInstance.InfoPrint(" >> How many account to making sell order(1~%d, 0 is cancel)? ", ths.assetBalanceValidCnt)
	num, b := ths.InputNumber()
	if b && num > 0 && num <= ths.assetBalanceValidCnt {
		ths.sellSize = num
	} else {
		step = -1
	}
	return step
}

func (ths *MakeOrderAB) getGroup(step int) int {
	_L.LoggerInstance.InfoPrint(" >> Make order random group ...\n")
	step++
	step = ths.getGroupFromDB(step)
	if step != -1 {
		ths.AccInDb = make([]*_DB.TTestAccount, 0)
		ths.AccInDb = append(ths.AccInDb, ths.selldb...)
		ths.AccInDb = append(ths.AccInDb, ths.buydb...)
		ths.AccInDb = ths.randSlice(ths.AccInDb)
		ths.LengthOfTestAccounts = ths.buySize + ths.sellSize
		ths.TestAccounts = make([]*api.AccountInfo, ths.LengthOfTestAccounts)
		for i := 0; i < ths.LengthOfTestAccounts; i++ {
			ths.TestAccounts[i] = &api.AccountInfo{}
			ths.TestAccounts[i].Init(ths.AccInDb[i].AccountID, ths.AccInDb[i].SecertAddr)
		}
	}
	return step
}

func (ths *MakeOrderAB) getGroupFromDB(step int) int {
	up0Acc := _OP.GetNumberofTestAccForAssetBalance(ths.assetBalanceValidCnt, fmt.Sprintf("%f", ths.amountHigh), ths.Code)
	up0Acc = ths.randSlice(up0Acc)
	ths.selldb = make([]*_DB.TTestAccount, ths.sellSize)
	copy(ths.selldb, up0Acc[0:ths.sellSize])
	left := ths.assetBalanceValidCnt - ths.sellSize
	fromtoAcc := _OP.GetNumberofTestAccForAssetBalanceFromTo(fmt.Sprintf("%f", ths.amountHigh), "0", ths.Code)
	if len(fromtoAcc)+left < ths.buySize {
		step = -1
		_L.LoggerInstance.ErrorPrint(" >> There are not enough test accounts make buy order!\r\n")
	} else {
		if left > 0 {
			fromtoAcc = append(fromtoAcc, up0Acc[ths.sellSize:]...)
		}
		fromtoAcc = ths.randSlice(fromtoAcc)
		ths.buydb = make([]*_DB.TTestAccount, ths.buySize)
		copy(ths.buydb, fromtoAcc[0:ths.buySize])
	}
	return step
}

func (ths *MakeOrderAB) randSlice(src []*_DB.TTestAccount) []*_DB.TTestAccount {
	srcLen := len(src)
	if srcLen != 0 {
		for i := 0; i < srcLen; i++ {
			index := ths.randCalc.Intn(srcLen)
			tmp := src[i]
			src[i] = src[index]
			src[index] = tmp
		}
	}
	return src
}

func (ths *MakeOrderAB) typeModeSelect(step int) int {
	_L.LoggerInstance.InfoPrint(" >> What type you want to run (1.Linear 2.Multiple others:Cancel)? : ")
	selectNumber, b := ths.InputNumber()
	if b {
		step++
		switch selectNumber {
		case 1:
			ths.RunType = api.Linear
		case 2:
			ths.RunType = api.Multiple
		default:
			step = -1
		}

	} else {
		step = -1
	}
	return step
}

func (ths *MakeOrderAB) run() {
	_L.LoggerInstance.InfoPrint(" >> Checking test account information ...\r\n")
	if ths.checkTestAccInfo() {
		ths.signerListSize = ths.LengthOfTestAccounts
		ths.signerList = make([]string, 0)
		order := &api.MakeOrder{}
		for i := 0; i < ths.LengthOfTestAccounts; i++ {
			accOrder := &api.AccMakeOrder{}
			accOrder.SrcAcc = ths.TestAccounts[i]
			buy := &api.OrderInfo{}
			sell := &api.OrderInfo{}
			if ths.AccInDb[i].AssetBalance >= ths.amountHigh {
				sell.Code = ths.Code
				sell.Issuer = ths.Issuer
			} else {
				buy.Code = ths.Code
				buy.Issuer = ths.Issuer
			}
			accOrder.SetOrderInfo(buy, sell, ths.getRandomPrice(), ths.getRandomAmount())
			ths.signerList = append(ths.signerList, order.GetSignature(ths.TestName, accOrder)...)
		}

		_L.LoggerInstance.InfoPrint(" >> Random send signatures ...\r\n")
		addr := _AC.ConfigInstance.GetCoreServer()
		err := order.SendCore(addr, ths.RunType, ths.signerList...)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Send signatures to core has error :\r\n%+v\r\n", err)
		}

		ths.checkTestAccInfo()

		_L.LoggerInstance.InfoPrint(" >> Save static datas to database...\r\n")
		ths.saveStatics(order.GetCoreSocket())
		err = _OP.SaveStaticsToDB(nil, ths.static)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Save statics to database has error : \r\n %+v\r\n", err)
		}

		_L.LoggerInstance.InfoPrint(" >> Update test account informations to database...\r\n")
		ths.saveTestAccInfo()
		_OP.UpdateTestAccInfoToDB(ths.AccInDb)
		_L.LoggerInstance.InfoPrint(" >> Make order operation had completed!\r\n")
	} else {
		_L.LoggerInstance.InfoPrint(" >> Checking test account information failure!\r\n")
	}
}

func (ths *MakeOrderAB) getRandomPrice() string {
	if ths.priceSpan <= 0 {
		return fmt.Sprintf("%f", ths.priceHigh)
	}
	r := ths.randCalc.Intn(int(ths.priceSpan * 1000))
	return fmt.Sprintf("%f", float64(r)/float64(1000.0)+ths.priceLow)
}

func (ths *MakeOrderAB) getRandomAmount() string {
	if ths.amountSpan <= 0 {
		return fmt.Sprintf("%f", ths.amountHigh)
	}
	r := ths.randCalc.Intn(int(ths.amountSpan * 1000))
	return fmt.Sprintf("%f", float64(r)/float64(1000.0)+ths.amountLow)
}

func (ths *MakeOrderAB) saveStatics(sockets []*net.SocketHttp) {
	ths.ItemGroupBase.saveStatics(sockets)
	for i := 0; i < ths.signerListSize; i++ {
		ths.static[i].MainAccID = ths.TestAccounts[i].Address
		ths.static[i].Success = string(sockets[i].CoreResult.Result)
		if !sockets[i].CoreResult.IsCorrect() {
			ths.static[i].FailureCause = sockets[i].CoreResult.Error
		}
	}
}
