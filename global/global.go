package global

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gorm.io/gorm"
	"sync"
)

// 任务相关
const (
	TimeFormatDay             = "2006-01-02"          // 固定format时间，2006-12345
	TimeFormatSecond          = "2006-01-02 15:04:05" // 固定format时间，2006-12345
	WorkStatusPending         = "pending"
	WorkStatusDoingGo         = "doingGo"
	WorkStatusDoing           = "doing"
	WorkStatusDone            = "done"
	WorkStatusFailed          = "failed"
	WorkStatusPause           = "pause"
	WorkStatusStop            = "stop"
	WorkStatusRestart         = "restart"
	CallbackWorkStatusSuccess = "success"
	CallbackWorkStatusFailed  = "failed"
	CallbackTypeApi           = "API"
	CallbackTypeMQ            = "MQ"
	CallbackTypegRPC          = "gRPC"
)

var (
	Config        *config.Cfg
	ElasticClient *elastic.Client
	Db            *gorm.DB
	Cache         *redis.Client
)

var (
	ValidExecutorIns  ExecutorInterface
	ValidExecutorChan ExecutorChan
	ValidDoingWork    DoingWork
	ValidModels       []interface{}
	ValidRouter       []func(*gin.Engine) gin.IRoutes
	ValidProgressChan chan *Progress
	ValidResultChan   chan *Result
)

type Work struct {
	WorkUUID string
	Context  context.Context
	Cancel   context.CancelFunc
}

type Progress struct {
	WorkUUID     string
	Progress     float64
	ProgressType string
	ProgressUrl  string
}

type Result struct {
	WorkUUID     string
	TaskUUID     string
	Result       map[string]interface{}
	CallbackType string
	CallbackUrl  string
}

type DoingWork struct {
	sync.Mutex
	DoingWorkMap map[string]*Work
}

type ExecutorChan struct {
	sync.Mutex
	WorkExecute chan bool
}

func init() {
	doingWorkMap := make(map[string]*Work)
	ValidDoingWork = DoingWork{DoingWorkMap: doingWorkMap}
	ValidProgressChan = make(chan *Progress, 1000)
	ValidResultChan = make(chan *Result, 1000)
}
