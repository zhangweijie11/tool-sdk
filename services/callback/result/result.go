package result

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services/callback/util"
	"strconv"
	"strings"
)

// PushResult 发送任务结果
func PushResult(result *global.Result) error {
	err := models.DeleteResultByWorkUUID(result.WorkUUID)
	if err != nil {
		return err
	}
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

	switch strings.ToLower(result.CallbackType) {
	case strings.ToLower(global.CallbackTypeApi):
		err = util.CallbackAPI(result)
		return err
	case strings.ToLower(global.CallbackTypeMQ):
		err = util.CallbackMQ(result)
		return err
	case strings.ToLower(global.CallbackTypegRPC):
		err = util.CallbackgRPC(result)
		return err
	default:
		return errors.New(schemas.WorkCallbackTypeErr)
	}
}
