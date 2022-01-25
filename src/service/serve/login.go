package serve

import (
	c "HeDa/src/client"
	s "HeDa/src/service/skeleton"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Henu_Res struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  string `json:"status"`
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
	var myClient c.MyClient
	defer rq.Body.Close()
	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		s.Res(rs, "错误", false)
		fmt.Println("11")
		return
	}
	fmt.Println("tttt")
	token, err := myClient.Logon(body)
	if err != nil {
		s.Res(rs, "错误", false)
		return
	}
	var tokenRes TokenRes
	tokenRes.Token = token
	s.Res(rs, tokenRes, true)
}
