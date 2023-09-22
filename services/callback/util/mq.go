package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"strings"
)

// CallbackMQ 通过 MQ 方式回调
func CallbackMQ(validParams interface{}) error {
	var (
		params       []string
		validMessage []byte
	)

	switch validParams.(type) {
	case *global.Progress:
		params = strings.Split(validParams.(*global.Progress).ProgressUrl, ",")
		// 要发送的 JSON 数据
		message := progressData{
			WorkUUID:  validParams.(*global.Progress).WorkUUID,
			Progress:  validParams.(*global.Progress).Progress,
			SeverName: global.Config.Server.ServerName,
		}
		// 将 JSON 数据编码为字节切片
		messageBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("无法编码 JSON 数据:", err)
			return errors.New(schemas.JsonParseErr)
		}
		validMessage = messageBytes
	case *global.Result:
		params = strings.Split(validParams.(*global.Result).CallbackUrl, ",")
		// 要发送的 JSON 数据
		message := resultData{
			WorkUUID:  validParams.(*global.Result).WorkUUID,
			Result:    validParams.(*global.Result).Result,
			SeverName: global.Config.Server.ServerName,
		}
		// 将 JSON 数据编码为字节切片
		messageBytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("无法编码 JSON 数据:", err)
			return errors.New(schemas.JsonParseErr)
		}
		validMessage = messageBytes
	default:
		return errors.New(schemas.WorkCallbackErr)
	}

	conn, err := amqp.Dial(params[0])
	if err != nil {
		// 处理连接错误
		return errors.New(schemas.MQConnectErr)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		// 处理通道创建错误
		return errors.New(schemas.MQChannelErr)
	}
	defer ch.Close()
	queueName := params[2]
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 是否持久化
		false,     // 是否自动删除
		false,     // 是否具有排他性
		false,     // 是否等待服务器的确认
		nil,       // 其他参数
	)
	if err != nil {
		// 处理队列声明错误
		return errors.New(schemas.MQQueueErr)
	}

	err = ch.Publish(
		params[1], // 交换机名称（空字符串表示使用默认交换机）
		queueName, // 队列名称
		false,     // 是否强制
		false,     // 是否立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        validMessage,
		},
	)
	if err != nil {
		// 处理消息发布错误
		return errors.New(schemas.MQMessageErr)
	}

	return err
}
