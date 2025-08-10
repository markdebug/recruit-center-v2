package api

import (
	"fmt"
	"net/http"
	"time"

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
func SetupRouter(jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler, notificationHandler *handler.NotificationHandler, jobStatsHandler *handler.JobStatisticsHandler) *gin.Engine {
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	// 配置全局中间件
	setupGlobalMiddleware(r)

	// 配置API路由
	apiGroup := r.Group(fmt.Sprintf("/api/%s", APIVersion))
	setupAPIRoutes(apiGroup, jobHandler, jobApplyHandler, resumeHandler, notificationHandler, jobStatsHandler)

	// 配置工具路由
	setupToolRoutes(r)

	return r
}

// setupGlobalMiddleware 配置全局中间件
func setupGlobalMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
}

// setupAPIRoutes 配置API路由
func setupAPIRoutes(api *gin.RouterGroup, jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler, notificationHandler *handler.NotificationHandler, jobStatsHandler *handler.JobStatisticsHandler) {
	// 职位相关路由
	setupJobRoutes(api.Group("/jobs"), jobHandler, jobStatsHandler)

	// 申请相关路由
	setupApplyRoutes(api.Group("/applies"), jobApplyHandler)

	// 简历相关路由
	setupResumeRoutes(api.Group("/resumes"), resumeHandler)
	// 通知相关路由
	setUpNotificationsRouter(api.Group("/notifications"), notificationHandler)
}

// setupJobRoutes 配置职位相关路由
func setupJobRoutes(jobs *gin.RouterGroup, handler *handler.JobHandler, jobStatsHandler *handler.JobStatisticsHandler) {
	jobs.POST("/", middleware.AuthRequired(), handler.Create)
	jobs.PUT("/:id", middleware.AuthRequired(), handler.Update)
	jobs.DELETE("/:id", middleware.AuthRequired(), handler.Delete)
	jobs.GET("/:id", handler.GetByID)
	jobs.GET("/", handler.List)

	// 职位统计相关路由
	jobs.GET("/jobs/:jobId/statistics", jobStatsHandler.GetJobStats)
	jobs.GET("/companies/:companyId/statistics", jobStatsHandler.GetCompanyStats)
	// jobs.GET("/search", handler.Search)
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
func setUpNotificationsRouter(notifications *gin.RouterGroup, notificationHandler *handler.NotificationHandler) {

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

	// Swagger文档
	// Swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
}
