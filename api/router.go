package api

import (
	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/internal/handler"
)

// SetupRouter 初始化路由
func SetupRouter(
	jobHandler *handler.JobHandler,
	jobApplyHandler *handler.JobApplyHandler,
	jobStatsHandler *handler.JobStatisticsHandler,
	// ...other handlers
) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		// 职位相关路由
		v1.GET("/jobs", jobHandler.List)
		v1.POST("/jobs", jobHandler.Create)
		v1.GET("/jobs/:id", jobHandler.GetByID)
		v1.PUT("/jobs/:id", jobHandler.Update)
		v1.DELETE("/jobs/:id", jobHandler.Delete)

		// 职位申请相关路由
		v1.GET("/job-applies", jobApplyHandler.List)
		v1.POST("/job-applies", jobApplyHandler.Create)
		v1.GET("/job-applies/:id", jobApplyHandler.GetByID)
		v1.PUT("/job-applies/:id/status", jobApplyHandler.UpdateStatus)
		v1.DELETE("/job-applies/:id", jobApplyHandler.Delete)

		// 职位统计相关路由
		v1.GET("/jobs/:jobId/statistics", jobStatsHandler.GetJobStats)
		v1.GET("/companies/:companyId/statistics", jobStatsHandler.GetCompanyStats)
	}

	return router
}
