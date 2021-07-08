package utils

import (
	"net/http"
)

type MqttResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MResult(code int, data interface{}, msg string) string{
	return DataByJsonStr(MqttResponse{code, msg,data})
}

func Success(pub string) string{
	return MResult(http.StatusOK, map[string]interface{}{}, "成功")
}

func SuccessWithData(data interface{}) string{
	return MResult(http.StatusOK, data, "成功")
}

//500	执行中服务器内部错误
func MFailInternalServerError(pub string,message string) string{
	return MResult(http.StatusInternalServerError, map[string]interface{}{}, message)
}
//501   参数错误
func MFailNotImplemented(pub string, message string)string{
	return MResult(http.StatusNotImplemented, map[string]interface{}{}, message)
}
