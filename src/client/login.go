package client

import (
	s "HeDa/src/service/skeleton"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type MyClient struct {
	http.Client
}

type LoginParams struct {
	SessionID string `json:"sessionid"`
	Deskey    string `json:"deskey"`
	Nowtime   string `json:"nowtime"`
}

type LogonPrams struct {
	Params     string `json:"params"`
	Username   string `json:"username"`
	JSESSIONID string `json:"jsessionid"`
}

type Henu_Res struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  string `json:"status"`
}

//	用户登陆
func (m MyClient) Login() (params LoginParams, err string) {
	params = getSDN(&m)
	return params, ""
}

func (m MyClient) Logon(body []byte) (string, error) {
	var logonParams LogonPrams
	json.Unmarshal(body, &logonParams)
	req, _ := http.NewRequest("POST", Urls["logon"], strings.NewReader(logonParams.Params))
	req.Header.Set("Host", "xk.henu.edu.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "JSESSIONID="+logonParams.JSESSIONID)
	fmt.Println("1531515")
	res, err := m.Do(req)
	if err != nil {
		return "", err
	}
	// 请求henu成功后获得访问其他页面jsessionid权限
	data, _ := ioutil.ReadAll(res.Body)
	var henu_Res = Henu_Res{}
	json.Unmarshal(data, &henu_Res)
	if henu_Res.Status != "200" {
		return "", err
	}
	token, _ := s.GenerateToken(logonParams.JSESSIONID)
	// 定时延长session
	go func() {
		t := time.NewTicker(2 * time.Hour)
		for {
			<-t.C
			m.LongSession(logonParams.JSESSIONID, henu_Res.Result)
		}
	}()
	return token, nil
}

//	客户端获取henu的sessionid timenow deskey
func getSDN(m *MyClient) (params LoginParams) {
	req, err := http.NewRequest("GET", Urls["login"], nil)
	// SetHeaders(req, MainHeaders)
	req.Header.Set("Accept", MainHeaders["Accept"])
	req.Header.Set("Accept-Encoding", "")
	if err != nil {
		params = LoginParams{}
		return params
	}
	res, err := m.Do(req)
	if err != nil {
		return LoginParams{}
	}
	fmt.Println(res.Cookies())
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return LoginParams{}
	}
	params.SessionID, params.Deskey, params.Nowtime = getSDN_query(doc.Find("head script"), m, req)
	return params
}

// 客户端获取jsp页面中的sessionid timenow deskey封装
func getSDN_query(doc *goquery.Selection, m *MyClient, req *http.Request) (string, string, string) {
	reg := regexp.MustCompile(`\w{32}\.kingo154`)
	sessionid := reg.FindString(doc.Nodes[0].LastChild.Data)
	jspUrl := doc.Nodes[1].Attr[1].Val
	req, _ = http.NewRequest("GET", jspUrl, nil)
	req.Header.Set("Cookie", "JSESSIONID="+sessionid)
	req.Header.Set("Host", "xk.henu.edu.cn")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	res, err := m.Do(req)
	if err != nil {
		fmt.Println("错了")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	reg = regexp.MustCompile(`_deskey\s=\s'\d{19}`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	deskey := result[0][0][11:]
	reg = regexp.MustCompile(`_nowtime\s=\s'\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`)
	result = reg.FindAllStringSubmatch(string(body), -1)
	nowtime := result[0][0][12:]

	return sessionid, deskey, nowtime
}

// 使session长期有效
func (m MyClient) LongSession(jessionid string, url string) {
	req, _ := http.NewRequest("GET", "https://xk.henu.edu.cn"+url, nil)
	req.Header.Set("Host", "xk.henu.edu.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "JSESSIONID="+jessionid)
	m.Do(req)
}
