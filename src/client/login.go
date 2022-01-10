package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type BaseLogin struct {
	Username string
	Password string
}

type MyClient struct {
	http.Client
}

type LoginParams struct {
	SessionID string `json:"sessionid"`
	Deskey    string `json:"deskey"`
	Nowtime   string `json:"nowtime"`
}

//	用户登陆
func (m MyClient) Login() (params LoginParams, err string) {
	params = getSDN(&m)
	fmt.Println(params.Deskey)
	return params, ""
}

func (m MyClient) Logon(param string, username string, sessionid string) {
	// req, err := http.NewRequest("GET", Urls["logon"], strings.NewReader(param))
	// if err != nil {
	// 	return
	// }
	// req.Header.Add("")
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
	fmt.Println(params.SessionID)
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
