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
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/api"
	"org.thinkinai.com/recruit-center/api/handler"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/config"
	"org.thinkinai.com/recruit-center/pkg/database"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// App 应用程序结构体
type App struct {
	cfg    *config.Config
	server *http.Server
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{}
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	var err error

	// 1. 加载配置
	if err = a.loadConfig(); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 2. 初始化日志
	if err = a.initLogger(); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	// 3. 初始化HTTP服务器
	if err = a.initHTTPServer(); err != nil {
		return fmt.Errorf("初始化HTTP服务器失败: %w", err)
	}

	return nil
}

// loadConfig 加载配置
func (a *App) loadConfig() error {
	if err := config.InitGlobalConfig(""); err != nil {
		return err
	}
	a.cfg = config.GetConfig()
	return nil
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return err
	}
	logger.L.Info("应用启动",
		zap.String("env", a.cfg.Env),
		zap.String("version", a.cfg.Version))
	return nil
}

// initHTTPServer 初始化HTTP服务器
func (a *App) initHTTPServer() error {
	// 初始化数据库
	db, err := database.Init(&a.cfg.DB)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化依赖
	handlers, err := a.initializeDependencies(db)
	if err != nil {
		return fmt.Errorf("初始化依赖失败: %w", err)
	}

	// 设置路由
	router := api.SetupRouter(handlers.job, handlers.jobApply, handlers.resume)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

// Handlers 处理器集合
type Handlers struct {
	job      *handler.JobHandler
	jobApply *handler.JobApplyHandler
	resume   *handler.ResumeHandler
}

// initializeDependencies 初始化所有依赖
func (a *App) initializeDependencies(db *gorm.DB) (*Handlers, error) {
	// 初始化 DAO 层
	jobDao := dao.NewJobDAO(db)
	jobApplyDao := dao.NewJobApplyDAO(db)
	resumeDao := dao.NewResumeDAO(db)

	// 初始化 Service 层
	jobService := service.NewJobService(jobDao)
	jobApplyService := service.NewJobApplyService(jobApplyDao, jobService)
	resumeService := service.NewResumeService(resumeDao)

	// 初始化 Handler 层
	return &Handlers{
		job:      handler.NewJobHandler(jobService),
		jobApply: handler.NewJobApplyHandler(jobApplyService),
		resume:   handler.NewResumeHandler(resumeService),
	}, nil
}

// Run 运行应用
func (a *App) Run() error {
	// 启动HTTP服务器
	go func() {
		logger.L.Info("HTTP服务启动", zap.String("address", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("HTTP服务启动失败", zap.Error(err))
		}
	}()

	return a.waitForShutdown()
}

// waitForShutdown 等待关闭信号
func (a *App) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L.Info("开始关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务关闭失败: %w", err)
	}

	logger.L.Info("服务已关闭")
	return nil
}

func main() {
	app := NewApp()
	if err := app.Initialize(); err != nil {
		fmt.Printf("应用初始化失败: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		logger.L.Fatal("应用运行失败", zap.Error(err))
	}
}
