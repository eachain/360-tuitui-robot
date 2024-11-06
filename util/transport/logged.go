package transport

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unsafe"
)

// 记录每个http请求的请求内容和响应结果。
// 用法：设置http.Client.Transport为LoggedTransport即可。
// 如果LoggedTransport为nil，所有字段均采用默认值，即http.Client.Transport可以设置为(*LoggedTransport)(nil)。
type LoggedTransport struct {
	// 如果为nil，默认用http.DefaultTransport。
	http.RoundTripper

	// 如果为nil，默认用log.Printf。
	Logf func(string, ...any)
}

// 实现http.RoundTripper。
func (lt *LoggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := http.DefaultTransport
	if lt != nil && lt.RoundTripper != nil {
		rt = lt.RoundTripper
	}
	logf := log.Printf
	if lt != nil && lt.Logf != nil {
		logf = lt.Logf
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logf("logged transport: request %v: read request body: %v", req.URL.Path, err)
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	var reqStr string
	if len(reqBody) > 0 {
		if strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
			reqStr = fmt.Sprintf("%v bytes multipart/form-data", len(reqBody))
		} else {
			reqStr = unsafe.String(&reqBody[0], len(reqBody))
		}
	} else {
		reqStr = "(NoBody)"
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		logf("logged transport: request %v %v, roundtrip: %v",
			req.URL.Path, reqStr, err)
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	if err != nil {
		logf("logged transport: request %v %v, read response body: %v",
			req.URL.Path, reqStr, err)
		return resp, nil
	}

	var respStr string
	if len(respBody) > 0 {
		if strings.Contains(resp.Header.Get("Content-Type"), "multipart/form-data") {
			respStr = fmt.Sprintf("%v bytes multipart/form-data", len(respBody))
		} else {
			respStr = unsafe.String(&respBody[0], len(respBody))
		}
	} else {
		respStr = "(NoBody)"
	}

	logf("logged transport: request %v %v, response %v %v",
		req.URL.Path, reqStr, resp.Status, respStr)
	return resp, nil
}
