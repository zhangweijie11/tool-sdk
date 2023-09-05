package initizlize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"strings"
)

func InitCache(cfg *config.CacheConfig) error {
	if cfg.Sentinel {
		global.Cache = sentinelCache(*cfg)
	} else {
		global.Cache = redisCache(*cfg)
	}
	return global.Cache.Ping(context.Background()).Err()
}

func sentinelCache(cfg config.CacheConfig) *redis.Client {
	option := redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: strings.Split(cfg.Hosts, ","),
		DB:            cfg.Database,
		PoolSize:      cfg.PoolSize,
		Password:      cfg.Password,
	}
	return redis.NewFailoverClient(&option)
}

func redisCache(cfg config.CacheConfig) *redis.Client {
	option := redis.Options{
		Addr:     cfg.Hosts,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
		Password: cfg.Password,
	}
	return redis.NewClient(&option)
}
