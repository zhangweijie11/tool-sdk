package progress

import (
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services/callback/util"
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
		switch strings.ToLower(progress.ProgressType) {
		case strings.ToLower(global.CallbackTypeApi):
			err = util.CallbackAPI(progress)
			return err
		case strings.ToLower(global.CallbackTypeMQ):
			err = util.CallbackMQ(progress)
			return nil
		case strings.ToLower(global.CallbackTypegRPC):
			fmt.Println("------------>", global.CallbackTypegRPC)
			return nil
		default:
			return errors.New(schemas.WorkCallbackTypeErr)
		}
	}

	return err
}
