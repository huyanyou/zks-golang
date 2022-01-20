package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// sjxz: sjxz3
// ysyx: yscj
// zx: 1
// fx: 1
// btnExport: （无法解码值）
// xn: 2020
// xn1: 2021
// xq: 0		0上学期  1下学期
// ysyxS: on
// sjxzS: on
// zxC: on
// fxC: on
// menucode_current:
type AchievementPrams struct {
	Sjxz             string `json:"sjxz"`
	Ysyx             string `json:"ysyx"`
	Zx               string `json:"zx"`
	Fx               string `json:"fx"`
	BtnExport        string `json:"btnExport"`
	Xn               string `json:"xn"`
	Xn1              string `json:"xn1"`
	Xq               string `json:"xq"`
	YsyxS            string `json:"ysyxS"`
	SjxzS            string `json:"sjxzS"`
	ZxC              string `json:"zxC"`
	FxC              string `json:"fxC"`
	Menucode_current string `json:"menucode_current"`
}

type AchiementRes struct {
	Course       string `json:"course"`        //课程名称
	Credit       string `json:"credit"`        //学分
	StudyTime    string `json:"study_time"`    //总学时
	Category     string `json:"category"`      // 类别
	StudyNature  string `json:"study_nature"`  // 修读性质
	AccessMethod string `json:"access_method"` //考核方式
	GetMethod    string `json:"get_method"`    //取得方式
	Fraction     string `json:"fraction"`      //分数
}

func GetAchieve(serviceBody io.ReadCloser, jsessionid string) []AchiementRes {
	serviceData, _ := ioutil.ReadAll(serviceBody)
	var AchievementPrams AchievementPrams
	json.Unmarshal(serviceData, &AchievementPrams)

	// 客户端请求数据
	req, _ := http.NewRequest("POST", Urls["getAchievement"], strings.NewReader(AchievementPrams.splicing()))
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", UA)
	req.Header.Set("Cookie", "JSESSIONID="+jsessionid)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	var myClient MyClient
	res, _ := myClient.Do(req)
	fmt.Println(jsessionid)
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	data, _ = GbkToUtf8(data)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))

	//	将请求的数据封装到数组中
	var c []AchiementRes
	var achiementRes AchiementRes
	d := doc.Find("body table tbody").Eq(1).Find("tr")
	d.Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch {
			case i == 1:
				achiementRes.Course = s.Text()
			case i == 2:
				achiementRes.Credit = s.Text()
			case i == 3:
				achiementRes.StudyTime = s.Text()
			case i == 4:
				achiementRes.Category = s.Text()
			case i == 5:
				achiementRes.StudyNature = s.Text()
			case i == 6:
				achiementRes.AccessMethod = s.Text()
			case i == 7:
				achiementRes.GetMethod = s.Text()
			case i == 8:
				achiementRes.Fraction = s.Text()
			default:
				break
			}
		})
		c = append(c, achiementRes)
	})
	return c
}

func (a AchievementPrams) splicing() string {
	return "sjxz=" + a.Sjxz + "&ysyx=" + a.Ysyx + "&zx=" + a.Zx + "&fx=" + a.Fx + "&btnExport=" + a.BtnExport + "&xn=" + a.Xn + "&xn1=" + a.Xn1 + "&xq=" + a.Xq +
		"&ysyxS=" + a.YsyxS + "&sjxzS=" + a.YsyxS + "&zxC=" + a.ZxC + "&fxC=" + a.FxC + "&menucode_current=" + a.Menucode_current
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
