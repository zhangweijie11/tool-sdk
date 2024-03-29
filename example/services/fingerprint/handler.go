package fingerprint

import (
	"context"
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"sync"
	"time"
)

type Worker struct {
	ID         int // 任务执行者 ID
	Ctx        context.Context
	Wg         *sync.WaitGroup
	TaskChan   chan Task     // 子任务通道
	ResultChan chan []string // 子任务结果通道
}

type Task struct {
	WorkUUID  string // 总任务 ID
	TaskUUID  string // 子任务 ID
	TargetUrl string // 子任务目标网站
}

// NewWorker 初始化 worker
func NewWorker(ctx context.Context, wg *sync.WaitGroup, id int, taskChan chan Task, resultChan chan []string) *Worker {
	return &Worker{
		ID:         id,
		Ctx:        ctx,
		Wg:         wg,
		TaskChan:   taskChan,
		ResultChan: resultChan,
	}
}

type FingerprintParams struct {
	Urls []string `json:"urls"`
}

// GroupFingerprintWorker 指纹识别方法
func (w *Worker) GroupFingerprintWorker() {
	go func() {
		defer w.Wg.Done()
		for task := range w.TaskChan {
			select {
			case <-w.Ctx.Done():
				return
			default:
				var fingerResults = []string{task.TargetUrl}
				time.Sleep(time.Second)
				select {
				case <-w.Ctx.Done():
					return
				default:
					w.ResultChan <- fingerResults
				}
			}
		}
	}()
}

func FingerprintMainWorker(ctx context.Context, work *toolModels.Work, validParams *FingerprintParams) error {
	quit := make(chan struct{})
	go func() {
		defer close(quit)
		onePercent := float64(100 / len(validParams.Urls))
		taskChan := make(chan Task, len(validParams.Urls))
		resultChan := make(chan []string, len(validParams.Urls))
		var wg sync.WaitGroup
		// 创建并启动多个工作者
		for i := 0; i < global.Config.Server.Worker; i++ {
			worker := NewWorker(ctx, &wg, i, taskChan, resultChan)
			worker.GroupFingerprintWorker()
			wg.Add(1)
		}
		go func() {
			// 通知消费者所有任务已经推送完毕
			defer close(taskChan)
			for _, url := range validParams.Urls {
				task := Task{
					WorkUUID:  work.UUID, // 总任务 ID
					TargetUrl: url,       // 子任务目标网站
				}
				taskChan <- task
			}
		}()

		go func() {
			wg.Wait()
			// 通知消费者所有任务结果已经推送完毕
			close(resultChan)
		}()

		var finalResult [][]string
		for fingerprintResult := range resultChan {
			if work.ProgressType != "" && work.ProgressUrl != "" {
				pushProgress := &global.Progress{WorkUUID: work.UUID, ProgressType: work.ProgressType, ProgressUrl: work.ProgressUrl, Progress: 0}
				pushProgress.Progress += onePercent
				// 回传进度
				global.ValidProgressChan <- pushProgress
			}
			fmt.Println("------------>", fingerprintResult)
		}

		if work.CallbackType != "" && work.CallbackUrl != "" {
			pushResult := &global.Result{WorkUUID: work.UUID, CallbackType: work.CallbackType, CallbackUrl: work.CallbackUrl, Result: map[string]interface{}{"result": finalResult}}
			// 回传结果
			global.ValidResultChan <- pushResult
		}

	}()

	select {
	case <-ctx.Done():
		return errors.New(schemas.WorkCancelErr)
	case <-quit:
		return nil
	}
}
