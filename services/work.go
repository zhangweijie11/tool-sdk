package services

import (
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"time"
)

// getPendingWork 获取状态为 pending 的任务
func getPendingWork() (interface{}, error) {
	global.ValidWorkChan.Lock()
	defer global.ValidWorkChan.Unlock()
	pendingWork, err := models.GetWorkOrderCreateTime()
	if err != nil {
		logger.Error(schemas.GetWorkErr, err)
		return nil, err
	}

	if pendingWork.ID > 0 {
		// 更新任务状态为进行中
		err = models.UpdateWork(pendingWork.UUID, "status", global.WorkStatusDoing)
		if err != nil {
			logger.Error(schemas.UpdateWorkErr, err)
		}

		return pendingWork, err

	}
	return nil, nil
}

// LoopExecutionWork  执行任务
func LoopExecutionWork() {
	for {
		select {
		// 限制全局并发任务执行数量
		case <-global.ValidWorkChan.WorkExecute:
			go func() {
				oldSchema, err := getPendingWork()
				if err == nil {
					// 开始执行任务
					_, err = global.ValidExecutorIns.ExecutorMainFunc(oldSchema.(models.Work))
					if err != nil {
						logger.Error(schemas.ExecuteWorkErr, err)
					}
					// 更新任务状态为已完成
					err = models.UpdateWork(oldSchema.(models.Work).UUID, "status", global.WorkStatusDone)
					if err != nil {
						logger.Error(schemas.UpdateWorkErr, err)
					}

				}
				global.ValidWorkChan.WorkExecute <- true
			}()
		default:
			time.Sleep(time.Second * 5)
		}
	}
}
