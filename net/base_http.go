package net

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// BaseHttp go基础http的调用
type BaseHttp struct {
}

// GetHTTPClient 得到一个http实例
func (ths *BaseHttp) GetHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:  time.Second * 120,
		Deadline: time.Now().Add(90 * time.Second),
		// KeepAlive: time.Second * 30,
	}
	trans := &http.Transport{
		Dial: dialer.Dial,
		ResponseHeaderTimeout: 60 * time.Second,

		// DialTLS:             dialer.Dial,
		// TLSHandshakeTimeout: 20 * time.Second,
		// TLSHandshakeTimeout: 10 * time.Second,
	}

	ret := &http.Client{
		Transport: trans,
		Timeout:   180 * time.Second,
	}
	return ret
}

// HttpGet http get function
func (ths *BaseHttp) HttpGet(address string, ret interface{}) error {
	phttp := ths.GetHTTPClient()
	resp, err := phttp.Get(address)
	if err != nil {
		return fmt.Errorf("  **** BaseHttp 'get' has error : \r\n %v", err)
	}
	defer resp.Body.Close()

	return ths.getResponseDecode(resp, ret)
}

// HttpPostForm http post form
func (ths *BaseHttp) HttpPostForm(address, data string, ret interface{}) error {

	resp, err := http.Post(address,
		"application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("  **** BaseHttp 'post' has error : \r\n %v", err)
	}

	return ths.getResponseDecode(resp, ret)
}

func (ths *BaseHttp) getResponseDecode(resp *http.Response, ret interface{}) (err error) {
	if resp == nil {
		return fmt.Errorf("http.Response is nil")
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(ret)
}
