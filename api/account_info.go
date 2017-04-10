package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_AC "jojopoper/NBi/StressTest/appconf"
	_L "jojopoper/NBi/StressTest/log"
	"net/http"
	"strconv"
	"sync"
)

// AssetInfo asset info
type AssetInfo struct {
	Type    string `json:"asset_type"`
	Balance string `json:"balance"`
	Code    string `json:"asset_code"`
	Issuer  string `json:"asset_issuer"`
}

// AccountInfo account info
type AccountInfo struct {
	ID       string      `json:"id"`
	Sequence string      `json:"sequence"`
	Assets   []AssetInfo `json:"balances"`
	Address  string      `json:"address"`
	Status   int         `json:"status"`
	Balance  float64
	sequence uint64
	Secret   string
}

// Init set address
func (ths *AccountInfo) Init(id, secret string) {
	ths.Address = id
	ths.Secret = secret
}

// GetInfo get base information
func (ths *AccountInfo) GetInfo(wt *sync.WaitGroup) error {
	if wt != nil {
		defer wt.Done()
	}
	addr := fmt.Sprintf("%s/accounts/%s", _AC.ConfigInstance.GetHorizonServer(), ths.Address)

	resp, err := http.Get(addr)
	if err != nil {
		return fmt.Errorf("\r\n  **** Http get '%s' has error : \r\n %+v", ths.Address, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("\r\n **** HTTP Response ERROR\r\n\tReadAll: %+v", err)
	}

	err = json.Unmarshal(body, ths)
	if err != nil {
		return fmt.Errorf("\r\n  **** Unmarshal body has error : %+v\r\n[%s]", err, string(body))
	}
	if ths.Status == 0 {
		for _, itm := range ths.Assets {
			if itm.Type == "native" {
				ths.Balance, _ = strconv.ParseFloat(itm.Balance, 64)
				break
			}
		}

		ths.sequence, err = strconv.ParseUint(ths.Sequence, 10, 64)
		_L.LoggerInstance.Debug(" Current Account info : %+v\r\n", ths)
		return err
	}
	return fmt.Errorf("Account '%s' is not exist\r\n[%s]", ths.Address, string(body))
}

// GetNextSequence get next sequence
func (ths *AccountInfo) GetNextSequence() uint64 {
	ths.sequence++
	return ths.sequence
}

// GetCurrentSequence get currnt sequence
func (ths *AccountInfo) GetCurrentSequence() uint64 {
	return ths.sequence
}

// ResetSequence reset currnt sequence
func (ths *AccountInfo) ResetSequence() {
	ths.sequence--
}
