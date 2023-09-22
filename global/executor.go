package global

import (
	"context"
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"time"
)

type ExecutorInterface interface {
	ValidWorkCreateParams(map[string]interface{}) error
	ExecutorMainFunc(context.Context, map[string]interface{}) error
}

type ExecutorIns struct{}

func NewExecutorIns() *ExecutorIns {
	return &ExecutorIns{}
}

func (ei *ExecutorIns) ExecutorMainFunc(ctx context.Context, params map[string]interface{}) error {
	for i := 0; i < 10; i++ {
		// 检查任务是否被取消
		select {
		case <-ctx.Done():
			return errors.New(schemas.WorkCancelErr)
		default:
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func (ei *ExecutorIns) ValidWorkCreateParams(map[string]interface{}) error {
	return nil
}
