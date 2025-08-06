package main

// @title           招聘中心 API
// @version         1.0
// @description     招聘中心系统 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
import (
	"context"
	"flag"
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

var (
	configPath string
	env        string
)

func init() {
	// 添加命令行参数支持
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.StringVar(&env, "env", "dev", "运行环境(dev/prod)")
	flag.Parse()
}

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
	// 确定配置文件路径
	if configPath == "" {
		// 优先使用环境变量
		if envPath := os.Getenv("CONFIG_FILE"); envPath != "" {
			configPath = envPath
		} else {
			// 默认使用环境对应的配置文件
			configPath = fmt.Sprintf("config/config.%s.yaml", env)
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	a.cfg = cfg
	return nil
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return err
	}
	logger.L.Info("应用启动",
		zap.String("env", a.cfg.Env),
		zap.String("version", a.cfg.Version),
		zap.String("config", configPath))
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
	resumeInteractionDao := dao.NewResumeInteractionDAO(db)

	// 初始化 Service 层
	jobService := service.NewJobService(jobDao)
	jobApplyService := service.NewJobApplyService(jobApplyDao, jobService)
	resumeService := service.NewResumeService(resumeDao)
	resumeInteractionService := service.NewResumeInteractionService(resumeInteractionDao)

	// 初始化 Handler 层
	return &Handlers{
		job:      handler.NewJobHandler(jobService),
		jobApply: handler.NewJobApplyHandler(jobApplyService),
		resume:   handler.NewResumeHandler(resumeService, resumeInteractionService),
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
