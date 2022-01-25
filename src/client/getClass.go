package client

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GetClassPrams struct {
	Xn string `json:"xn"` //开始学年
	Xq string `json:"xq"` //上或下学期
	Xh string `json:"xh"` //学号
}

type ClassRes struct {
	Num         string `json:"classNum"`    // 上课班号
	ClassName   string `json:"className"`   // 课程名称
	TotleTime   string `json:"totleTime"`   // 总学时
	Credit      string `json:"credit"`      // 学分
	StudyNature string `json:"studyNature"` // 修读性质
	Teacher     string `json:"teacher"`     // 任课老师
	Status      string `json:"status"`      // 选课状态
	Book        string `json:"book"`        // 教材
	Place       string `json:"place"`       // 上课地点
}

func GetClass(data io.ReadCloser, jesssionid string) []ClassRes {
	defer data.Close()
	serverData, _ := ioutil.ReadAll(data)
	var getClassPrams GetClassPrams
	json.Unmarshal(serverData, &getClassPrams)
	resData := ClassReq(getClassPrams, jesssionid)
	return resData

}

// 获取课表  (暂时先这样写，后面再抽离，写的太丑了)
func ClassReq(data GetClassPrams, jsessionid string) []ClassRes {
	var myClient MyClient
	req, _ := http.NewRequest("GET", "https://xk.henu.edu.cn/student/xkjg.wdkb.jsp?menucode=S20301", nil)
	req.Header.Set("Referer", "https://xk.henu.edu.cn//frame/jw/teacherstudentmenu.jsp?menucode=S203")
	req.Header.Set("Cookie", "JSESSIONID="+jsessionid)
	req.Header.Set("User-Agent", UA)
	res, _ := myClient.Do(req)
	defer res.Body.Close()
	resHtmlGBK, _ := ioutil.ReadAll(res.Body)
	resHtmlUTF, _ := GbkToUtf8(resHtmlGBK)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(resHtmlUTF)))
	data.Xh, _ = dom.Find("#xh").Attr("value")
	encodeParams := paramsTran(data)
	encodeParams = base64.RawURLEncoding.EncodeToString([]byte(encodeParams))
	req, _ = http.NewRequest("GET", Urls["getClass"]+"?params="+encodeParams, nil)
	req.Header.Set("Referer", "https://xk.henu.edu.cn/student/xkjg.wdkb.jsp?menucode=S20301")
	req.Header.Set("Cookie", "JSESSIONID="+jsessionid)
	req.Header.Set("User-Agent", UA)
	res, _ = myClient.Do(req)
	defer res.Body.Close()
	resHtmlGBK, _ = ioutil.ReadAll(res.Body)
	resHtmlUTF, _ = GbkToUtf8(resHtmlGBK)
	var classRes ClassRes
	var resData []ClassRes
	v := reflect.ValueOf(&classRes).Elem()
	dom, _ = goquery.NewDocumentFromReader(strings.NewReader(string(resHtmlUTF)))
	dom.Find("body tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i > 8 {
				return false
			}
			v.Field(i).SetString(s.Text())
			return true
		})
		resData = append(resData, classRes)
	})
	return resData
}

func paramsTran(data GetClassPrams) string {
	return "xn=" + data.Xn + "&xq=" + data.Xq + "&xh=" + data.Xh
}
