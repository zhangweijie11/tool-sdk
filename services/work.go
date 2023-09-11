package services

import (
	"context"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
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
		err := models.UpdateWork(work.WorkUUID, "status", global.WorkStatusDoing)
		if err != nil {
			logger.Error(schemas.UpdateWorkErr, err)
			return
		}

		// 开始执行任务
		err = global.ValidExecutorIns.ExecutorMainFunc(work.Context, work.WorkUUID)
		if err != nil {
			logger.Error(schemas.ExecuteWorkErr, err)
			return
		}
		// 更新任务状态为已完成
		err = models.UpdateWork(work.WorkUUID, "status", global.WorkStatusDone)
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
				work := &global.Work{WorkID: oldSchema.(models.Work).ID, WorkUUID: oldSchema.(models.Work).UUID, Context: ctx, Cancel: cancel}
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
