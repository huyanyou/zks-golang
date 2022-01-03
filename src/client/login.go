package client

import "net/http"

type Htc struct {
	http.Client
}

var c *Htc

//	用户登陆
func Login(username string, password string) (resp *http.Response, err error) {
	req, err := http.NewRequest("Get", U["login"], nil)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, MainHeaders)
	res, err := c.Do(req)
	if err != nil {
		return res, err
	}
	return
}
