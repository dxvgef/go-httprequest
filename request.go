package httprequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	err    error
	config *Config
	url    []string
	header map[string]string
	values url.Values
	body   []byte
	req    *http.Request
}

// 添加URL
func (request *Request) AddURL(url string) *Request {
	request.url = append(request.url, url)
	return request
}

// 设置URL
func (request *Request) SetURL(urls []string) *Request {
	request.url = urls
	return request
}

// 设置Header参数值
func (request *Request) SetHeaders(headers map[string]string) *Request {
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
		body        io.Reader
		urlIndex    int                    // 当前URL索引
		urlIndexMax = len(request.url) - 1 // URL最大索引
		retry       uint8                  // 当前URL重试次数
		resp        *http.Response
		response    Response
	)

	if len(request.url) == 0 {
		request.err = errors.New("未定义请求URL")
		response.request = request
		return &response
	}

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
		request.err = errors.New("不支持的请求方法：" + method)
		response.request = request
		return &response
	}

	for {
		client := &http.Client{
			Timeout: time.Duration(request.config.Timeout) * time.Second,
		}
		request.req, request.err = http.NewRequest(method, request.url[urlIndex], body)
		if request.err != nil {
			if retry < request.config.RetryCount {
				retry++
				time.Sleep(time.Duration(request.config.RetryInterval) * time.Millisecond)
				continue
			}
			if urlIndex < urlIndexMax {
				urlIndex++
				continue
			}
		}
		for k, v := range request.header {
			request.req.Header.Set(k, v)
		}
		resp, request.err = client.Do(request.req)
		if request.err != nil {
			if retry < request.config.RetryCount {
				retry++
				time.Sleep(time.Duration(request.config.RetryInterval) * time.Millisecond)
				continue
			}
			if urlIndex < urlIndexMax {
				urlIndex++
				continue
			}
		}
		if inIntSlice(resp.StatusCode, request.config.RetryStatus) {
			if retry < request.config.RetryCount {
				retry++
				time.Sleep(time.Duration(request.config.RetryInterval) * time.Millisecond)
				continue
			}
			if urlIndex < urlIndexMax {
				urlIndex++
				continue
			}
		}
		break
	}
	response.request = request
	response.resp = resp
	if request.err == nil && resp.Body != nil {
		response.body, request.err = ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close() // nolint:errcheck
	}
	return &response
}
