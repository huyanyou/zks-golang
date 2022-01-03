package client

const baseUrl = "http://xk.henu.edu.cn"

//	henu Server 接口
var U map[string]string = map[string]string{
	// 登陆接口
	"login": baseUrl + "/cas/logon.action",
}

// request Headers
var UA string = "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
var MainHeaders = map[string]string{
	"User-Agent":      UA,
	"Accept":          "*/*",
	"Accept-Encoding": "gzip, deflate",
	"Connection":      "keep-alive",
}
