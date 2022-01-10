package skeleton

import (
	"log"
	"net/http"
	"net/url"
	"time"
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
	(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	baseInfo := baseInfo{path: r.URL, ip: r.Host, method: r.Method}
	L.Logger(baseInfo.path.String(), baseInfo.ip, baseInfo.method, r.RemoteAddr)
}
