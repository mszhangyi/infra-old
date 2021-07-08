package base

import (
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/mszhangyi/infra"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var Service common.Service

func Server() common.Service {
	Check(Service)
	return Service
}

type ServerStarter struct {
	infra.BaseStarter
}

func (i *ServerStarter) Init() {
	Service = daprd.NewService(":" + strconv.Itoa(props.Port))
}

func (i *ServerStarter) Start() {
	logrus.Info("服务器正在运行,端口：" + strconv.Itoa(props.Port))
	if err := Service.Start(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("error listenning: %v", err)
	}
	logrus.Println("服务器正在退出 ")
}
