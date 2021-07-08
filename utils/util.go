package utils

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"sync"
)

// ----------------------------------------------------------    json  操作
var(
	JsonApi jsoniter.API
	oJson sync.Once
)

func init()  {
	oJson.Do(func() {
		JsonApi = jsoniter.ConfigCompatibleWithStandardLibrary
	})
}

func DataByJsonByte(params interface{}) []byte {
	by ,err := JsonApi.Marshal(params)
	if err != nil {
		logrus.Error("DataByJsonStr err ", err)
	}
	return by
}

func StrJsonByData(str string, data interface{}) {
	err := JsonApi.Unmarshal([]byte(str), data)
	if err != nil {
		logrus.Error("StrJsonByData err ", err)
	}
}

func DataByJsonStr(params interface{}) string {
	by ,err := JsonApi.Marshal(params)
	if err != nil {
		logrus.Error("DataByJsonStr err ", err)
	}
	return string(by)
}

func ByteJsonByData(by []byte,data interface{}) error {
	err := JsonApi.Unmarshal(by, data)
	if err != nil {
		return err
	}
	return nil
}

type ClientResp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}
func SendClient(result interface{},appId ,methodName string, content interface{}) error{
	client, err := dapr.NewClient()
	if err != nil {
		logrus.Info(err)
		return err
	}
	resp,err := client.InvokeMethodWithCustomContent(context.Background(),appId,methodName,"post","application/json",content)
	if err != nil {
		logrus.Info("请求"+methodName+"错误",err)
		return  err
	}
	err = ByteJsonByData(resp, result)
	if err != nil {
		logrus.Info("解析"+methodName+"错误", err)
		return err
	}
	return nil
}

/*func SendClient(appId ,methodName string, content interface{}) (interface{}, error){
	client, err := dapr.NewClient()
	if err != nil {
		logrus.Info(err)
		return nil, err
	}
	defer client.Close()
	resp,err := client.InvokeMethodWithCustomContent(context.Background(),appId,methodName,"post","application/json",content)
	if err != nil {
		logrus.Info("请求"+methodName+"错误",err)
		return nil, err
	}
	result := make(map[string]interface{})
	err = ByteJsonByData(resp, &result)
	if err != nil {
		logrus.Info("解析"+methodName+"错误", err)
		return nil, err
	}
	if result["code"].(float64) > 0 {
		return nil, errors.New(result["msg"].(string))
	}
	return result["data"], nil
}
*/