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

type Progress struct {
	sync.Mutex
	WorkUUID     string
	Progress     float32
	ProgressType string
	ProgressUrl  string
}

// PushProgress 发送任务进度
func (p *Progress) PushProgress() error {
	p.Lock()
	defer p.Unlock()

	if p.Progress > 95 && p.Progress != 100 {
		p.Progress = 95
	}
	switch p.ProgressType {
	case global.CallbackTypeApi:
		err := p.callbackAPI()
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

func (p *Progress) callbackAPI() error {
	validUrl := strings.TrimRight(p.ProgressUrl, "/") + "/progress"
	var progressParams = map[string]interface{}{
		"workUUID": p.WorkUUID,
		"progress": p.Progress,
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
