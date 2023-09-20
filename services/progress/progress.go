package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"net/http"
	"strconv"
	"strings"
)

// PushProgress 发送任务进度
func PushProgress(progress *global.Progress) (err error) {
	if progress.Progress > 95 && progress.Progress != 100 {
		progress.Progress = 95
	}

	// 修改任务进度
	err = models.UpdateWork(progress.WorkUUID, "progress", strconv.FormatFloat(progress.Progress, 'f', 2, 64))
	if err == nil {
		switch progress.ProgressType {
		case global.CallbackTypeApi:
			err = callbackAPI(progress)
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

	return err
}

func callbackAPI(progress *global.Progress) error {
	validUrl := strings.TrimRight(progress.ProgressUrl, "/") + "/progress"
	var progressParams = map[string]interface{}{
		"workUUID": progress.WorkUUID,
		"progress": progress.Progress,
	}
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
