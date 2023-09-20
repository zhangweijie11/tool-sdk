package schemas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"reflect"
	"strings"
)

var taskValidatorErrorMessage = map[string]string{
	"WorkUUIDrequired": "缺少任务唯一标识",
	"Sourcerequired":   "缺少任务来源参数",
	"Prioritymax":      "无效的任务优先级（0-9）",
	"Prioritymin":      "无效的任务优先级（0-9）",
	"Priorityrequired": "缺少任务基础参数",
	"Paramsrequired":   "缺少任务基础参数",
}

// registerValidatorRule 注册参数验证错误消息, Key = e.StructNamespace(), value.key = e.Field()+e.Tag()
var registerValidatorRule = map[string]map[string]string{
	"WorkCreateSchema":  taskValidatorErrorMessage,
	"WorkDeleteSchema":  taskValidatorErrorMessage,
	"WorkGetInfoSchema": taskValidatorErrorMessage,
	"WorkPauseSchema":   taskValidatorErrorMessage,
	"WorkStopSchema":    taskValidatorErrorMessage,
	"WorkRestartSchema": taskValidatorErrorMessage,
}

var (
	JsonParseErr          = "json解析失败"
	ParameterErr          = "参数错误"
	LengthErr             = "参数长度错误"
	DBErr                 = "数据库错误"
	RecordNotFoundErr     = "数据不存在"
	RecordDeleteErr       = "数据删除错误"
	RecordUpdateErr       = "数据更新错误"
	GetWorkErr            = "查询任务错误"
	UpdateWorkErr         = "更新任务错误"
	ExecuteWorkErr        = "执行任务错误"
	CancelWorkErr         = "停止任务错误"
	TimeoutWorkErr        = "任务超时错误"
	ProgressErr           = "任务进度推送错误"
	TypeErr               = "回调/推送类型错误"
	CallbackWorkErr       = "回调任务错误"
	UnSupportOperationErr = "不支持的操作"
	InternalErr           = "内部错误"
	WorkTargetErr         = "任务目标错误"
)

// serializeValidatorError 参数tag验证失败转换
func serializeValidatorError(e validator.FieldError, validatorRule map[string]map[string]string) (message string) {
	if validatorRule == nil {
		validatorRule = registerValidatorRule
	}

	switch e.Field() {
	case "Page", "Size":
		switch e.Tag() {
		case "min":
			message = e.Field() + "最小值为" + e.Param()
		case "max":
			message = e.Field() + "最大值为" + e.Param()
		}
	default:
		message = validatorRule[strings.Split(e.StructNamespace(), ".")[0]][e.Field()+e.Tag()]
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
			msg = serializeValidatorError(e[0], nil)
		default:
			msg = ParameterErr
		}
		return errors.New(msg)
	}
	return err
}

func mapToReader(data map[string]interface{}) (io.Reader, error) {
	// 将 map 转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 创建一个 bytes.Buffer 来存储 JSON 数据
	buffer := bytes.NewBuffer(jsonData)

	// 返回 bytes.Buffer 作为 io.Reader
	return buffer, nil
}

// CustomBindSchema 绑定单类型的请求参数
func CustomBindSchema(params map[string]interface{}, obj interface{}, validatorRule map[string]map[string]string) (err error) {
	validParams, _ := mapToReader(params)
	if err = customBind(validParams, obj); err != nil {
		var msg string
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			msg = serializeTypeError(e)
		case validator.ValidationErrors:
			msg = serializeValidatorError(e[0], validatorRule)
		default:
			msg = ParameterErr
		}
		return errors.New(msg)
	}
	return err
}

func customBind(params io.Reader, obj any) error {
	if params == nil {
		return errors.New(ParameterErr)
	}
	return decodeJSON(params, obj)
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if binding.EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if binding.EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj any) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

func ValidateLength(params interface{}, min, max int) error {
	// 检查 data 是否是一个切片类型
	dataValue := reflect.ValueOf(params)
	if dataValue.Kind() != reflect.Slice {
		return errors.New(ParameterErr)
	}

	// 获取切片的长度
	sliceLength := dataValue.Len()

	// 检查切片长度是否超过限定值
	if sliceLength >= max || sliceLength <= min {
		return errors.New(LengthErr)
	}

	return nil
}
