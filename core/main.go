package core

import (
	"github.com/gin-gonic/gin"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/initizlize"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gitlab.example.com/zhangweijie/tool-sdk/models"
	"gitlab.example.com/zhangweijie/tool-sdk/services"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func Start() error {
	err := logger.Setup(global.Config.Server.Loglevel)
	if err != nil {
		logger.Panic("设置 Logger 异常", err)
	}

	if err = initizlize.InitWorker(global.Config.Server.Concurrency); err != nil {
		logger.Panic("初始化任务执行者异常", err)
	}

	if global.Config.Elastic.Activate {
		if err = initizlize.InitElastic(&global.Config.Elastic); err != nil {
			logger.Panic("ElasticSearch 连接异常", err)
		}
	}

	if global.Config.Database.Activate == false || global.Config.Database.Activate == true {
		if err = initizlize.InitDatabase(&global.Config.Database); err != nil {
			logger.Panic("数据库连接异常", err)
		}
		if err = models.UpdateWorkDoingToPending(); err != nil {
			logger.Panic("任务状态变更错误", err)
		}
	}

	if global.Config.Cache.Activate {
		if err = initizlize.InitCache(&global.Config.Cache); err != nil {
			logger.Panic("缓存连接异常", err)
		}
	}
	// 开启 pprof 性能分析
	//go func() {
	//	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪，block
	//	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪，mutex
	//	log.Println(http.ListenAndServe("localhost:8080", nil))
	//}()

	go services.LoopExecuteWork()
	go services.LoopProgressResult()

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
	for _, f := range global.ValidRouter {
		f(engine)
	}
	server := &http.Server{
		Addr:           ":" + global.Config.Server.RunPort,                             // 监听地址
		MaxHeaderBytes: 1 << 20,                                                        // 1048576
		Handler:        engine,                                                         // 服务引擎
		ReadTimeout:    time.Duration(global.Config.Server.ReadTimeout) * time.Second,  // 请求超市
		WriteTimeout:   time.Duration(global.Config.Server.WriteTimeout) * time.Second, // 响应超时
	}

	return server.ListenAndServe()
}
