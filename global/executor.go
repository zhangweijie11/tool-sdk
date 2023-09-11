package global

import (
	"context"
	"errors"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"time"
)

type ExecutorInterface interface {
	ValidWorkCreateParams(map[string]interface{}) (string, error)
	ExecutorMainFunc(context.Context, interface{}) error
}

type ExecutorIns struct{}

func NewExecutorIns() *ExecutorIns {
	return &ExecutorIns{}
}

func (ei *ExecutorIns) ExecutorMainFunc(ctx context.Context, params interface{}) error {
	for i := 0; i < 20; i++ {
		fmt.Println("------------>任务", params, i)
		time.Sleep(1 * time.Second)
		// 检查任务是否被取消
		select {
		case <-ctx.Done():
			return errors.New(schemas.CancelWorkErr)
		default:

		}
	}
	return nil
}

func (ei *ExecutorIns) ValidWorkCreateParams(params map[string]interface{}) (string, error) {
	return "", nil
}
