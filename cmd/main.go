package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/initizlize"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/routers"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func Start() error {
	err := logger.Setup(global.Config.Server.Loglevel)
	if err != nil {
		logger.Panic("设置 Logger 异常", err)
	}

	if global.Config.Elastic.Activate {
		if err = initizlize.InitElastic(&global.Config.Elastic); err != nil {
			logger.Panic("ElasticSearch 连接异常", err)
		}
	}

	if global.Config.Database.Activate {
		if err = initizlize.InitDatabase(&global.Config.Database); err != nil {
			logger.Panic("数据库连接异常", err)
		}
	}

	if global.Config.Cache.Activate {
		if err = initizlize.InitCache(&global.Config.Cache); err != nil {
			logger.Panic("缓存连接异常", err)
		}
	}

	// 开启 pprof 性能分析
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8080", nil))
	//}()

	var engine = gin.New()
	switch global.Config.Server.RunMode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	//引用中间件
	engine.Use(schemas.Cors())
	engine.Use(logger.GinLogger())
	engine.Use(logger.GinRecovery(true))

	// 初始化路由
	routers.InitPingRouter(engine)
	routers.InitWorkRouter(engine)
	server := &http.Server{
		Addr:           ":" + global.Config.Server.RunPort,                             // 监听地址
		MaxHeaderBytes: 1 << 20,                                                        // 1048576
		Handler:        engine,                                                         // 服务引擎
		ReadTimeout:    time.Duration(global.Config.Server.ReadTimeout) * time.Second,  // 请求超市
		WriteTimeout:   time.Duration(global.Config.Server.WriteTimeout) * time.Second, // 响应超时
	}

	return server.ListenAndServe()
}
