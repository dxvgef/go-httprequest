package httprequest

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
)

const ErrNoData = "no data"

type Response struct {
	request *Request
	resp    *http.Response
	body    []byte
}

// 错误信息
func (response *Response) Error() error {
	return response.request.err
}

// 原生对象
func (response *Response) Raw() *http.Response {
	if response.request.err != nil {
		return nil
	}
	return response.resp
}

// 状态码
func (response *Response) StatusCode() int {
	if response.request.err != nil {
		return 0
	}
	return response.resp.StatusCode
}

// 将响应数据转为[]byte
func (response *Response) Bytes() (data []byte, err error) {
	if response.request.err != nil {
		err = response.request.err
		return
	}
	return response.body, nil
}

// 将响应数据转为字符串
func (response *Response) String() (str string, err error) {
	if response.request.err != nil {
		err = response.request.err
		return
	}
	return BytesToStr(response.body), nil
}

// 将响应数据做为JSON解码
// nolint:govet
func (response *Response) UnmarshalJSON(obj interface{}) (err error) {
	if response.request.err != nil {
		return response.request.err
	}
	if len(response.body) == 0 {
		return errors.New(ErrNoData)
	}
	err = json.Unmarshal(response.body, obj)
	return
}

// 将响应数据做为XML解码
// nolint:govet
func (response *Response) UnmarshalXML(obj interface{}) (err error) {
	if response.request.err != nil {
		return response.request.err
	}
	if len(response.body) == 0 {
		return errors.New(ErrNoData)
	}
	err = xml.Unmarshal(response.body, obj)
	return
}

// 将响应数据做为查询参数解码
func (response *Response) ParseQuery() (values url.Values, err error) {
	if response.request.err != nil {
		err = response.request.err
		return
	}
	if len(response.body) == 0 {
		err = errors.New(ErrNoData)
		return
	}
	var data string
	data, err = response.String()
	if err != nil {
		return
	}
	return url.ParseQuery(data)
}
