package api

import (
	_b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// CreateAccount 创建账户
type CreateAccount struct {
	PostSend2Horizon
}

// GetSignature 得到签名
func (ths *CreateAccount) GetSignature(srcAcc *AccountInfo, startBalance, memo string, accIDs ...string) string {
	cnt := len(accIDs)
	if cnt == 0 {
		return ""
	}
	tx := &_b.TransactionBuilder{}
	tx.Mutate(_b.SourceAccount{AddressOrSeed: srcAcc.Address})

	for _, itm := range accIDs {
		tx.Mutate(_b.CreateAccount(
			_b.Destination{AddressOrSeed: itm},
			_b.NativeAmount{Amount: startBalance},
			_b.SourceAccount{AddressOrSeed: srcAcc.Address},
		))
	}
	tx.Mutate(_b.Sequence{Sequence: srcAcc.GetNextSequence()})
	tx.Mutate(GetNetwork())
	lenMemo := len(memo)
	if lenMemo > 0 && lenMemo <= 28 {
		tx.Mutate(_b.MemoText{Value: memo})
	}

	tx.TX.Fee = xdr.Uint32(100 * cnt)
	ret := tx.Sign(srcAcc.Secret)
	base64, _ := ret.Base64()
	return base64
}
