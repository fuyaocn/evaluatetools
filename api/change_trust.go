package api

import (
	_b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// ChangeTrust 修改信任
type ChangeTrust struct {
	PostSend2Horizon
	SendToCore
}

// GetSignature 得到签名
func (ths *ChangeTrust) GetSignature(issuer, code, memo string, srcAcc ...*AccountInfo) (retB64 []string) {
	if len(issuer) == 0 || len(srcAcc) == 0 {
		return
	}
	retB64 = make([]string, 0)
	lenMemo := len(memo)

	for _, itm := range srcAcc {
		tx := &_b.TransactionBuilder{}
		tx.Mutate(_b.SourceAccount{AddressOrSeed: itm.Address})
		tx.Mutate(_b.ChangeTrust(
			_b.Asset{
				Code:   code,
				Issuer: issuer,
				Native: false,
			},
			_b.MaxLimit,
			_b.SourceAccount{AddressOrSeed: itm.Address},
		))
		tx.Mutate(_b.Sequence{Sequence: itm.GetNextSequence()})
		if lenMemo > 0 && lenMemo <= 28 {
			tx.Mutate(_b.MemoText{Value: memo})
		}
		tx.Mutate(GetNetwork())
		tx.TX.Fee = xdr.Uint32(100)
		ret := tx.Sign(itm.Secret)
		base64, _ := ret.Base64()
		retB64 = append(retB64, base64)
	}
	return
}
