package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type JobApplyHandler struct {
	jobApplyService *service.JobApplyService
}

func NewJobApplyHandler(service *service.JobApplyService) *JobApplyHandler {
	return &JobApplyHandler{jobApplyService: service}
}

// Create 创建职位申请
// @Summary 创建职位申请
// @Description 创建一个新的职位申请
// @Tags 职位申请
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param apply body request.JobApply true "申请信息"
// @Success 200 {object} response.Response{data=model.JobApply}
// @Failure 400 {object} response.Response
// @Router /api/v1/applies [post]
func (h *JobApplyHandler) Create(c *gin.Context) {
	var req request.JobApply
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	userID := c.GetUint("userId")
	req.UserID = userID
	req.ApplyTime = time.Now()
	apply := req.ToModel()

	if err := h.jobApplyService.Create(apply); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(apply))
}

// Delete 删除职位申请
func (h *JobApplyHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	if err := h.jobApplyService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetByID 获取职位申请详情
// @Summary 获取申请详情
// @Description 获取指定ID的职位申请详情
// @Tags 职位申请
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response{data=model.JobApply}
// @Failure 404 {object} response.Response
// @Router /api/v1/applies/{id} [get]
func (h *JobApplyHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	apply, err := h.jobApplyService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(apply))
}

// List 获取职位申请列表
// @Summary 获取申请列表
// @Description 分页获取所有职位申请记录
// @Tags 职位申请
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer JWT"
// @Param page query integer false "页码 (默认值: 1)" minimum(1) default(1)
// @Param size query integer false "每页数量 (默认值: 10)" minimum(1) maximum(100) default(10)
// @Success 200 {object} response.PageResponse{data=[]model.JobApply} "成功"
// @Failure 400 {object} response.Response{} "请求参数错误"
// @Failure 401 {object} response.Response{} "未授权"
// @Failure 500 {object} response.Response{} "服务器内部错误"
// @Router /api/v1/applies [get]
func (h *JobApplyHandler) List(c *gin.Context) {
	page, size := parsePageSize(c)
	applies, total, err := h.jobApplyService.List(page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(response.NewPage(applies, total, page, size)))
}

// parsePageSize 解析分页参数
func parsePageSize(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	return page, size
}
