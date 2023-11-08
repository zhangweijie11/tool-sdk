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
	"time"
)

// 任务相关
const (
	TimeFormatDay    = "2006-01-02"          // 固定format时间，2006-12345
	TimeFormatSecond = "2006-01-02 15:04:05" // 固定format时间，2006-12345
)

const (
	callbackTypeApi  = "API"
	callbackTypeMQ   = "MQ"
	callbackTypegRPC = "gRPC"
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

// 通用参数错误
var (
	JsonParseErr          = "json解析失败"
	ParameterErr          = "参数错误"
	ParameterTypeErr      = "参数类型错误"
	ParamsLengthErr       = "参数长度错误"
	UnSupportOperationErr = "不支持的操作"
	InternalErr           = "内部错误"
	DecryptConfigErr      = "解密配置错误"
	TimeCreateErr         = "无效的创建时间区间"
	TimeUpdateErr         = "无效的更新时间区间"
)

// 数据库错误
var (
	DBErr             = "数据库错误"
	RecordNotFoundErr = "数据不存在"
	RecordExistsErr   = "数据已存在"
	RecordDeleteErr   = "数据删除错误"
	RecordUpdateErr   = "数据更新错误"
)

// 任务相关错误
var (
	WorkTargetErr       = "任务目标错误"
	WorkGetErr          = "查询任务错误"
	WorkUpdateErr       = "更新任务错误"
	WorkExecuteErr      = "执行任务错误"
	WorkCancelErr       = "停止任务"
	WorkTimeoutErr      = "任务超时错误"
	WorkProgressErr     = "任务进度错误"
	WorkCallbackTypeErr = "回调/推送类型错误"
	WorkCallbackErr     = "回调任务错误"
	WorkCallbackUrlErr  = "回调地址错误"
)

// 消息队列错误
var (
	MQConnectErr = "MQ连接错误"
	MQChannelErr = "MQ通道创建错误"
	MQQueueErr   = "MQ队列声明错误"
	MQMessageErr = "MQ消息发布错误"
)

// gRPC错误
var (
	RPCConnectErr = "RPC 连接错误"
	RPCPushErr    = "RPC 发送错误"
)

// API 错误
var (
	APIConnectErr  = "API连接错误"
	APIRequestErr  = "API请求错误"
	APIResponseErr = "API响应错误"
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
		return errors.New(ParameterTypeErr)
	}

	// 获取切片的长度
	sliceLength := dataValue.Len()

	// 检查切片长度是否超过限定值
	if sliceLength >= max || sliceLength <= min {
		return errors.New(ParamsLengthErr)
	}

	return nil
}

func ValidateCallbackUrlAndType(callbackType, callbackUrl string) error {
	switch strings.ToLower(callbackType) {
	case strings.ToLower(callbackTypeApi):
		if strings.ToLower(callbackUrl) == "" {
			return errors.New(WorkCallbackUrlErr)
		}
	case strings.ToLower(callbackTypeMQ):
		if len(strings.Split(strings.ToLower(callbackUrl), ",")) != 3 {
			return errors.New(WorkCallbackUrlErr)
		}
	case strings.ToLower(callbackTypegRPC):
		return nil
	default:
		return errors.New(WorkCallbackTypeErr)
	}

	return nil
}

// TimeRangeValidator 时间范围验证
func TimeRangeValidator(times []string) (ts [2]time.Time, err error) {
	times[0] = times[0] + " 00:00:00"
	times[1] = times[1] + " 23:59:59"
	if ts[0], err = time.Parse(TimeFormatSecond, times[0]); err != nil {
		return ts, err
	}
	if ts[1], err = time.Parse(TimeFormatSecond, times[1]); err != nil {
		return ts, err
	}
	if ts[0].Before(ts[1]) {
		return ts, nil
	}
	return ts, err
}
