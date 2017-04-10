package api

import _b "github.com/stellar/go/build"

// OrderInfo 挂单信息
type OrderInfo struct {
	Code   string
	Issuer string
}

// GetAsset 得到Asset结构
func (ths *OrderInfo) GetAsset() _b.Asset {
	if len(ths.Issuer) == 0 {
		return _b.NativeAsset()
	}
	return _b.CreditAsset(ths.Code, ths.Issuer)
}

// AccMakeOrder 挂单信息
type AccMakeOrder struct {
	SrcAcc     *AccountInfo
	buyInfo    []*OrderInfo
	sellInfo   []*OrderInfo
	LengthInfo int
	price      []string
	amount     []string
}

// SetOrderInfo 设置买卖信息
func (ths *AccMakeOrder) SetOrderInfo(buy, sell *OrderInfo, price, amount string) {
	if ths.buyInfo == nil {
		ths.buyInfo = make([]*OrderInfo, 0)
	}
	if ths.sellInfo == nil {
		ths.sellInfo = make([]*OrderInfo, 0)
	}
	if ths.price == nil {
		ths.price = make([]string, 0)
	}
	if ths.amount == nil {
		ths.amount = make([]string, 0)
	}
	ths.buyInfo = append(ths.buyInfo, buy)
	ths.sellInfo = append(ths.sellInfo, sell)
	ths.amount = append(ths.amount, amount)
	ths.price = append(ths.price, price)
	ths.LengthInfo++
}

// GetSellInfo 得到卖出信息
func (ths *AccMakeOrder) GetSellInfo(i int) *OrderInfo {
	if ths.sellInfo == nil || i >= ths.LengthInfo {
		return nil
	}
	return ths.sellInfo[i]
}

// GetBuyInfo 得到买入信息
func (ths *AccMakeOrder) GetBuyInfo(i int) *OrderInfo {
	if ths.buyInfo == nil || i >= ths.LengthInfo {
		return nil
	}
	return ths.buyInfo[i]
}

// GetAmount 得到数量信息
func (ths *AccMakeOrder) GetAmount(i int) _b.Amount {
	if ths.amount == nil || i >= ths.LengthInfo {
		return ""
	}
	return _b.Amount(ths.amount[i])
}

// GetPrice 得到价格信息
func (ths *AccMakeOrder) GetPrice(i int) _b.Price {
	if ths.price == nil || i >= ths.LengthInfo {
		return ""
	}
	return _b.Price(ths.price[i])
}
