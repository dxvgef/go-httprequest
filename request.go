package httprequest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 请求结构体
type Request struct {
	err      error
	config   *Config
	endpoint []string
	header   map[string]string
	values   url.Values
	body     []byte
}

// 添加端点URL
func (request *Request) AddEndpoint(url string) *Request {
	request.endpoint = append(request.endpoint, url)
	return request
}

// 设置端点URL
func (request *Request) SetEndpoint(urls []string) *Request {
	request.endpoint = urls
	return request
}

// 设置Header参数值
func (request *Request) SetHeader(headers map[string]string) *Request {
	request.header = headers
	return request
}

// 添加Header参数值
func (request *Request) AddHeader(key, value string) *Request {
	if request.header == nil {
		request.header = make(map[string]string)
	}
	request.header[key] = value
	return request
}

// 设置请求参数值
func (request *Request) AddValue(key, value string) *Request {
	request.body = nil
	if request.values == nil {
		request.values = url.Values{}
	}
	request.values.Set(key, value)
	return request
}

// 设置请求参数值
func (request *Request) SetValue(values url.Values) *Request {
	request.body = nil
	request.values = values
	return request
}

// 设置请求数据
func (request *Request) SetBody(body []byte) *Request {
	request.values = nil
	request.body = body
	return request
}

// 设置请求JSON数据
func (request *Request) SetJSON(data interface{}) *Request {
	buf, err := json.Marshal(data)
	if err != nil {
		request.err = err
	} else {
		request.body = buf
	}
	return request
}

// 设置请求XML数据
func (request *Request) SetXML(data interface{}) *Request {
	buf, err := xml.Marshal(data)
	if err != nil {
		request.err = err
	} else {
		request.body = buf
	}
	return request
}

// 获取error
func (request *Request) Error() error {
	return request.err
}

// 发起GET请求
func (request *Request) GET() *Response {
	return request.do(http.MethodGet)
}

// 发起POST请求
func (request *Request) POST() *Response {
	return request.do(http.MethodPost)
}

// 发起PUT请求
func (request *Request) PUT() *Response {
	return request.do(http.MethodPut)
}

// 发起PATCH请求
func (request *Request) PATCH() *Response {
	return request.do(http.MethodPatch)
}

// 发起HEAD请求
func (request *Request) HEAD() *Response {
	return request.do(http.MethodHead)
}

// 发起DELETE请求
func (request *Request) DELETE() *Response {
	return request.do(http.MethodDelete)
}

// 发起OPTIONS请求
func (request *Request) OPTIONS() *Response {
	return request.do(http.MethodOptions)
}

// 发起TRACE请求
func (request *Request) TRACE() *Response {
	return request.do(http.MethodTrace)
}

// 发起请求
func (request *Request) do(method string) *Response {
	var (
		body           io.Reader
		endpointIndex  int                         // 当前URL索引
		endpointLength = len(request.endpoint) - 1 // URL最大索引
		retry          uint8                       // 当前URL重试次数
		req            *http.Request
		resp           *http.Response
		response       Response
	)
	if len(request.endpoint) == 0 {
		request.err = errors.New("没有定义目标端点URL")
		response.request = request
		return &response
	}
	for {
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodTrace:
			if len(request.values) > 0 {
				body = strings.NewReader(request.values.Encode())
			}
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			if len(request.values) > 0 {
				body = strings.NewReader(request.values.Encode())
			} else if len(request.body) > 0 {
				body = bytes.NewBuffer(request.body)
			}
		default:
			request.err = errors.New("不支持 " + method + " 方法的请求")
			response.request = request
			return &response
		}
		client := &http.Client{
			Timeout: time.Duration(request.config.Timeout) * time.Second,
		}
		// nolint:noctx
		req, request.err = http.NewRequest(method, request.endpoint[endpointIndex], body)
		if request.err != nil {
			if endpointIndex < endpointLength {
				endpointIndex++
				client.CloseIdleConnections()
				continue
			}
			client.CloseIdleConnections()
			break
		}
		for k, v := range request.header {
			req.Header.Set(k, v)
		}
		resp, request.err = client.Do(req) // nolint:bodyclose
		if request.err != nil || resp == nil || inIntSlice(resp.StatusCode, request.config.RetryStatus) {
			if retry < request.config.RetryCount {
				retry++
				client.CloseIdleConnections()
				time.Sleep(time.Duration(request.config.RetryInterval) * time.Millisecond)
				continue
			}
			if endpointIndex < endpointLength {
				endpointIndex++
				client.CloseIdleConnections()
				continue
			}
		}
		client.CloseIdleConnections()
		break
	}
	response.request = request
	response.resp = resp
	if request.err == nil && resp.Body != nil {
		response.body, request.err = ioutil.ReadAll(resp.Body)
		request.err = resp.Body.Close()
	}
	return &response
}
