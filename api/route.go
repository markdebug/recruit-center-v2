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
func SetupRouter(jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler) *gin.Engine {
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	// 配置全局中间件
	setupGlobalMiddleware(r)

	// 配置API路由
	apiGroup := r.Group(fmt.Sprintf("/api/%s", APIVersion))
	setupAPIRoutes(apiGroup, jobHandler, jobApplyHandler, resumeHandler)

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
func setupAPIRoutes(api *gin.RouterGroup, jobHandler *handler.JobHandler, jobApplyHandler *handler.JobApplyHandler, resumeHandler *handler.ResumeHandler) {
	// 职位相关路由
	setupJobRoutes(api.Group("/jobs"), jobHandler)

	// 申请相关路由
	setupApplyRoutes(api.Group("/applies"), jobApplyHandler)

	// 简历相关路由
	setupResumeRoutes(api.Group("/resumes"), resumeHandler)
}

// setupJobRoutes 配置职位相关路由
func setupJobRoutes(jobs *gin.RouterGroup, handler *handler.JobHandler) {
	jobs.POST("/", middleware.AuthRequired(), handler.Create)
	jobs.PUT("/:id", middleware.AuthRequired(), handler.Update)
	jobs.DELETE("/:id", middleware.AuthRequired(), handler.Delete)
	jobs.GET("/:id", handler.GetByID)
	jobs.GET("/", handler.List)
	// jobs.GET("/search", handler.Search)
}

// setupApplyRoutes 配置申请相关路由
func setupApplyRoutes(applies *gin.RouterGroup, handler *handler.JobApplyHandler) {
	applies.POST("/", middleware.AuthRequired(), handler.Create)
	// applies.GET("/my", middleware.AuthRequired(), handler.ListByUser)
	// applies.GET("/job/:id", handler.ListByJob)
	// applies.PUT("/:id/status", middleware.AuthRequired(), handler.UpdateStatus)
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
