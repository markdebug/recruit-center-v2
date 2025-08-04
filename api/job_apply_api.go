package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobApplyHandler 职位申请处理器
type JobApplyHandler struct {
	jobApplyService *service.JobApplyService
}

// NewJobApplyHandler 创建职位申请处理器
func NewJobApplyHandler(service *service.JobApplyService) *JobApplyHandler {
	return &JobApplyHandler{
		jobApplyService: service,
	}
}

// Create 创建职位申请
func (h *JobApplyHandler) Create(c *gin.Context) {
	var apply dao.JobApply
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	// 获取当前用户ID
	userID := c.GetUint("user_id")
	apply.UserID = userID

	if err := h.jobApplyService.Create(&apply); err != nil {
		c.JSON(200, response.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(200, response.NewSuccess(apply))
}

// ListByUser 获取用户的申请列表
func (h *JobApplyHandler) ListByUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	applies, total, err := h.jobApplyService.ListByUser(userID, page, size)
	if err != nil {
		c.JSON(200, response.NewError(enums.InternalServerError))
		return
	}

	c.JSON(200, response.NewPage(applies, total, page, size))
}

// ListByJob 获取职位的申请列表
func (h *JobApplyHandler) ListByJob(c *gin.Context) {
	jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	applies, total, err := h.jobApplyService.ListByJob(uint(jobID), page, size)
	if err != nil {
		c.JSON(200, response.NewError(enums.InternalServerError))
		return
	}

	c.JSON(200, response.NewPage(applies, total, page, size))
}

// UpdateStatus 更新申请状态
func (h *JobApplyHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	// 验证状态值是否有效
	if !enums.JobApplyEnum(req.Status).IsValid() {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	if err := h.jobApplyService.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(200, response.NewErrorWithMsg(enums.InternalServerError, err.Error()))
		return
	}

	c.JSON(200, response.NewSuccess(nil))
}

// GetByID 获取申请详情
func (h *JobApplyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(200, response.NewError(enums.BadRequest))
		return
	}

	apply, err := h.jobApplyService.GetByID(uint(id))
	if err != nil {
		c.JSON(200, response.NewError(enums.NotFound))
		return
	}

	c.JSON(200, response.NewSuccess(apply))
}
