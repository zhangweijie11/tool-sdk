package global

import (
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gorm.io/gorm"
)

// 端口扫描相关
const (
	TimeFormatDay    = "2006-01-02"          // 固定format时间，2006-12345
	TimeFormatSecond = "2006-01-02 15:04:05" // 固定format时间，2006-12345
)

var (
	Config        *config.Cfg
	ElasticClient *elastic.Client
	Db            *gorm.DB
	Cache         *redis.Client
)

var ParamsIns paramsIns

// 任务相关
const (
	WorkStatusPending   = "pending"
	WorkStatusDoing     = "doing"
	WorkStatusDone      = "done"
	WorkStatusFailed    = "failed"
	WorkStatusPause     = "pause"
	WorkStatusCancelled = "cancelled"
)
