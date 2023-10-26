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
			return
		}
		if schema.CallbackType != "" {
			err = schemas.ValidateCallbackUrlAndType(schema.CallbackType, schema.CallbackUrl)
			if err != nil {
				schemas.Fail(c, err.Error())
				return
			}
		}
		if schema.ProgressType != "" {
			err = schemas.ValidateCallbackUrlAndType(schema.ProgressType, schema.ProgressUrl)
			if err != nil {
				schemas.Fail(c, err.Error())
				return
			}
		}

		if global.Config.Database.Activate == true || global.Config.Database.Activate == false {
			jsonBytes, err := json.Marshal(schema.Params)
			if err != nil {
				logger.Error(schemas.JsonParseErr, err)
				schemas.Fail(c, schemas.JsonParseErr)
				return
			}

			if schema.WorkUUID == "" {
				schema.WorkUUID = uuid.New().String()
			}

			dbWork, err := models.GetWorkByUUID(schema.WorkUUID)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				work := &models.Work{
					UUID:           schema.WorkUUID,
					Params:         jsonBytes,
					Status:         global.WorkStatusPending,
					Source:         schema.Source,
					Priority:       schema.Priority,
					Retry:          3,
					CallbackUrl:    schema.CallbackUrl,
					CallbackType:   schema.CallbackType,
					CallbackStatus: global.WorkStatusPending,
					WorkType:       global.Config.Server.ServerName,
					ProgressType:   schema.ProgressType,
					ProgressUrl:    schema.ProgressUrl,
				}

				err = models.CreateWok(work)

				if err != nil {
					logger.Error(schemas.DBErr, err)
					schemas.Fail(c, schemas.DBErr)
					return
				} else {
					data := make(map[string]string)
					data["work_uuid"] = work.UUID
					schemas.SuccessCreate(c, data)
					return
				}
			} else {
				if dbWork.ID > 0 {
					schemas.Fail(c, schemas.RecordExistsErr)
					return
				}
			}

		} else {
			schemas.SuccessCreate(c, nil)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkDeleteApi 删除总任务
func WorkDeleteApi(c *gin.Context) {
	var schema = new(schemas.WorkDeleteSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.WorkDeleteByWorkUUID(schema.WorkUUID)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			schemas.SuccessDelete(c, schemas.CurdStatusOkMsg)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkGetInfoApi 获取总任务数据
func WorkGetInfoApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		work, err := models.GetWorkByUUID(schema.WorkUUID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			schemas.Fail(c, schemas.RecordNotFoundErr)
			return
		} else {
			schemas.SuccessGet(c, work)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkPauseApi 暂停总任务
func WorkPauseApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.WorkUpdateByWorkUUID(schema.WorkUUID, "status", global.WorkStatusPause)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			services.PauseWork(schema.WorkUUID)
			schemas.SuccessUpdate(c, nil)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkStopApi 停止总任务
func WorkStopApi(c *gin.Context) {
	var schema = new(schemas.WorkGetInfoSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.WorkUpdateByWorkUUID(schema.WorkUUID, "status", global.WorkStatusStop)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			services.PauseWork(schema.WorkUUID)
			schemas.SuccessUpdate(c, nil)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkRestartApi 重启总任务
func WorkRestartApi(c *gin.Context) {
	var schema = new(schemas.WorkRestartSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = controller.WorkUpdateByWorkUUID(schema.WorkUUID, "status", global.WorkStatusRestart)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			schemas.SuccessUpdate(c, nil)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkListApi 获取任务结果
func WorkListApi(c *gin.Context) {
	var schema = new(schemas.WorkListSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		data, err := controller.WorkList(schema)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			schemas.SuccessGet(c, data)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}

// WorkResultListApi 获取任务结果
func WorkResultListApi(c *gin.Context) {
	var schema = new(schemas.WorkResultListSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		data, err := controller.WorkResultList(schema)
		if err != nil {
			schemas.Fail(c, err.Error())
			return
		} else {
			schemas.SuccessGet(c, data)
			return
		}
	} else {
		schemas.Fail(c, err.Error())
		return
	}
}
