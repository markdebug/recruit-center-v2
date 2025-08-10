package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type JobStatisticsHandler struct {
	statsService *service.JobStatisticsService
}

func NewJobStatisticsHandler(service *service.JobStatisticsService) *JobStatisticsHandler {
	return &JobStatisticsHandler{statsService: service}
}

// GetJobStats 获取职位统计信息
// @Summary 获取职位统计信息
// @Description 获取指定职位的统计信息
// @Tags 职位统计
// @Accept json
// @Produce json
// @Param jobId path int true "职位ID"
// @Success 200 {object} response.Response{data=response.JobStatisticsResponse}
// @Router /api/v1/jobs/{jobId}/statistics [get]
func (h *JobStatisticsHandler) GetJobStats(c *gin.Context) {
	jobID, err := strconv.ParseUint(c.Param("jobId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	stats, err := h.statsService.GetJobStatisticsByJobID(uint(jobID))
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(stats))
}

// GetCompanyStats 获取公司职位统计信息
// @Summary 获取公司职位统计信息
// @Description 获取指定公司所有职位的统计信息
// @Tags 职位统计
// @Accept json
// @Produce json
// @Param companyId path int true "公司ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.JobStatisticsListResponse}
// @Router /api/v1/companies/{companyId}/statistics [get]
func (h *JobStatisticsHandler) GetCompanyStats(c *gin.Context) {
	companyID, _ := strconv.ParseUint(c.Param("companyId"), 10, 32)
	page, size := parsePageSize(c)

	stats, err := h.statsService.GetCompanyStats(uint(companyID), page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}
	c.JSON(http.StatusOK, response.NewSuccess(stats))
}
