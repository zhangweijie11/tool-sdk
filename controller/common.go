package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"strings"
)

var taskValidatorErrorMessage = map[string]string{}

// registerValidatorRule 注册参数验证错误消息, Key = e.StructNamespace(), value.key = e.Field()+e.Tag()
var registerValidatorRule = map[string]map[string]string{}

var (
	JsonParseErr          = "json解析失败"
	ParameterErr          = "参数错误"
	DBErr                 = "数据库错误"
	NoTaskErr             = "没有相关任务"
	UnSupportOperationErr = "不支持的操作"
	ToolStatusErr         = "工具状态不允许操作"
	InternalErr           = "内部错误"
)

// serializeValidatorError 参数tag验证失败转换
func serializeValidatorError(e validator.FieldError) (message string) {
	switch e.Field() {
	default:
		message = registerValidatorRule[strings.Split(e.StructNamespace(), ".")[0]][e.Field()+e.Tag()]
	case "Page", "Size":
		switch e.Tag() {
		case "min":
			message = e.Field() + "最小值为" + e.Param()
		case "max":
			message = e.Field() + "最大值为" + e.Param()
		}
	}
	return message
}

// serializeTypeError 参数类型错误转换
func serializeTypeError(e *json.UnmarshalTypeError) string {
	return fmt.Sprintf("参数 %s 类型错误, 预期 %s, 接收到 %s", e.Field, e.Type, e.Value)
}

// BindSchema 绑定单类型的请求参数
func BindSchema(c *gin.Context, obj interface{}, bind binding.Binding) (err error) {
	if err = c.ShouldBindWith(obj, bind); err != nil {
		var msg string
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			msg = serializeTypeError(e)
		case validator.ValidationErrors:
			msg = serializeValidatorError(e[0])
		default:
			msg = "无效的请求参数"
		}
		schemas.Fail(c, msg)
	}
	return err
}
