package serve

import (
	c "HeDa/src/client"
	s "HeDa/src/service/skeleton"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Henu_Res struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  string `json:"status"`
}

type LogonPrams struct {
	Params     string `json:"params"`
	Username   string `json:"username"`
	JSESSIONID string `json:"jsessionid"`
}

type TokenRes struct {
	Token string `json:"token"`
}

func Login(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
	var myClient c.MyClient
	params, err := myClient.Login()
	fmt.Println(rq.RemoteAddr)
	if err != "" {
		s.Res(rs, nil, false)
		return
	}
	s.Res(rs, params, true)
	return
}

func Logon(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
	fmt.Println(rq.RemoteAddr)
	var myClient c.MyClient
	var logon LogonPrams
	defer rq.Body.Close()
	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		s.Res(rs, "错误", false)
		fmt.Println("11")
		return
	}
	json.Unmarshal(body, &logon)
	req, _ := http.NewRequest("POST", "https://xk.henu.edu.cn/cas/logon.action", strings.NewReader(logon.Params))
	req.Header.Set("Host", "xk.henu.edu.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "JSESSIONID="+logon.JSESSIONID)
	res, err := myClient.Do(req)
	if err != nil {
		s.Res(rs, "服务错误!", false)
		return
	}
	// 请求henu成功后获得访问其他页面jsessionid权限
	data, _ := ioutil.ReadAll(res.Body)
	var henu_Res = Henu_Res{}
	json.Unmarshal(data, &henu_Res)
	if henu_Res.Status != "200" {
		fmt.Println("22")
		s.Res(rs, "服务错误!", false)
		return
	}
	token, _ := s.GenerateToken(logon.JSESSIONID)
	tokenRes := TokenRes{Token: token}
	s.Res(rs, tokenRes, true)
}
