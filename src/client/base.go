package client

const (
	clientUrl = "https://xk.henu.edu.cn/"
	serverUrl = "https://xk.henu.edu.cn:80/"
)

//	henu Server 接口
var Urls map[string]string = map[string]string{
	// 登陆接口
	"login":          clientUrl + "cas/login.action",
	"logon":          clientUrl + "cas/logon.action",
	"jsp":            clientUrl + "custom/js/SetKingoEncypt.jsp",
	"getAchievement": clientUrl + "student/xscj.stuckcj_data.jsp",        //获取成绩接口
	"getClass":       clientUrl + "wsxk/xkjg.ckdgxsxdkchj_data10319.jsp", //获取课表接口（列表）
}

// request Headers
var UA string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
var MainHeaders = map[string]string{
	"User-Agent":      UA,
	"Accept":          "*/*",
	"Accept-Encoding": "gzip, deflate",
	"Connection":      "keep-alive",
}
