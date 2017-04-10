package statics

import (
	"jojopoper/NBi/StressTest/net"
	"time"
)

// BaseDefine 统计基础定义
type BaseDefine struct {
	StaticName   string
	GroupIndex   int
	ItemIndex    int
	StartTime    int64
	CompleteTime int64
	Success      string
	FailureCause string
	Action       string
}

// SetBaseInfo 设置基础信息内容
// grpIndex 组索引号；itmIndex 项目索引号；name 测试用例名称，最大程度255个字符（可以为空）
func (ths *BaseDefine) SetBaseInfo(grpIndex, itmIndex int, name string) {
	ths.GroupIndex = grpIndex
	ths.ItemIndex = itmIndex
	ths.StaticName = name
}

// SetTimeValue 设置时间的值
func (ths *BaseDefine) SetTimeValue(tt TimeType, t int64) {
	if tt == StartTimeFlag {
		ths.StartTime = t
	} else if tt == CompleteTimeFlag {
		ths.CompleteTime = t
	}
}

// SetTime 设置时间为当前时间UnixNano
func (ths *BaseDefine) SetTime(tt TimeType) {
	ths.SetTimeValue(tt, time.Now().UnixNano())
}

// SetResult 设置结果
func (ths *BaseDefine) SetResult(rslt net.SocketResult, cause string) {
	ths.Success = string(rslt)
	ths.FailureCause = cause
}
