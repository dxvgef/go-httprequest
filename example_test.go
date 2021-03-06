package httprequest

import (
	"net/http"
	"testing"
)

// 测试多个URL轮循
func TestURL(t *testing.T) {
	request := New(Config{
		Timeout:       10,
		RetryStatus:   []int{500, 502, 503, 504, 506, 507},
		RetryCount:    1,
		RetryInterval: 3000,
	})
	resp := request.AddHeader("test", "ok").
		AddEndpoint("http://127.0.0.1/").
		AddEndpoint("http://127.0.0.1/backup").
		GET()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	t.Log(resp.StatusCode())
}

// 测试添加value
func TestValues(t *testing.T) {
	request := New(Config{
		Timeout:       10,
		RetryStatus:   []int{500, 502, 503, 504, 506, 507},
		RetryCount:    1,
		RetryInterval: 3000,
	})
	resp := request.AddHeader("test", "ok").
		AddEndpoint("http://127.0.0.1/backup").
		AddValue("test", "ok").
		PUT()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	t.Log(resp.StatusCode())
}

// 测试结果转String
func TestString(t *testing.T) {
	request := New(Config{
		Timeout:       10,
		RetryStatus:   []int{500, 502, 503, 504, 506, 507},
		RetryCount:    1,
		RetryInterval: 3000,
	})
	resp := request.AddHeader("test", "ok").
		AddEndpoint("http://127.0.0.1/string").
		GET()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	t.Log(resp.StatusCode())
	t.Log(resp.String())
}

// 测试查询参数结果
func TestQuery(t *testing.T) {
	request := New(Config{
		Timeout:       10,
		RetryStatus:   []int{500, 502, 503, 504, 506, 507},
		RetryCount:    1,
		RetryInterval: 3000,
	})
	resp := request.AddHeader("test", "ok").
		AddEndpoint("http://127.0.0.1/query").
		GET()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	t.Log(resp.StatusCode())
	values, err := resp.ParseQuery()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(values.Encode())
}

// 测试JSON结果
func TestJSON(t *testing.T) {
	request := New(Config{
		Timeout:       10,
		RetryStatus:   []int{500, 502, 503, 504, 506, 507},
		RetryCount:    1,
		RetryInterval: 3000,
	})
	resp := request.AddHeader("test", "ok").
		AddEndpoint("http://127.0.0.1/json").
		GET()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	t.Log(resp.StatusCode())
	result := make(map[string]string)
	err := resp.UnmarshalJSON(&result)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

// 测试重试
func TestRetry(t *testing.T) {
	requestConfig := DefaultConfig
	requestConfig.RetryCount = 3
	requestConfig.RetryInterval = 3000
	request := New(requestConfig)
	resp := request.AddEndpoint("http://127.0.0.1:8080/").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddHeader("Authorization", "123456").
		AddValue("test", "ok").
		PUT()
	if resp.Error() != nil {
		t.Error(resp.Error())
		return
	}
	if resp.StatusCode() != http.StatusNoContent && resp.StatusCode() != http.StatusOK {
		t.Error("错误状态码", resp.StatusCode())
		return
	}
}
