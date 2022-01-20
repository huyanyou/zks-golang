package skeleton

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
)

//	日志对象
type Log struct{}

//	请求的基本信息
type baseInfo struct {
	path   *url.URL
	ip     string
	method string
}

func (l *Log) Logger(path string, ip string, method string, remoterArr string) {
	log.Default().Printf("%s ip:%s requestURL: %s Method:%s port: %s", time.Now().Format("2021-02-02 12:12:12"), ip, path, method, remoterArr)
}

var L Log

func GlobalMiddle(w http.ResponseWriter, r *http.Request) {
	path := r.RemoteAddr
	if path != "/" {
		tokenStr := r.Header.Get("Authorization")
		_, err := ParseToken(tokenStr)
		if err != nil {
			Res(w, "错误", false)
			return
		}
	}
	baseInfo := baseInfo{path: r.URL, ip: r.Host, method: r.Method}
	L.Logger(baseInfo.path.String(), baseInfo.ip, baseInfo.method, r.RemoteAddr)
}

//	jwt验证中间件
func MiddleAuth(router httprouter.Handle) httprouter.Handle {
	return func(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
		path := rq.RemoteAddr
		if path != "/" {
			tokenStr := rq.Header.Get("Authorization")
			claims, err := ParseToken(tokenStr)
			if err != nil {
				Res(rs, "错误", false)
				return
			}
			rq.Header.Set("Authorization", claims.Jsessionid)
			router(rs, rq, p)
			return
		}
		router(rs, rq, p)
	}
}
