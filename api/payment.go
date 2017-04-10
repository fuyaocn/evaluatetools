package api

import (
	_b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// Payment 发送资产或币
type Payment struct {
	SendToCore
	PostSend2Horizon
}

// GetSignature 得到签名
func (ths *Payment) GetSignature(issuer, code, amount, memo string, src []*AccountInfo, dest []*AccountInfo) (retB64 []string) {
	if len(issuer) == 0 || len(src) == 0 || len(dest) == 0 {
		return
	}
	retB64 = make([]string, 0)
	isNative := len(issuer) == 0
	lengthDest := len(dest)
	lenMemo := len(memo)

	for _, itm := range src {
		tx := &_b.TransactionBuilder{}
		tx.Mutate(_b.SourceAccount{AddressOrSeed: itm.Address})
		for _, d := range dest {
			pay := _b.Payment(
				_b.Destination{AddressOrSeed: d.Address},
				_b.SourceAccount{AddressOrSeed: itm.Address},
			)
			if isNative {
				pay.Mutate(
					_b.NativeAmount{Amount: amount},
				)
			} else {
				pay.Mutate(
					_b.CreditAmount{
						Code:   code,
						Issuer: issuer,
						Amount: amount,
					},
				)
			}
			tx.Mutate(pay)
		}
		tx.Mutate(_b.Sequence{Sequence: itm.GetNextSequence()})
		tx.Mutate(GetNetwork())
		if lenMemo > 0 && lenMemo <= 28 {
			tx.Mutate(_b.MemoText{Value: memo})
		}
		tx.TX.Fee = xdr.Uint32(100 * lengthDest)
		ret := tx.Sign(itm.Secret)
		base64, _ := ret.Base64()
		retB64 = append(retB64, base64)
	}
	return
}
