package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func callbackAPI(progressParams map[string]string) error {
	validUrl := strings.TrimRight(progressParams["progressUrl"], "/") + "/progress"
	delete(progressParams, "progressUrl")
	// 将JSON对象编码为JSON字符串
	jsonData, err := json.Marshal(progressParams)
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

// PushProgress 推送任务进度
func PushProgress(workUUID, progress string) error {
	// 查询数据是否存在
	work, err := models.GetWorkByUUID(workUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(schemas.RecordNotFoundErr)
	}
	progressParams := make(map[string]string)
	progressParams["workUUID"] = work.UUID
	progressParams["progress"] = progress
	progressParams["progressUrl"] = work.ProgressUrl
	switch work.ProgressType {
	case global.CallbackTypeApi:
		err = callbackAPI(progressParams)
		return err
	case global.CallbackTypeMQ:
		fmt.Println("------------>", global.CallbackTypeMQ)
	case global.CallbackTypegRPC:
		fmt.Println("------------>", global.CallbackTypegRPC)
	}
	return nil
}
