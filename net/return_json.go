package net

import (
	"github.com/stellar/go/xdr"
)

// ExtrasData extras data
type ExtrasData struct {
	ResultXdr string `json:"result_xdr"`
}

// JSONResult post or get result define
type JSONResult struct {
	Result SocketResult
	Hash   string      `json:"hash"`
	Extras *ExtrasData `json:"extras"`
}

// IsCorrect 是否存在错误
func (ths *JSONResult) IsCorrect() bool {
	return ths.Extras == nil && ths.Hash != ""
}

// DecodeError 解码错误信息
func (ths *JSONResult) DecodeError() {
	if ths.Extras != nil {
		tret := &xdr.TransactionResult{}
		tret.Scan(ths.Extras.ResultXdr)
		ths.Extras.ResultXdr = tret.Result.Code.String()
	}
}

// ===============================================================

// JSONCoreResult core返回结果解析
type JSONCoreResult struct {
	Result SocketResult
	Status string `json:"status"`
	Error  string `json:"error"`
}

// IsCorrect 是否存在错误
func (ths *JSONCoreResult) IsCorrect() bool {
	return len(ths.Error) == 0
}

// DecodeError 解码错误信息
func (ths *JSONCoreResult) DecodeError() {
	if !ths.IsCorrect() {
		tret := &xdr.TransactionResult{}
		tret.Scan(ths.Error)
		ths.Error = tret.Result.Code.String()
	}
}
