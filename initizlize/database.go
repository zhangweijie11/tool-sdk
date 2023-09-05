package initizlize

import (
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) (err error) {
	var mysqlCfg = mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName),
	}

	var gormCfg = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.GormLoggerNew(logger.Get(), &logger.GormLoggerConfig{
			LogLevel: cfg.LogLevel, SlowThreshold: time.Duration(cfg.SlowThreshold) * time.Millisecond}),
	}

	global.Db, err = gorm.Open(mysql.New(mysqlCfg), gormCfg)
	if err != nil {
		return
	}

	////自动生成迁移脚本
	//availableModels := []interface{}{}
	//err = global.Db.AutoMigrate(availableModels...)
	//if err != nil {
	//	logger.Panic("数据库连接失败", err)
	//}

	db, _ := global.Db.DB()
	//设置数据库最大空闲连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	//设置数据库最大打开的连接数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	//设置数据库连接的最大空闲时间，过期的连接可能会在重用之前惰性关闭。
	db.SetConnMaxIdleTime(60 * time.Second)
	//设置数据库一个连接可能被重用的最大时间，过期的连接可能会在重用之前惰性关闭。
	db.SetConnMaxLifetime(-1)
	return err
}
