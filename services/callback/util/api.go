package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"net/http"
	"strings"
	"time"
)

// CallbackAPI 通过 API形式回调
func CallbackAPI(validParams interface{}) error {
	var (
		validUrl     string
		validMessage map[string]interface{}
	)

	switch validParams.(type) {
	case *global.Progress:
		validUrl = strings.TrimRight(validParams.(*global.Progress).ProgressUrl, "/") + "/progress"
		validMessage = map[string]interface{}{
			"workUUID":   validParams.(*global.Progress).WorkUUID,
			"serverName": global.Config.Server.ServerName,
			"progress":   validParams.(*global.Progress).Progress,
		}
	case *global.Result:
		validUrl = strings.TrimRight(validParams.(*global.Result).CallbackUrl, "/") + "/callback/result"
		validMessage = map[string]interface{}{
			"workUUID":   validParams.(*global.Result).WorkUUID,
			"serverName": global.Config.Server.ServerName,
			"result":     validParams.(*global.Result).Result,
		}
	default:
		return errors.New(schemas.WorkCallbackErr)
	}

	// 将JSON对象编码为JSON字符串
	jsonData, err := json.Marshal(validMessage)
	if err != nil {
		return errors.New(schemas.JsonParseErr)
	}

	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", validUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New(schemas.APIConnectErr)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间为 5 秒
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 禁用 HTTPS 验证
		}}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(schemas.APIRequestErr)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New(schemas.APIResponseErr)
	}

	return nil
}
