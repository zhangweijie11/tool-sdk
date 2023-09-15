package example

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	tool "gitlab.example.com/zhangweijie/tool-sdk/cmd"
	"gitlab.example.com/zhangweijie/tool-sdk/example/middlerware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/example/models"
	"gitlab.example.com/zhangweijie/tool-sdk/example/routers"
	"gitlab.example.com/zhangweijie/tool-sdk/example/services/fingerprint"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	toolSchemas "gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	toolModels "gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/option"
	toolRouters "gitlab.example.com/zhangweijie/tool-sdk/routers"
)

type executorIns struct {
	global.ExecutorIns
}

// ValidWorkCreateParams 验证任务参数
func (ei *executorIns) ValidWorkCreateParams(params map[string]interface{}) (err error) {
	var schema = new(schemas.FingerprintParams)
	if err = toolSchemas.CustomBindSchema(params, schema, schemas.RegisterValidatorRule); err == nil || err.Error() == "" {
		err = toolSchemas.ValidateLength(schema.Urls, 0, 100)
	}

	return err
}

// ExecutorMainFunc 任务执行主函数（可自由发挥）
func (ei *executorIns) ExecutorMainFunc(ctx context.Context, params map[string]interface{}) error {
	errChan := make(chan error)
	go func() {
		work, err := toolModels.GetWorkByUUID(params["workUUID"].(string))
		if err != nil {
			logger.Error(toolSchemas.GetWorkErr, err)
			errChan <- err
		} else {
			var validParams fingerprint.FingerprintParams
			err = json.Unmarshal(work.Params, &validParams)
			if err != nil {
				logger.Error(toolSchemas.JsonParseErr, err)
				errChan <- err
			} else {
				if len(validParams.Urls) < 1 {
					errChan <- errors.New(toolSchemas.WorkTargetErr)
				} else {
					err = fingerprint.FingerprintMainWorker(ctx, &work, &validParams)
					errChan <- err
				}
			}
		}
	}()
	select {
	case <-ctx.Done():
		return errors.New(toolSchemas.CancelWorkErr)
	case err := <-errChan:
		return err
	}
}

// Start 工具主函数
func Start() {
	defaultOption := option.GetDefaultOption()
	defaultOption.ExecutorIns = &executorIns{}
	defaultOption.ValidModels = []interface{}{&toolModels.Work{}, &toolModels.Task{}, &toolModels.Result{}, &models.Technologies{}}
	defaultOption.ValidRouters = []func(*gin.Engine) gin.IRoutes{toolRouters.InitPingRouter, toolRouters.InitWorkRouter, routers.InitTechnologiesRouter}
	tool.Start(defaultOption)
}
