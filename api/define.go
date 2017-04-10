package api

// SendType 以什么方式发送
type SendType int

const (
	// Linear 发送一个立刻接收，之后再发送另外一个
	Linear SendType = 1
	// Multiple 先发送所有，然后再读取结果
	Multiple SendType = 2
)
