package client

import "net/http"

//	设置请求头
func SetHeaders(req *http.Request, headers map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}
