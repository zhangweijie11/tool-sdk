package result

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"net/http"
	"strconv"
	"strings"
)

// PushResult 发送任务结果
func PushResult(result *global.Result) error {
	jsonBytes, err := json.Marshal(result.Result)
	if err != nil {
		return err
	}

	validResult := &models.Result{
		UUID:     uuid.New().String(),
		WorkUUID: result.WorkUUID,
		TaskUUID: result.TaskUUID,
		Extra:    jsonBytes,
	}

	// 存储任务结果
	err = models.CreateResult(validResult)
	if err != nil {
		return err
	}

	// 修改任务进度
	err = models.UpdateWork(result.WorkUUID, "progress", strconv.FormatFloat(float64(100), 'f', 2, 64))
	if err != nil {
		return err
	}

	switch result.CallbackType {
	case global.CallbackTypeApi:
		err = callbackAPI(result)
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

func callbackAPI(result *global.Result) error {
	validUrl := strings.TrimRight(result.CallbackUrl, "/") + "/callback/result"
	var callbackResultParams = map[string]interface{}{
		"workUUID": result.CallbackType,
		"result":   result.Result,
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
