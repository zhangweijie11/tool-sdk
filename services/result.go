package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"net/http"
	"strings"
	"sync"
)

type Result struct {
	sync.Mutex
	WorkUUID     string
	Result       map[string]interface{}
	CallbackType string
	CallbackUrl  string
}

// PushResult 发送任务结果
func (cr *Result) PushResult() error {
	cr.Lock()
	defer cr.Unlock()

	switch cr.CallbackType {
	case global.CallbackTypeApi:
		err := cr.callbackAPI()
		return err
	case global.CallbackTypeMQ:
		fmt.Println("------------>", global.CallbackTypeMQ)
		return nil
	case global.CallbackTypegRPC:
		fmt.Println("------------>", global.CallbackTypegRPC)
		return nil
	default:
		return errors.New(schemas.TypeErr)
	}
}

func (cr *Result) callbackAPI() error {
	validUrl := strings.TrimRight(cr.CallbackUrl, "/") + "/callback/result"
	var callbackResultParams = map[string]interface{}{
		"workUUID": cr.CallbackType,
		"result":   cr.Result,
	}
	// 将JSON对象编码为JSON字符串
	jsonData, err := json.Marshal(callbackResultParams)
	if err != nil {
		return err
	}

	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", validUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
