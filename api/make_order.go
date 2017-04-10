package api

import (
	_b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// MakeOrder 挂单
type MakeOrder struct {
	SendToCore
	PostSend2Horizon
}

// GetSignature 得到签名
func (ths *MakeOrder) GetSignature(memo string, src ...*AccMakeOrder) (retB64 []string) {
	if len(src) == 0 {
		return nil
	}

	retB64 = make([]string, 0)
	lenMemo := len(memo)

	for _, itm := range src {
		tx := &_b.TransactionBuilder{}
		tx.Mutate(_b.SourceAccount{AddressOrSeed: itm.SrcAcc.Address})
		tx.Mutate(_b.Sequence{Sequence: itm.SrcAcc.GetNextSequence()})
		for i := 0; i < itm.LengthInfo; i++ {
			sell := itm.GetSellInfo(i)
			buy := itm.GetBuyInfo(i)
			order := _b.CreateOffer(
				_b.Rate{
					Selling: sell.GetAsset(),
					Buying:  buy.GetAsset(),
					Price:   itm.GetPrice(i),
				},
				itm.GetAmount(i))
			order.Mutate(_b.SourceAccount{AddressOrSeed: itm.SrcAcc.Address})
			tx.Mutate(order)
		}
		tx.Mutate(GetNetwork())
		if lenMemo > 0 && lenMemo <= 28 {
			tx.Mutate(_b.MemoText{Value: memo})
		}
		tx.TX.Fee = xdr.Uint32(100 * itm.LengthInfo)
		ret := tx.Sign(itm.SrcAcc.Secret)
		base64, _ := ret.Base64()
		retB64 = append(retB64, base64)
	}
	return
}
