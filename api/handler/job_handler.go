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
//
//	@Summary		创建职位
//	@Description	创建新职位,需要提供职位相关信息
//	@Tags			职位
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			Authorization	header		string						true	"Bearer JWT"
//	@Param			request			body		request.CreateJobRequest	true	"职位创建请求参数"
//	@Success		0000			{object}	response.Response{data=model.Job}
//	@Failure		2000			{object}	response.Response{}
//	@Router			/api/v1/jobs [post]
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

	c.JSON(http.StatusOK, response.NewSuccess(h.jobService.ConvertToJobResponse(job)))
}

// Update 更新职位
//
//	@Summary		更新职位
//	@Description	更新指定ID的职位信息
//	@Tags			职位
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			Authorization	header		string								true	"Bearer JWT"
//	@Param			id				path		integer								true	"职位ID"
//	@Param			request			body		request.UpdateJobRequest			true	"职位更新请求参数"
//	@Success		0000			{object}	response.Response{data=model.Job}	"成功"
//	@Failure		2000			{object}	response.Response{}					"请求参数错误"
//	@Router			/api/v1/jobs/{id} [put]
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
//
//	@Summary		删除职位
//	@Description	删除指定ID的职位
//	@Tags			职位
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			Authorization	header		string				true	"Bearer JWT"
//	@Param			id				path		integer				true	"职位ID"
//	@Success		0000			{object}	response.Response{}	"成功"
//	@Failure		2000			{object}	response.Response{}	"未授权"
//	@Router			/api/v1/jobs/{id} [delete]
func (h *JobHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.jobService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetByID 获取职位详情
//
//	@Summary		获取职位详情
//	@Description	获取指定ID的职位详细信息
//	@Tags			职位
//	@Produce		json
//	@Param			id		path		int	true	"职位ID"
//	@Success		0000	{object}	response.Response{data=model.Job}
//	@Failure		2000	{object}	response.Response
//	@Router			/api/v1/jobs/{id} [get]
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
//
//	@Summary		获取职位列表
//	@Description	分页获取职位列表
//	@Tags			职位
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(10)
//	@Success		0000	{object}	response.Response{data=response.JobListResponse}
//	@Failure		2000	{object}	response.Response{}
//	@Router			/api/v1/jobs [get]
func (h *JobHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	result, err := h.jobService.List(page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}

// SearchByKeyword 关键词搜索职位
//
//	@Summary		关键词搜索职位
//	@Description	根据关键词搜索职位信息
//	@Tags			职位
//	@Produce		json
//	@Param			keyword	query		string												true	"搜索关键词"
//	@Success		0000	{object}	response.Response{data=response.JobListResponse}	"成功"
//	@Failure		2000	{object}	response.Response{}									"服务器内部错误"
//	@Router			/api/v1/jobs/search [get]
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
//
//	@Summary		多条件搜索职位
//	@Description	根据多个条件筛选职位
//	@Tags			职位
//	@Produce		json
//	@Param			page		query		integer												false	"页码"	default(1)
//	@Param			size		query		integer												false	"每页数量"	default(10)
//	@Param			conditions	body		object												true	"搜索条件"
//	@Success		0000		{object}	response.Response{data=response.JobListResponse}	"成功"
//	@Failure		2000		{object}	response.Response{}									"请求参数错误"
//	@Router			/api/v1/jobs/search/condition [post]
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

// GetExpiredJobs 获取已过期职位
//
//	@Summary		获取已过期职位
//	@Description	获取所有已过期的职位列表
//	@Tags			职位
//	@Produce		json
//	@Security		Bearer
//	@Param			Authorization	header		string												true	"Bearer JWT"
//	@Success		0000			{object}	response.Response{data=response.JobListResponse}	"成功"
//	@Failure		2000			{object}	response.Response{}									"服务器内部错误"
//	@Router			/api/v1/jobs/expired [get]
func (h *JobHandler) GetExpiredJobs(c *gin.Context) {
	result, err := h.jobService.GetExpiredJobs()
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}

// SearchByCompany 获取公司发布的职位
//
//	@Summary		获取公司职位
//	@Description	获取指定公司发布的所有职位
//	@Tags			职位
//	@Produce		json
//	@Param			companyId	path		integer												true	"公司ID"
//	@Param			page		query		integer												false	"页码"	default(1)
//	@Param			size		query		integer												false	"每页数量"	default(10)
//	@Success		0000		{object}	response.Response{data=response.JobListResponse}	"成功"
//	@Failure		2000		{object}	response.Response{}									"服务器内部错误"
//	@Router			/api/v1/companies/{companyId}/jobs [get]
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
