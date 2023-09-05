package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
)

// CreateWorkApi 创建总任务
func CreateWorkApi(c *gin.Context) {
	var schema = new(schemas.WorkCreateSchema)
	if err := schemas.BindSchema(c, schema, binding.JSON); err == nil {
		err = global.ParamsIns.ValidWorkCreateParams(schema.Params)
		if err != nil {
			schemas.Fail(c, schemas.ParameterErr)
		}

		if global.Config.Database.Activate {
			jsonBytes, err := json.Marshal(schema)
			if err != nil {
				logger.Warn(schemas.JsonParseErr)
			}

			work := &models.Work{
				UUID:   schema.WorkID,
				Params: jsonBytes,
				Status: global.WorkStatusPending,
				Source: schema.Source,
			}

			err = models.CreateWok(work)

			if err != nil {
				logger.Error(schemas.DBErr, err)
				schemas.Fail(c, schemas.DBErr)
			}
		}

		schemas.SuccessCreate(c, "Success")
	}
}
