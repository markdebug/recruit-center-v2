package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

type JobApplyHandler struct {
	jobApplyService *service.JobApplyService
	jobService      *service.JobService
}

func NewJobApplyHandler(service *service.JobApplyService, jobService *service.JobService) *JobApplyHandler {
	return &JobApplyHandler{jobApplyService: service, jobService: jobService}
}

// Create 创建职位申请
//
//	@Summary		创建职位申请
//	@Description	创建一个新的职位申请
//	@Tags			职位申请
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer 用户令牌"
//	@Param			apply			body		request.JobApply	true	"申请信息"
//	@Success		0000			{object}	response.Response{data=model.JobApply}
//	@Failure		2000			{object}	response.Response{}
//	@Router			/api/v1/applies [post]
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
//
//	@Summary		删除职位申请
//	@Description	删除指定ID的职位申请
//	@Tags			职位申请
//	@Produce		json
//	@Param			id		path		int	true	"申请ID"
//	@Success		0000	{object}	response.Response{data=string}
//	@Failure		2000	{object}	response.Response{}
//	@Router			/api/v1/applies/{id} [delete]
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
//
//	@Summary		获取申请详情
//	@Description	获取指定ID的职位申请详情
//	@Tags			职位申请
//	@Produce		json
//	@Param			id		path		int	true	"申请ID"
//	@Success		0000	{object}	response.Response{data=model.JobApply}
//	@Failure		2000	{object}	response.Response{}
//	@Router			/api/v1/applies/{id} [get]
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

// 根据用户id，查询其全部的申请信息
// ListByUser 获取用户的申请记录
//
//	@Summary		获取用户申请记录
//	@Description	分页获取指定用户的职位申请记录
//	@Tags			职位申请
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			Authorization	header		string											true	"Bearer
//	@Param			userId			path		int												true	"用户ID"
//	@Param			page			query		integer											false	"页码 (默认
//	@Param			size			query		integer											false	"每页数量 (默认值: 10)"	minimum(1)	maximum(100)	default(10)
func (h *JobApplyHandler) ListByUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))
	page, size := parsePageSize(c)

	applies, err := h.jobApplyService.ListByUser(uint(userID), page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(response.NewPage(applies, applies.Total, page, size)))
}

// 根据公司信息，查询所有的职位申请记录
// ListByCompany 获取公司职位申请记录
//
//	@Summary		获取公司职位申请记录
//	@Description	分页获取指定公司的
//	职位申请记录
//	@Tags			职位申请
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			Authorization	header		string											true	"Bearer
//	@Param			companyId		path		int												true	"公司ID"
//	@Param			page			query		integer											false	"页码 (默认值: 1)"		minimum(1)	default(1)
//	@Param			size			query		integer											false	"每页数量 (默认值: 10)"	minimum(1)	maximum(100)	default(10)
func (h *JobApplyHandler) ListByCompany(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("companyId"))
	page, size := parsePageSize(c)

	applies, err := h.jobApplyService.ListByCompanyID(uint(companyID), page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(response.NewPage(applies, applies.Total, page, size)))
}

// List 获取职位申请列表
//
//	@Summary		获取申请列表
//	@Description	分页获取所有职位申请记录
//	@Tags			职位申请
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			Authorization	header		string											true	"Bearer JWT"
//	@Param			page			query		integer											false	"页码 (默认值: 1)"		minimum(1)	default(1)
//	@Param			size			query		integer											false	"每页数量 (默认值: 10)"	minimum(1)	maximum(100)	default(10)
//	@Success		0000			{object}	response.PageResponse{data=[]model.JobApply}	"成功"
//	@Failure		2000			{object}	response.Response{}								"错误"
//	@Router			/api/v1/applies [get]
func (h *JobApplyHandler) List(c *gin.Context) {
	page, size := parsePageSize(c)
	jobID := c.GetUint("jobId")
	applies, err := h.jobApplyService.ListByJob(jobID, page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(response.NewPage(applies, applies.Total, page, size)))
}

// 更新职位申请状态
func (h *JobApplyHandler) UpdateStatus(c *gin.Context) {
	var req request.JobApplyUpdateStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}
	//通过jobID查询job详情
	job, err := h.jobService.GetByID(req.JobID)
	if err != nil {
		logger.L.Error("查询职位信息失败",
			zap.Error(err),
			zap.Uint("jobId", req.JobID))
		c.JSON(http.StatusOK, errors.Wrap(err, errors.JobNotFound))
		return
	}
	//判断职位是否是当前公司的职位
	if job.CompanyID != c.GetUint("companyId") {
		logger.L.Warn("职位不属于当前公司",
			zap.Uint("jobId", req.JobID),
			zap.Uint("companyId", c.GetUint("companyId")))
		c.JSON(http.StatusOK, errors.Forbidden)
	}

	if err := h.jobApplyService.UpdateStatus(req.JobID, req.Status); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// parsePageSize 解析分页参数
func parsePageSize(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	return page, size
}
