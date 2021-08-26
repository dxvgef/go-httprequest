# go-httprequest

## 导入路径
> github.com/dxvgef/go-httprequest

## 简介
Go语言的HTTP请求包，功能特性：
- 线程安全
- 链式语法
- 可按间隔时间和最大次数自动重发请求
- 可定义多个URL轮循请求，当请求失败时自动使用下一个URL重发请求
- 对响应数据进行类型转换

## 基本示例
```go
package main

import (
	"log"
	
	"github.com/dxvgef/go-httprequest"
)

// 根据配置创建实例
func main() {
    request := httprequest.New(httprequest.Config{
        Timeout:       10,  // 每个URL请求的超时时间(秒)
        // 需要重试的响应状态码
        RetryStatus:   []int{500, 502, 503, 504, 506, 507},
        RetryCount:    1,   // 每个URL的重试次数
        RetryInterval: 3000,    // 每个URL的重试间隔时间
    })
    
    // 定义请求参数
    resp := request.AddHeader("test", "ok").    // 定义Header参数
    AddURL("http://localhost/").    // 添加主URL
    AddURL("http://loalhost/backup").  // 添加备URL
    GET()   // 执行GET请求
    // 判断请求过程是否出错
    if resp.Error() != nil {
        log.Println(err)
        return
    }
    // 响应状态码
    log.Println(resp.StatusCode())
    
    // 获取*http.Response
    log.Println(resp.Raw())
    
    // 将响应转为[]byte
    bytes, err := resp.Bytes()
    
    // 将响应数据转为string
    str, err := resp.String()
    
    // 将响应数据按查询参数解析到net.Values
    values, err := resp.ParseQuery()
    
    // 将响应数据按JSON解析
    err := resp.UnmarshalJSON(&target)
    
    // 将响应数据按XML解析
    err := resp.UnmarshalXML(&target)	
}
```
