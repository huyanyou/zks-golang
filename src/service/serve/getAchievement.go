package serve

import (
	c "HeDa/src/client"
	s "HeDa/src/service/skeleton"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//	获取成绩
func GetAchieve(rs http.ResponseWriter, rq *http.Request, p httprouter.Params) {
	defer rq.Body.Close()
	data := c.GetAchieve(rq.Body, rq.Header.Get("Authorization"))
	s.Res(rs, data, true)
}
