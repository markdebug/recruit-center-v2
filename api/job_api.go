package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
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
	var job request.CreateJobRequest
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.BadRequest))
		return
	}
	//校验参数

	// 获取当前用户公司ID（假设从JWT中获取）
	companyID := c.GetUint("companyId")
	job.CompanyID = companyID
	//判断公司是否存在
	if err := h.jobService.Create(&job); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(job))
}

// Update 更新职位
func (h *JobHandler) Update(c *gin.Context) {
	var job dao.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.BadRequest))
		return
	}
	//判断职位是否存在

	//判断公司是否正确

	// 检查权限（确保是职位所属公司）
	companyID := c.GetUint("company_id")
	if job.CompanyID != companyID {
		c.JSON(http.StatusOK, response.NewError(enums.Forbidden))
		return
	}

	if err := h.jobService.Update(&job); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// Delete 删除职位
func (h *JobHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.BadRequest))
		return
	}
	//判断职位是否存在

	//判断公司是否正确

	if err := h.jobService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetByID 获取职位详情
func (h *JobHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.BadRequest))
		return
	}

	job, err := h.jobService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.JobNotFound))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(job))
}

// List 获取职位列表
func (h *JobHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	jobs, total, err := h.jobService.List(page, size)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(enums.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewPage(jobs, total, page, size))
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
		c.JSON(http.StatusOK, response.NewError(enums.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, result)
}
