package services

import (
	"context"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services/callback/progress"
	"gitlab.example.com/zhangweijie/tool-sdk/services/callback/result"
	"time"
)

// getPendingWork 获取状态为 pending 的任务
func getPendingWork() (interface{}, error) {
	global.ValidExecutorChan.Lock()
	defer global.ValidExecutorChan.Unlock()
	pendingWork, err := models.GetWorkOrderCreateTime()
	if err != nil {
		logger.Error(schemas.WorkGetErr, err)
		return nil, err
	}

	if pendingWork.ID > 0 {
		// 更新任务状态为进行中
		err = models.UpdateWork(pendingWork.UUID, "status", global.WorkStatusDoingGo)
		if err != nil {
			logger.Error(schemas.WorkUpdateErr, err)
		}

		pendingWork, _ = models.GetWorkByUUID(pendingWork.UUID)

		return pendingWork, err

	}
	return nil, nil
}

func executeWork(work *global.Work) {
	go func() {
		defer func() {
			work.Cancel()
			global.ValidExecutorChan.WorkExecute <- true
		}()
		// 更新任务状态为进行中
		err := models.UpdateWork(work.WorkUUID, "status", global.WorkStatusDoing)
		if err != nil {
			logger.Error(schemas.WorkUpdateErr, err)
			return
		}

		params := make(map[string]interface{})
		validWork, err := models.GetWorkByUUID(work.WorkUUID)
		if err != nil {
			logger.Error(err.Error(), err)
			// 更新任务状态为进行中
			err = models.UpdateWork(work.WorkUUID, "status", global.WorkStatusPending)
			if err != nil {
				logger.Error(schemas.WorkUpdateErr, err)
				return
			}
			return
		}
		params["work"] = &validWork
		// 开始执行任务
		err = global.ValidExecutorIns.ExecutorMainFunc(work.Context, params)
		if err != nil {
			logger.Error(schemas.WorkExecuteErr, err)
			// 更新任务状态为失败
			err = models.UpdateWork(work.WorkUUID, "status", global.WorkStatusFailed)
			if err != nil {
				logger.Error(schemas.WorkUpdateErr, err)
			}
			return
		}
		// 更新任务状态为已完成
		err = models.UpdateWork(work.WorkUUID, "status", global.WorkStatusDone)
		if err != nil {
			logger.Error(schemas.WorkUpdateErr, err)
			return
		}
		return
	}()
	for {
		select {
		case <-work.Context.Done():
			return
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
				work := &global.Work{WorkUUID: oldSchema.(models.Work).UUID, Context: ctx, Cancel: cancel}
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
		// 推送结果
		case validResult := <-global.ValidResultChan:
			err := result.PushResult(validResult)
			if err != nil {
				logger.Warn(fmt.Sprintf("任务 %s 结果推送失败，错误为 %s !", validResult.WorkUUID, err))
			} else {
				// 修改任务推送状态为已完成
				err = models.UpdateWork(validResult.WorkUUID, "callback_status", global.WorkStatusDone)
				if err != nil {
					// 修改任务推送状态为已失败
					err = models.UpdateWork(validResult.WorkUUID, "callback_status", global.WorkStatusFailed)
					logger.Warn(fmt.Sprintf("任务 %s 结果推送失败，错误为 %s !", validResult.WorkUUID, err))
				}
			}
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
