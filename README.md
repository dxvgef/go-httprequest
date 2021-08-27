# go-httprequest

## 简介
Go语言的HTTP请求包，功能特性：
- 线程安全
- 链式语法
- 支持`GET`,`POST`,`PUT`,`PATCH`,`HEAD`,`DELETE`,`OPTIONS`,`TRACE`方法
- 可按间隔时间和最大次数自动重发请求
- 可定义多个URL轮循请求，当请求失败时自动使用下一个URL重发请求
- 对响应数据进行类型转换

## 导入路径
> github.com/dxvgef/go-httprequest

## 当前版本
> v0.0.3

## 基本示例
```go
package main

import (
	"log"
	
	"github.com/dxvgef/go-httprequest"
)

func main() {
    // 使用默认配置创建请求实例，可传入httprequest.Config{}来自定义配置
    request := httprequest.New()
    
    resp := request.AddHeader("test", "ok").    // 添加Header参数
    AddEndpoint("http://localhost/").    // 添加端点
    GET()   // 执行GET请求
    // 判断请求过程是否出错
    if resp.Error() != nil {
        log.Println(err)
        return
    }

    // 响应状态码
    log.Println(resp.StatusCode())

    // 将响应数据转为string
    str, err := resp.String()
}
```

## 参数配置
如果执行`httprequest.New()`时不传入参数，则自动使用默认配置。

默认配置及具体参数说明如下：

```go
var DefaultConfig = Config{
    // 每次请求的超时时间(秒)
    Timeout:       10,
    // 触发重试的响应状态码
    RetryStatus:   []int{500, 502, 503, 504, 506, 507},
    // 每个端点的重试次数
    RetryCount:    1,
    // 每个端点的重试间隔时间(毫秒)
    RetryInterval: 1000,
}
```

## 响应数据处理

- `Error()` 判断请求过程是否存在错误
- `Raw()` 获取`*http.Response`
- `Bytes() ([]byte, error)` 转为`[]byte`类型
- `ParseQuery() (url.Values, error)` 将响应数据按查询参数格式解析到`net.Values`类型的变量
- `UnmarshalJSON(interface{}) error` 将响应数据按JSON编码反序列化到入参变量
- `UnmarshalXML(interface{}) error` 将响应数据按XML编码反序列化到入参变量    
