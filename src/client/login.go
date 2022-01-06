package client

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

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
	P_username        string
	P_password        string
	Username          string
	Password          string
	Randnumber        int8
	IsPasswordPolicy  int8
	Txt_mm_expression int8
	Txt_mm_length     int8 // 密码长度
	Txt_mm_userzh     int8 //  判断密码是否包含帐号
	Hid_flag          int8 // 默认是1
	SessionID         string
	Deskey            string
	Nowtime           string
}

//	用户登陆
func (m MyClient) Login(username string, pwd string) {
	params := getSDN(&m)
	params.Username = base64Encode(username + ";;" + params.SessionID)
	params.Password = pwd
	params.IsPasswordPolicy = isPasswordPolicy(username, pwd)
	params.Txt_mm_expression, params.Txt_mm_length, params.Txt_mm_userzh = checkpwd(username, pwd)
	params.P_password = "_u"
	params.P_username = "_p"
}

// 密码验证规则
func checkpwd(username string, pwd string) (txt_mm_expression int8, txt_mm_length int8, txt_mm_userzh int8) {
	txt_mm_expression = 0
	txt_mm_length = int8(utf8.RuneCountInString(pwd))
	txt_mm_userzh = '0'
	for _, v := range pwd {
		txt_mm_expression |= charType(v)
	}
	if l := strings.Contains(strings.ToLower(pwd), strings.ToLower(username)); l {
		txt_mm_userzh = '1'
		return txt_mm_expression, txt_mm_length, txt_mm_userzh
	}
	return txt_mm_expression, txt_mm_length, txt_mm_userzh
}

func charType(uni rune) int8 {
	if uni >= 48 && uni <= 57 {
		return 8
	}
	if uni >= 97 && uni <= 122 {
		return 4
	}
	if uni >= 65 && uni <= 90 {
		return 2
	}
	return 1
}

func isPasswordPolicy(username string, pwd string) int8 {
	if username == "" || pwd == "" || username == pwd {
		return '0'
	}
	pwdLen := len(pwd)
	if pwdLen < 6 {
		return '0'
	}
	return '1'
}

// base64加密
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

//	对params进行拼接
func paramsSplicing(l LoginParams) string {
	return l.P_username + "=" + l.Username + "&" + l.P_password + "=" + l.Password + "&randnumber=" + string(rune(l.Randnumber)) +
		"&isPasswordPolicy" + string(rune(l.IsPasswordPolicy)) + "&txt_mm_expression=" + string(rune(l.Txt_mm_expression)) +
		"&txt_mm_length=" + string(rune(l.Txt_mm_userzh)) + "&txt_mm_userzh=" + string(rune(l.Txt_mm_userzh)) + "&hid_flag=" + string(rune(l.Hid_flag)) +
		"&hidlag=1"
}

//	对拼接后的params进行编码
func paramGetEncParams(p string) string {
	return ""
}

//	客户端获取henu的sessionid timenow deskey
func getSDN(m *MyClient) *LoginParams {
	var params LoginParams
	req, err := http.NewRequest("GET", Urls["login"], nil)
	if err != nil {
		return nil
	}
	SetHeaders(req, MainHeaders)
	req.Header.Set("Accept-Encoding", "")
	res, err := m.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}
	params.SessionID, params.Deskey, params.Nowtime = getSDN_query(doc.Find("head script"), m, req)
	return &params
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
