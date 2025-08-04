package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/pkg/middleware"
)

// SetupRouter 初始化路由配置
func SetupRouter(jobHandler *JobHandler, jobApplyHandler *JobApplyHandler) *gin.Engine {
	// 创建路由引擎
	r := gin.Default()

	// panic恢复
	r.Use(gin.Recovery())
	// 日志中间件
	r.Use(middleware.LoggerMiddleware())
	// CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API 版本分组
	v1 := r.Group("/api/v1")
	{
		// 职位相关路由
		jobs := v1.Group("/jobs")
		{
			jobs.POST("/", middleware.AuthRequired(), jobHandler.Create)      // 创建职位
			jobs.PUT("/:id", middleware.AuthRequired(), jobHandler.Update)    // 更新职位
			jobs.DELETE("/:id", middleware.AuthRequired(), jobHandler.Delete) // 删除职位
			jobs.GET("/:id", jobHandler.GetByID)                              // 获取职位详情
			jobs.GET("/", jobHandler.List)                                    // 获取职位列表
			jobs.GET("/search", jobHandler.Search)                            // 搜索职位
			// jobs.PUT("/:id/status", middleware.AuthRequired(), jobHandler.UpdateStatus) // 更新状态
		}

		// 职位申请相关路由
		applies := v1.Group("/applies")
		{
			applies.POST("/", middleware.AuthRequired(), jobApplyHandler.Create)                // 创建申请
			applies.GET("/my", middleware.AuthRequired(), jobApplyHandler.ListByUser)           // 我的申请
			applies.GET("/job/:id", jobApplyHandler.ListByJob)                                  // 职位申请列表
			applies.PUT("/:id/status", middleware.AuthRequired(), jobApplyHandler.UpdateStatus) // 更新状态
		}
	}

	return r
}
