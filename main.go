package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/config"
	"org.thinkinai.com/recruit-center/pkg/database"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

func main() {
	// 1. 加载配置文件
	cfg, err := config.LoadConfig("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 2. 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	logger.L.Info("应用启动",
		zap.String("env", cfg.Env),
		zap.String("version", cfg.Version))

	// 3. 初始化数据库连接
	db, err := database.Init(&cfg.DB)
	if err != nil {
		logger.L.Fatal("初始化数据库失败", zap.Error(err))
	}

	// 4. 初始化依赖
	// 4.1 初始化 DAO 层
	jobDao := dao.NewJobDao(db)
	jobApplyDao := dao.NewJobApplyDao(db)

	// 4.2 初始化 Service 层
	jobService := service.NewJobService(jobDao)
	jobApplyService := service.NewJobApplyService(jobApplyDao, jobDao)

	// 4.3 初始化 API 层
	jobHandler := api.NewJobHandler(jobService)
	jobApplyHandler := api.NewJobApplyHandler(jobApplyService)

	// 5. 设置路由
	r := api.SetupRouter(jobHandler, jobApplyHandler)

	// 6. 启动HTTP服务
	// 8. 启动HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// 启动服务器的goroutine
	go func() {
		logger.L.Info("HTTP服务启动",
			zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("HTTP服务启动失败", zap.Error(err))
		}
	}()

	// 7. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L.Info("开始关闭服务...")

	// 给服务器30秒时间来完成当前请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.L.Fatal("服务关闭失败", zap.Error(err))
	}

	logger.L.Info("服务已关闭")
}
