package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestDoubleNumHandler 测试对应文件中的doubleNumHandler函数
// t是单元测试。
// 本测试中，我们只测试了单个输入
func TestDoubleNumHandler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "localhost:9000?v=3", nil)
	if err != nil {
		// Fatalf is equivalent to Logf followed by FailNow.
		t.Fatalf("cannot create a new request: %v, error: %v", request, err)
	}

	// NewRecorder用于大规模测试，recorder内部嵌套了response
	rec := httptest.NewRecorder()
	// TODO：直接调用要测试的函数
	doubleNumHandler(rec, request)

	// 测试响应是否正确
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		// Errorf is equivalent to Logf followed by Fail. 因此手动return
		t.Errorf("received status code %d, expect %d\n", res.StatusCode, http.StatusOK)
		return
	}

	// 测试结果是否正确
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read all from response body, err: %v\n", err)
	}
	result, err := strconv.Atoi(strings.TrimSpace(string(resBytes)))
	if err != nil {
		t.Errorf("cannot convert response body to int, err: %v", err)
		return
	}
	if result != 6 {
		t.Errorf("received result %d, expect 6", result)
		return
	}

}

// 一举测试多组输入
func TestDoubleNumHandler2(t *testing.T) {
	// 编写单个测试的结构体，应该包含输入数据、输出数据等fields
	testCases := []struct {
		name   string
		input  string
		result int
		status int
		err    string
	}{
		{name: "double of 2", input: "2", result: 4, status: http.StatusOK, err: ""},
		{name: "double of 3", input: "3", result: 6, status: http.StatusOK, err: ""},
		{name: "double of -1", input: "-1", result: -2, status: http.StatusOK, err: ""},
		{name: "double of 100", input: "100", result: 200, status: http.StatusOK, err: ""},
		{name: "double of 123", input: "123", result: 246, status: http.StatusOK, err: ""},
		{name: "double of -123", input: "-123", result: -246, status: http.StatusOK, err: "missing value v"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "localhost:9000?v="+testCase.input, nil)
			if err != nil {
				// Fatalf is equivalent to Logf followed by FailNow.
				t.Fatalf("cannot create a new request: %v, error: %v", request, err)
			}

			rec := httptest.NewRecorder()
			doubleNumHandler(rec, request)

			res := rec.Result()
			if res.StatusCode != testCase.status {
				// Errorf is equivalent to Logf followed by Fail. 因此手动return
				t.Errorf("received status code %d, expect %d\n", res.StatusCode, testCase.status)
				return
			}

			resBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("cannot read all from response body, err: %v\n", err)
			}
			defer func() {
				if err := res.Body.Close(); err != nil {
					t.Fatalf("cannot close response body")
				}
			}()

			trimedResult := strings.TrimSpace(string(resBytes))
			if res.StatusCode != http.StatusOK {
				// 用户输入错误，要判定函数的error handling是否是我们想要的结果（对比err字符串）
				if trimedResult != testCase.err {
					t.Errorf("received error msg %s, expect %s\n", trimedResult, testCase.err)
				}
				return
			}

			doubleVal, err := strconv.Atoi(trimedResult)
			if err != nil {
				t.Errorf("cannot convert response body to int, err: %v", err)
				return
			}
			if doubleVal != testCase.result {
				t.Errorf("received result %d, expect 6", doubleVal)
				return
			}
		})
	}
}
