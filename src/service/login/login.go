package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	s "HeDa/src/service/skeleton"

	"github.com/julienschmidt/httprouter"
)

//	 用户登陆帐号密码
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
	user := isNameAndPass(rq)
	if user == nil {
		fmt.Println("false")
		s.Res(rs, nil, false)
		return
	}
	s.Res(rs, user, true)
}

//	判断是否有帐号密码
func isNameAndPass(rq *http.Request) *User {
	defer rq.Body.Close()
	var user *User
	if err := json.NewDecoder(rq.Body).Decode(&user); err != nil {
		user = nil
		return user
	}
	if user.Username == "" {
		user = nil
		return user
	}
	if user.Password == "" {
		user = nil
		return user
	}
	return user
}
