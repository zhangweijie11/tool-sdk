package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"gitlab.example.com/zhangweijie/tool-sdk/controller"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services"
	"gorm.io/gorm"
)

// WorkCreateApi 创建总任务
func WorkCreateApi(c *gin.Context) {
	var schema = new(schemas.WorkCreateSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = global.ValidExecutorIns.ValidWorkCreateParams(schema.Params)
		if err != nil {
			schemas.Fail(c, err.Error())
		} else {
			if global.Config.Database.Activate {
				jsonBytes, err := json.Marshal(schema.Params)
				if err != nil {
					logger.Error(schemas.JsonParseErr, err)
					schemas.Fail(c, schemas.JsonParseErr)
				}

				if schema.WorkUUID == "" {
					schema.WorkUUID = uuid.New().String()
				}

				work := &models.Work{
					UUID:     schema.WorkUUID,
					Params:   jsonBytes,
					Status:   global.WorkStatusPending,
					Source:   schema.Source,
					Priority: schema.Priority,
					Retry:    3,
				}

				err = models.CreateWok(work)

				if err != nil {
					logger.Error(schemas.DBErr, err)
					schemas.Fail(c, schemas.DBErr)
				} else {
					data := make(map[string]string)
					data["work_id"] = work.UUID
					schemas.SuccessCreate(c, data)
				}
			} else {
				schemas.SuccessCreate(c, nil)
			}
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}

// WorkDeleteApi 删除总任务
func WorkDeleteApi(c *gin.Context) {
	var schema = new(schemas.WorkDeleteSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.DeleteWorkByWorkUUID(schema.WorkUUID)
		if err != nil {
			schemas.Fail(c, err.Error())
		} else {
			schemas.SuccessDelete(c, schemas.CurdStatusOkMsg)
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}

// WorkGetInfoApi 获取总任务数据
func WorkGetInfoApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		work, err := models.GetWorkByUUID(schema.WorkUUID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			schemas.Fail(c, schemas.RecordNotFoundErr)
		} else {
			schemas.SuccessGet(c, work)
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}

// WorkPauseApi 暂停总任务
func WorkPauseApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.UpdateWorkByWorkUUID(schema.WorkUUID, "status", global.WorkStatusPause)
		if err != nil {
			schemas.Fail(c, schemas.RecordUpdateErr)
		} else {
			services.PauseWork(schema.WorkUUID)
			schemas.SuccessUpdate(c, nil)
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}

// WorkStopApi 停止总任务
func WorkStopApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.UpdateWorkByWorkUUID(schema.WorkUUID, "status", global.WorkStatusStop)
		if err != nil {
			schemas.Fail(c, schemas.RecordUpdateErr)
		} else {
			services.PauseWork(schema.WorkUUID)
			schemas.SuccessUpdate(c, nil)
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}

// WorkRestartApi 重启总任务
func WorkRestartApi(c *gin.Context) {
	var schema = new(schemas.WorkRestartSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.UpdateWorkByWorkUUID(schema.WorkUUID, "status", global.WorkStatusRestart)
		if err != nil {
			schemas.Fail(c, schemas.RecordUpdateErr)
		} else {
			schemas.SuccessUpdate(c, nil)
		}
	} else {
		schemas.Fail(c, err.Error())
	}
}
