package initizlize

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"tool-sdk/config"
	"tool-sdk/global"
	"tool-sdk/middleware/logger"
)

// InitElastic  初始化 Elasticsearch，并创建相对应的索引
func InitElastic(cfg *config.ElasticConfig) (err error) {
	//创建 client
	global.ElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		//由于部署为单节点模式，所以需要设置嗅探器和健康检查为 False，这两个机制都是在集群模式下起作用
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	if err != nil {
		logger.Panic("Elasticsearch 连接失败", err)
	}

	//执行 ES 请求需要提供上下文对象
	ctx := context.Background()

	//首先检测索引是否存在
	_, err = global.ElasticClient.IndexExists(cfg.Index).Do(ctx)
	if err != nil {
		logger.Panic(fmt.Sprintf("Elasticsearch 索引 %s 不存在", cfg.Index), err)
	}

	return err
}
