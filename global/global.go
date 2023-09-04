package global

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"tool-sdk/config"
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
)
