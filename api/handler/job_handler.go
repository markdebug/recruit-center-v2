package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

// Create 创建职位
// @Summary 创建新职位
// @Description 创建一个新的职位信息
// @Tags 职位管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param job body request.CreateJobRequest true "职位信息"
// @Success 200 {object} response.Response{data=model.Job}
// @Failure 400 {object} response.Response
// @Router /api/v1/jobs [post]
func (h *JobHandler) Create(c *gin.Context) {
	var req request.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	companyID := c.GetUint("companyId")
	req.CompanyID = companyID
	job := req.ToModel()

	if err := h.jobService.Create(job); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(job))
}

// Update 更新职位
// @Summary 更新职位信息
// @Description 更新指定职位的信息
// @Tags 职位管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "职位ID"
// @Param job body request.UpdateJobRequest true "职位信息"
// @Success 200 {object} response.Response{data=model.Job}
// @Failure 400 {object} response.Response
// @Router /api/v1/jobs/{id} [put]
func (h *JobHandler) Update(c *gin.Context) {
	var req request.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.BadRequest))
		return
	}

	job := req.ToModel()
	if err := h.jobService.Update(job); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(job))
}

// Delete 删除职位
func (h *JobHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.jobService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetByID 获取职位详情
// @Summary 获取职位详情
// @Description 获取指定ID的职位详细信息
// @Tags 职位管理
// @Produce json
// @Param id path int true "职位ID"
// @Success 200 {object} response.Response{data=model.Job}
// @Failure 404 {object} response.Response
// @Router /api/v1/jobs/{id} [get]
func (h *JobHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	job, err := h.jobService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(job))
}

// List 获取职位列表
// @Summary 获取职位列表
// @Description 分页获取职位列表
// @Tags 职位管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.PageResponse{data=[]model.Job}
// @Router /api/v1/jobs [get]
func (h *JobHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	jobs, total, err := h.jobService.List(page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewPage(jobs, total, page, size))
}

// SearchByKeyword 关键词搜索职位
func (h *JobHandler) SearchByKeyword(c *gin.Context) {
	keyword := c.Query("keyword")

	jobs, err := h.jobService.SearchByKeyword(keyword)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(jobs))
}

// SearchByCondition 多条件搜索职位
func (h *JobHandler) SearchByCondition(c *gin.Context) {
	var conditions map[string]interface{}
	if err := c.ShouldBindJSON(&conditions); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.BadRequest))
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	result, err := h.jobService.SearchByCondition(conditions, page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}

// UpdateStatus 更新职位状态
func (h *JobHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	status, _ := strconv.Atoi(c.Query("status"))

	if err := h.jobService.UpdateStatus(uint(id), status); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetExpiredJobs 获取已过期职位
func (h *JobHandler) GetExpiredJobs(c *gin.Context) {
	jobs, err := h.jobService.GetExpiredJobs()
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(jobs))
}

// SearchByCompany 获取公司发布的职位
func (h *JobHandler) SearchByCompany(c *gin.Context) {
	companyID, _ := strconv.ParseUint(c.Param("companyId"), 10, 32)
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	result, err := h.jobService.SearchByCompany(uint(companyID), page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}
