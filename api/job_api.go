package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/entity/rsp"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobHandler 职位接口处理器
type JobHandler struct {
	jobService *service.JobService
}

// NewJobHandler 创建职位处理器实例
func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

// Create 创建职位
func (h *JobHandler) Create(c *gin.Context) {
	var job dao.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(200, rsp.NewError(enums.BadRequest))
		return
	}

	// 获取当前用户公司ID（假设从JWT中获取）
	companyID := c.GetUint("company_id")
	job.CompanyID = companyID

	if err := h.jobService.Create(&job); err != nil {
		c.JSON(200, rsp.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(200, rsp.NewSuccess(job))
}

// Update 更新职位
func (h *JobHandler) Update(c *gin.Context) {
	var job dao.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(200, rsp.NewError(enums.BadRequest))
		return
	}

	// 检查权限（确保是职位所属公司）
	companyID := c.GetUint("company_id")
	if job.CompanyID != companyID {
		c.JSON(200, rsp.NewError(enums.Forbidden))
		return
	}

	if err := h.jobService.Update(&job); err != nil {
		c.JSON(200, rsp.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(200, rsp.NewSuccess(nil))
}

// Delete 删除职位
func (h *JobHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, rsp.NewError(enums.BadRequest))
		return
	}

	if err := h.jobService.Delete(uint(id)); err != nil {
		c.JSON(200, rsp.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(200, rsp.NewSuccess(nil))
}

// GetByID 获取职位详情
func (h *JobHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, rsp.NewError(enums.BadRequest))
		return
	}

	job, err := h.jobService.GetByID(uint(id))
	if err != nil {
		c.JSON(200, rsp.NewError(enums.JobNotFound))
		return
	}

	c.JSON(200, rsp.NewSuccess(job))
}

// List 获取职位列表
func (h *JobHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	jobs, total, err := h.jobService.List(page, size)
	if err != nil {
		c.JSON(200, rsp.NewError(enums.InternalServerError))
		return
	}

	c.JSON(200, rsp.NewPage(jobs, total, page, size))
}

// Search 搜索职位
func (h *JobHandler) Search(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	jobType := c.Query("type")
	location := c.Query("location")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 构建查询条件
	conditions := make(map[string]interface{})
	if keyword != "" {
		conditions["keyword"] = keyword
	}
	if jobType != "" {
		conditions["job_type"] = jobType
	}
	if location != "" {
		conditions["job_location"] = location
	}

	// 执行搜索
	result, err := h.jobService.SearchByCondition(conditions, page, size)
	if err != nil {
		c.JSON(200, rsp.NewError(enums.InternalServerError))
		return
	}

	c.JSON(200, result)
}
