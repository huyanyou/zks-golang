package login

import (
	c "HeDa/src/client"
	s "HeDa/src/service/skeleton"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//	 用户登陆帐号密码
type User struct {
	Key string `json:"key"`
}

func Login(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
	var myClient c.MyClient
	params, err := myClient.Login()
	if err != "" {
		s.Res(rs, nil, false)
	}
	s.Res(rs, params, true)
}

func Logon(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {

}

//	判断是否有帐号密码
func isNameAndPass(rq *http.Request) *User {
	defer rq.Body.Close()
	var user *User
	if err := json.NewDecoder(rq.Body).Decode(&user); err != nil {
		user = nil
		return user
	}
	if user.Key == "" {
		user = nil
		return user
	}
	return user
}
