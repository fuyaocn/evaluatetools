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

// PostSend2Horizon Post方式发送数据到Horizon
type PostSend2Horizon struct {
	h    []*net.SocketHttp
	wait *sync.WaitGroup
}

// Send 发送签名到网络(以Post的方式发送)
func (ths *PostSend2Horizon) Send(addr string, t SendType, s ...string) (err error) {
	length := len(s)
	if length == 0 {
		err = fmt.Errorf("[PostSend2Horizon:Send] Can not send empty signature")
		_L.LoggerInstance.ErrorPrint("  **** %+v\r\n", err)
		return
	}
	ths.h = make([]*net.SocketHttp, length)
	for i := 0; i < length; i++ {
		ths.h[i] = new(net.SocketHttp)
		err = ths.h[i].Init(addr)
		if err != nil {
			ths.h[i].Result.Result = net.InitErr
			return
		}
	}
	ths.wait = new(sync.WaitGroup)
	// ths.wait.Init()
	switch t {
	case Linear:
		return ths.sendLinear(s...)
	case Multiple:
		return ths.sendMultiple(s...)
	default:
		err = fmt.Errorf("Unknown 'SendType' [%d]", t)
	}
	return
}

// GetSocket 获取Socket http
func (ths *PostSend2Horizon) GetSocket() []*net.SocketHttp {
	return ths.h
}

// sendLinear 线性方式发送签名结果到Horizon
func (ths *PostSend2Horizon) sendLinear(s ...string) error {
	_L.LoggerInstance.Debug(" >> Linear post send datas...")
	length := len(s)
	header := _AC.ConfigInstance.GetHorizonHeader()
	for i := 0; i < length; i++ {
		time.Sleep(time.Duration(50 * time.Millisecond))
		data := "tx=" + url.QueryEscape(s[i])
		err := ths.h[i].PostForm(data, header)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Post [Linear] data has error :\r\n * addr=%s\r\n * post data = %s\r\n * error = %+v\r\n",
				ths.h[i].Address, s[i], err)
			ths.h[i].Result.Result = net.Failure
			ths.h[i].Result.Extras.ResultXdr = err.Error()
			continue
		}

		err = ths.h[i].Response(ths.h[i].Result)
		if err == nil {
			if ths.h[i].Result.IsCorrect() {
				ths.h[i].Result.Result = net.Success
				continue
			}
		}
		ths.h[i].Result.Result = net.Failure
		if ths.h[i].Result.Extras != nil {
			ths.h[i].Result.DecodeError()
			_L.LoggerInstance.ErrorPrint(" ##[%d]## Post transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", i, err, ths.h[i].Result.Extras.ResultXdr)
		} else {
			ths.h[i].Result.Extras = &net.ExtrasData{
				ResultXdr: "Timeout",
			}
			_L.LoggerInstance.ErrorPrint(" ##[%d]## Post transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : Timeout\r\n", i, err)
		}
	}
	return nil
}

func (ths *PostSend2Horizon) sendMultiple(s ...string) error {
	_L.LoggerInstance.DebugPrint(" >> Multiple post send datas...\r\n")
	length := len(s)
	header := _AC.ConfigInstance.GetHorizonHeader()
	for i := 0; i < length; i++ {
		data := "tx=" + url.QueryEscape(s[i])
		err := ths.h[i].PostForm(data, header)
		if err != nil {
			_L.LoggerInstance.ErrorPrint("  **** Post [Linear] data has error :\r\n * addr=%s\r\n * post data = %s\r\n * error = %+v\r\n",
				ths.h[i].Address, s[i], err)
			ths.h[i].Result.Result = net.Failure
			continue
		}
		ths.wait.Add(1)
		go func(index int, http *net.SocketHttp, wg *sync.WaitGroup) {
			defer wg.Done()
			err, body := http.ResponseWithBody(http.Result)
			if err == nil {
				if http.Result.IsCorrect() {
					http.Result.Result = net.Success
				} else {
					http.Result.Result = net.Failure
					http.Result.Extras = &net.ExtrasData{
						ResultXdr: "Reject",
					}
					_L.LoggerInstance.ErrorPrint(" ##[%d]## Post transaction is reject!! ###\r\n ### body : %s\r\n ### Detail : Reject\r\n", index, body)
				}
				return
			}
			http.Result.Result = net.Failure
			if http.Result.Extras != nil {
				http.Result.DecodeError()
				_L.LoggerInstance.ErrorPrint(" ##[%d]## Post transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : %s\r\n", index, err, http.Result.Extras.ResultXdr)
			} else {
				http.Result.Extras = new(net.ExtrasData)
				http.Result.Extras.ResultXdr = "Timeout"
				_L.LoggerInstance.ErrorPrint(" ##[%d]## Post transaction is fail!! ###\r\n ### error : %v\r\n ### Detail : Timeout\r\n", index, err)
			}
		}(i, ths.h[i], ths.wait)
	}
	_L.LoggerInstance.DebugPrint(" >> Multiple post send datas end amd waiting response ...\r\n")
	ths.wait.Wait()
	return nil
}
