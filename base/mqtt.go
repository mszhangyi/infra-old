package base

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mszhangyi/infra"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	PubClientSize = 6
)

var (
	pubChanMap map[int]chan *ResponseMsg
	mClient    mqtt.Client
)

type ResponseMsg struct {
	Topic string
	Msg   interface{}
}

type MQttStarter struct {
	infra.BaseStarter
}

func (t *MQttStarter) Init() {
	pubChanMap = make(map[int]chan *ResponseMsg)
	for i := 0; i < PubClientSize; i++ {
		ch := make(chan *ResponseMsg, 3000)
		pubChanMap[i] = ch
		go t.startPublishMQtt(i)
	}
	//主程序客户端
	opts := getMqttOpts(props.Name)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mClient = client
}
func (t *MQttStarter) Stop() {
	//关闭服务的时候
}

/*func Start(ctx infra.StarterContext) {
	errChan := make(chan error)
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()
	getErr := <-errChan //只要报错 或者service关闭阻塞在这里的会进行下去
	logrus.Error(getErr)
}*/

//--------------------------------------------------------------------------	启动订阅Topic
func StartSubscribe(topic string, handler func(message mqtt.Message)) {
	var handlerMeg mqtt.MessageHandler
	if handler == nil {
		handlerMeg = func(client mqtt.Client, message mqtt.Message) {}
	} else {
		handlerMeg = func(client mqtt.Client, message mqtt.Message) {
			go handlerRecover(handler, message)
		}
	}
	token := mClient.Subscribe(topic, 0, handlerMeg)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

//异常捕捉   防止panic终止程序
func handlerRecover(handler func(message mqtt.Message), message mqtt.Message) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	handler(message)
}

/**
发送消息
*/
func SendMsg(pub string, msg string) {
	response := &ResponseMsg{
		Topic: pub,
		Msg:   msg,
	}
	rand.Seed(time.Now().UnixNano())
	pubChanMap[rand.Intn(PubClientSize)] <- response
}

//---------------------------------------------------------------------------------   启动发布Topic
func (t *MQttStarter) startPublishMQtt(index int) {
	clientId := fmt.Sprintf(props.Name+"-"+fmt.Sprint(rand.Int63n(time.Now().UnixNano()))+"-%d", index)
	opts := getMqttOpts(clientId)
	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for {
		select {
		case t, ok := <-pubChanMap[index]:
			if !ok {
				return
			}
			c.Publish(t.Topic, 0, false, t.Msg)
		}
	}
	c.Disconnect(3000)
}

func getMqttOpts(ClientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(props.EmqAddr).SetClientID(ClientId)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectTimeout(60 * time.Second)
	opts.SetUsername(props.MqttUser)
	opts.SetPassword(props.MqttPwd)
	opts.SetProtocolVersion(4)
	opts.SetAutoReconnect(true)                 //设置自动重新连接
	opts.SetOnConnectHandler(onConnectCallBack) //设置初始连接时和自动重新连接时调用的函数。
	return opts
}

var onConnectCallBack mqtt.OnConnectHandler = func(client mqtt.Client) {
	/*options := client.OptionsReader()
	clientId := options.ClientID()
	logrus.Info("mqtt " +clientId + " client connect success ")*/
}
