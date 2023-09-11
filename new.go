package httprequest

import (
	"net"
	"net/http"
	"time"
)

// 配置
type Config struct {
	HTTPProxy     *net.Addr // HTTP代理地址
	Timeout       int       // 每次请求的超时时间(秒)
	RetryStatus   []int     // 触发重试的状态码
	RetryCount    uint8     // 请求重试次数，如果URL数量大于1，建议败为0，则使用下一个URL重试，而不是对同一个URL重试
	RetryInterval uint16    // 每次重试的间隔时间(毫秒)
}

// 默认配置
var DefaultConfig = Config{
	Timeout:       10,
	RetryStatus:   []int{500, 502, 503, 504, 506, 507},
	RetryCount:    1,
	RetryInterval: 1000,
}

// 新建请求实例
func New(config ...Config) *Request {
	var request Request
	if len(config) == 0 {
		// 使用默认配置
		request.config = &DefaultConfig
	} else {
		// 使用指定配置
		request.config = &config[0]
	}
	// 设置每次请求的超时时间
	request.client = &http.Client{
		Timeout: time.Duration(request.config.Timeout) * time.Second,
	}
	return &request
}
