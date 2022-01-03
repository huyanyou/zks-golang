package skeleton

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResStruct struct {
	Code int16       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

//	响应成功的格式
func Res(rs http.ResponseWriter, data interface{}, isSuccess bool) {
	switch isSuccess {
	case true:
		var resStruct ResStruct
		rs.WriteHeader(200)
		resStruct.Code = 200
		resStruct.Msg = "服务完成"
		resStruct.Data = data
		rs.Header().Set("Content-Type", "application/json")
		jsonData, _ := json.Marshal(resStruct)
		fmt.Println("true")
		rs.Write(jsonData)
		break
	case false:
		var resStruct ResStruct
		rs.WriteHeader(404)
		resStruct.Code = 404
		resStruct.Msg = "请求失败"
		jsonData, _ := json.Marshal(resStruct)
		rs.Write(jsonData)
		break
	default:
		break
	}

}
