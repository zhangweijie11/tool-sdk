package services

import (
	"context"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services/progress"
	"gitlab.example.com/zhangweijie/tool-sdk/services/result"
	"time"
)

// getPendingWork 获取状态为 pending 的任务
func getPendingWork() (interface{}, error) {
	global.ValidExecutorChan.Lock()
	defer global.ValidExecutorChan.Unlock()
	pendingWork, err := models.GetWorkOrderCreateTime()
	if err != nil {
		logger.Error(schemas.GetWorkErr, err)
		return nil, err
	}

	if pendingWork.ID > 0 {
		// 更新任务状态为进行中
		err = models.UpdateWork(pendingWork.UUID, "status", global.WorkStatusDoingGo)
		if err != nil {
			logger.Error(schemas.UpdateWorkErr, err)
		}

		pendingWork, _ = models.GetWorkByUUID(pendingWork.UUID)

		return pendingWork, err

	}
	return nil, nil
}

func executeWork(work *global.Work) {
	go func() {
		defer func() {
			global.ValidExecutorChan.WorkExecute <- true
		}()
		// 更新任务状态为进行中
		err := models.UpdateWork(work.Work.UUID, "status", global.WorkStatusDoing)
		if err != nil {
			logger.Error(schemas.UpdateWorkErr, err)
			return
		}

		params := make(map[string]interface{})
		params["work"] = &work.Work
		// 开始执行任务
		err = global.ValidExecutorIns.ExecutorMainFunc(work.Context, params)
		if err != nil {
			logger.Error(schemas.ExecuteWorkErr, err)
			// 更新任务状态为失败
			err = models.UpdateWork(work.Work.UUID, "status", global.WorkStatusFailed)
			if err != nil {
				logger.Error(schemas.UpdateWorkErr, err)
			}
			return
		}
		// 更新任务状态为已完成
		err = models.UpdateWork(work.Work.UUID, "status", global.WorkStatusDone)
		if err != nil {
			logger.Error(schemas.UpdateWorkErr, err)
			return
		}
		return
	}()
	for {
		select {
		case <-work.Context.Done():
			return
		default:
		}
	}
}

// LoopExecuteWork  执行任务
func LoopExecuteWork() {
	for {
		select {
		// 限制全局并发任务执行数量
		case <-global.ValidExecutorChan.WorkExecute:
			oldSchema, err := getPendingWork()
			if err == nil && oldSchema != nil {
				global.ValidDoingWork.Lock()
				ctx, cancel := context.WithCancel(context.Background())
				work := &global.Work{Work: oldSchema.(models.Work), Context: ctx, Cancel: cancel}
				global.ValidDoingWork.DoingWorkMap[oldSchema.(models.Work).UUID] = work
				global.ValidDoingWork.Unlock()
				go executeWork(work)
			} else {
				time.Sleep(time.Second * 5)
				global.ValidExecutorChan.WorkExecute <- true
			}
		}
	}
}

// LoopProgressResult  执行任务进度推送和结果推送
func LoopProgressResult() {
	for {
		select {
		// 推送进度
		case validProgress := <-global.ValidProgressChan:
			err := progress.PushProgress(validProgress)
			logger.Warn(fmt.Sprintf("任务 %s 进度推送失败，错误为 %s !", validProgress.WorkUUID, err))
		case validResult := <-global.ValidResultChan:
			err := result.PushResult(validResult)
			logger.Warn(fmt.Sprintf("任务 %s 结果推送失败，错误为 %s !", validResult.WorkUUID, err))
		default:

		}
	}
}

// PauseWork 暂停任务
func PauseWork(workUUID string) {
	global.ValidDoingWork.Lock()
	defer global.ValidDoingWork.Unlock()
	work, ok := global.ValidDoingWork.DoingWorkMap[workUUID]
	if ok {
		work.Cancel()
		delete(global.ValidDoingWork.DoingWorkMap, workUUID)
	}
}
