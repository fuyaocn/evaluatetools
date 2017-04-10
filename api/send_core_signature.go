package api

import (
	"fmt"
	_AC "jojopoper/NBi/StressTest/appconf"
	_L "jojopoper/NBi/StressTest/log"
	"jojopoper/NBi/StressTest/net"
	"net/url"
	"sync"
	"time"
)

// SendToCore 直接发送签名到core
type SendToCore struct {
	h    []*net.SocketHttp
	wait *sync.WaitGroup
}

// SendCore 发送签名到core网络(以get的方式发送)
func (ths *SendToCore) SendCore(addr string, t SendType, s ...string) (err error) {
	length := len(s)
	if length == 0 {
		err = fmt.Errorf("[SendToCore:SendCore] Can not send empty signature")
		_L.LoggerInstance.ErrorPrint("  **** %+v\r\n", err)
		return
	}
	ths.h = make([]*net.SocketHttp, length)
	for i := 0; i < length; i++ {
		ths.h[i] = new(net.SocketHttp)
		address := addr + "/tx?blob=" + url.QueryEscape(s[i])
		err = ths.h[i].Init(address)
		if err != nil {
			ths.h[i].Result.Result = net.InitErr
			return
		}
	}
	ths.wait = new(sync.WaitGroup)
	// ths.wait.Init()
	switch t {
	case Linear:
		return ths.sendCoreLinear(s...)
	case Multiple:
		return ths.sendCoreMultiple(s...)
	default:
		err = fmt.Errorf("Unknown 'SendType' [%d]", t)
	}
	return
}

// GetCoreSocket 获取Socket http
func (ths *SendToCore) GetCoreSocket() []*net.SocketHttp {
	return ths.h
}

// sendCoreLinear 线性方式发送签名结果到core
func (ths *SendToCore) sendCoreLinear(s ...string) error {
	_L.LoggerInstance.DebugPrint(" >> Linear get send datas to core...")
	length := len(s)
	body := ""
	header := _AC.ConfigInstance.GetHorizonHeader()
	for i := 0; i < length; i++ {
		time.Sleep(time.Duration(50 * time.Millisecond))
		err := ths.h[i].Get(header)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Get [Linear] data has error :\r\n * addr=%s\r\n * data = %s\r\n * error = %+v\r\n",
				ths.h[i].Address, s[i], err)
			ths.h[i].CoreResult.Result = net.Failure
			ths.h[i].CoreResult.Error = err.Error()
			continue
		}

		err, body = ths.h[i].ResponseWithBody(ths.h[i].CoreResult)
		if err == nil {
			if ths.h[i].CoreResult.IsCorrect() {
				ths.h[i].CoreResult.Result = net.Success
			} else {
				ths.h[i].CoreResult.Result = net.Failure
				ths.h[i].CoreResult.Error = "Reject"
				_L.LoggerInstance.ErrorPrint(" ##[%d]## Send to core transaction is reject!! ###\r\n ### body : %s\r\n ### Detail : Reject\r\n", i, body)
			}
			continue
		}
		ths.h[i].CoreResult.Result = net.Failure
		if !ths.h[i].CoreResult.IsCorrect() {
			ths.h[i].CoreResult.DecodeError()
		} else {
			ths.h[i].CoreResult.Error = "Timeout"
		}
		_L.LoggerInstance.ErrorPrint(" ##[%d]## Send to core transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", i, err, ths.h[i].CoreResult.Error)
	}
	return nil
}

// sendCoreMultiple 以并行方式发送签名到core
func (ths *SendToCore) sendCoreMultiple(s ...string) error {
	_L.LoggerInstance.DebugPrint(" >> Multiple get send datas to core...")
	length := len(s)
	header := _AC.ConfigInstance.GetHorizonHeader()
	for i := 0; i < length; i++ {
		err := ths.h[i].Get(header)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Get [Linear] data has error :\r\n * addr=%s\r\n * post data = %s\r\n * error = %+v\r\n",
				ths.h[i].Address, s[i], err)
			ths.h[i].CoreResult.Result = net.Failure
			continue
		}
		ths.wait.Add(1)
		go func(index int, http *net.SocketHttp, wg *sync.WaitGroup) {
			defer wg.Done()
			err, body := http.ResponseWithBody(http.CoreResult)
			if err == nil {
				if http.CoreResult.IsCorrect() {
					http.CoreResult.Result = net.Success
				} else {
					http.CoreResult.Result = net.Failure
					http.CoreResult.Error = "Reject"
					_L.LoggerInstance.ErrorPrint(" ##[%d]## Send to core transaction is reject!! ###\r\n ### body : %s\r\n ### Detail : Reject\r\n", index, body)
				}
				return
			}
			http.CoreResult.Result = net.Failure
			if !http.CoreResult.IsCorrect() {
				http.CoreResult.DecodeError()
			} else {
				http.CoreResult.Error = "Timeout"
			}
			_L.LoggerInstance.ErrorPrint(" ##[%d]## Send to core transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", index, err, http.CoreResult.Error)
		}(i, ths.h[i], ths.wait)
	}
	ths.wait.Wait()
	return nil
}
