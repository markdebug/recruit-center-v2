package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"org.thinkinai.com/recruit-center/api/handler"
	"org.thinkinai.com/recruit-center/pkg/config"
	"org.thinkinai.com/recruit-center/pkg/middleware"
)

// APIVersion 当前API版本
const APIVersion = "v1"

// SetupRouter 初始化路由配置
func SetupRouter(jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler, notificationHandler *handler.NotificationHandler, jobStatsHandler *handler.JobStatisticsHandler, jobFavoriteHandler *handler.JobFavoriteHandler) *gin.Engine {
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	// 配置全局中间件
	setupGlobalMiddleware(r)

	// 配置API路由
	apiGroup := r.Group(fmt.Sprintf("/api/%s", APIVersion))
	setupAPIRoutes(apiGroup, jobHandler, jobApplyHandler, resumeHandler, notificationHandler, jobStatsHandler, jobFavoriteHandler)

	// 配置工具路由
	setupToolRoutes(r)

	return r
}

// setupGlobalMiddleware 配置全局中间件
func setupGlobalMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// setupAPIRoutes 配置API路由
func setupAPIRoutes(api *gin.RouterGroup, jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler, notificationHandler *handler.NotificationHandler, jobStatsHandler *handler.JobStatisticsHandler, jobFavoriteHandler *handler.JobFavoriteHandler) {
	// 职位相关路由
	setupJobRoutes(api.Group("/jobs"), jobHandler, jobStatsHandler, jobFavoriteHandler)

	// 申请相关路由
	setupApplyRoutes(api.Group("/applies"), jobApplyHandler)

	// 简历相关路由
	setupResumeRoutes(api.Group("/resumes"), resumeHandler)
	// 通知相关路由
	setupNotificationsRouter(api.Group("/notifications"), notificationHandler)
}

// setupJobRoutes 配置职位相关路由
func setupJobRoutes(jobs *gin.RouterGroup, handler *handler.JobHandler, jobStatsHandler *handler.JobStatisticsHandler, jobFavoriteHandler *handler.JobFavoriteHandler) {
	jobs.POST("/", middleware.AuthRequired(), handler.Create)
	jobs.PUT("/:id", middleware.AuthRequired(), handler.Update)
	jobs.DELETE("/:id", middleware.AuthRequired(), handler.Delete)
	jobs.GET("/:id", handler.GetByID)
	jobs.GET("/", handler.List)

	// 职位统计相关路由
	jobs.GET("/jobs/:jobId/statistics", jobStatsHandler.GetJobStats)
	jobs.GET("/companies/:companyId/statistics", jobStatsHandler.GetCompanyStats)
	// 更新职位状态
	jobs.PUT("/:id/status", middleware.AuthRequired(), handler.UpdateStatus)
	// 根据公司搜索职位信息
	jobs.GET("/companies/:companyId/search", handler.SearchByCompany) // 假设有搜索功能

	// 收藏相关路由
	jobs.POST("/favorite/:jobId", middleware.AuthRequired(), jobFavoriteHandler.AddFavorite)
	jobs.DELETE("/favorite/:jobId", middleware.AuthRequired(), jobFavoriteHandler.RemoveFavorite)
	//获取用户收藏的职位
	jobs.GET("/favorites", middleware.AuthRequired(), jobFavoriteHandler.ListFavorites)
	//获取用户收藏职位的统计信息
	jobs.GET("/favorites/stats", middleware.AuthRequired(), jobFavoriteHandler.GetUserStatistics)
}

// setupApplyRoutes 配置申请相关路由
func setupApplyRoutes(applies *gin.RouterGroup, handler *handler.JobApplyHandler) {
	applies.POST("/", middleware.AuthRequired(), handler.Create)
	applies.GET("/my", middleware.AuthRequired(), handler.ListByUser)
	applies.GET("/job/:id", handler.List)
	applies.GET("/:id", handler.GetByID)
	applies.DELETE("/:id", middleware.AuthRequired(), handler.Delete)
	//根据公司id查询职位申请信息
	applies.GET("/company/:id", middleware.AuthRequired(), handler.ListByCompany)
	applies.PUT("/:id/status", middleware.AuthRequired(), handler.UpdateStatus)
}

// setupResumeRoutes 配置简历相关路由
func setupResumeRoutes(resumes *gin.RouterGroup, handler *handler.ResumeHandler) {
	resumes.POST("/", middleware.AuthRequired(), handler.Create)
	resumes.PUT("/:id", middleware.AuthRequired(), handler.Update)
	resumes.GET("/my", middleware.AuthRequired(), handler.GetByUser)
	// 添加文件上传路由
	resumes.POST("/upload",
		middleware.AuthRequired(),
		middleware.FileUploadValidator(middleware.FileUploadConfig(config.GetConfig().FileUploadConfig)),
		handler.UploadResume,
	)
	resumes.PUT("/access-status", middleware.AuthRequired(), handler.UpdateAccessStatus)
	resumes.PUT("/working-status", middleware.AuthRequired(), handler.UpdateWorkingStatus)
	resumes.GET("/:id", handler.GetByID)
	//格局share token获取简历
	resumes.GET("/share/:token", handler.GetByShareToken)
	//查看简历收藏相关信息
	resumes.GET("/:id/view", middleware.AuthRequired(), handler.ViewResume)
	//切换简历收藏状态
	resumes.PUT("/:id/favorite", middleware.AuthRequired(), handler.ToggleFavorite)
	//获取简历的统计信息
	resumes.GET("/:id/stats", middleware.AuthRequired(), handler.GetStats)
}

// SetupNotificationsRouter 通知相关路由配置
func setupNotificationsRouter(notifications *gin.RouterGroup, notificationHandler *handler.NotificationHandler) {

	notifications.GET("", notificationHandler.List)                        // 获取通知列表
	notifications.GET("/unread/count", notificationHandler.GetUnreadCount) // 获取未读数量
	notifications.POST("/read", notificationHandler.MarkAsRead)            // 标记已读
	notifications.POST("/send", notificationHandler.Send)                  // 发送通知

}

// setupToolRoutes 配置工具相关路由
func setupToolRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 创建 swagger 路由组并添加 CORS 中间件
	swagger := r.Group("/swagger")
	swagger.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 从配置获取服务地址和端口
	cfg := config.GetConfig()
	swaggerURL := fmt.Sprintf("http://%s:%d/swagger/doc.json",
		cfg.Host,
		cfg.Port,
	)

	// Swagger配置
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL(swaggerURL),
		ginSwagger.DefaultModelsExpandDepth(-1),
		// ginSwagger.InstanceName("default"),
		// ginSwagger.PersistAuthorization(true),
	))
}
