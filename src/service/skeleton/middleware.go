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

func (l *Log) Logger(path string, ip string, method string) {
	log.Default().Printf("%s ip:%s requestURL: %s Method:%s", time.Now().Format("2021-02-02 12:12:12"), ip, path, method)
}

var L Log

func GlobalMiddle(w http.ResponseWriter, r *http.Request) {
	baseInfo := baseInfo{path: r.URL, ip: r.Host, method: r.Method}
	L.Logger(baseInfo.path.String(), baseInfo.ip, baseInfo.method)
}
