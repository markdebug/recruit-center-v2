package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type ResumeHandler struct {
	resumeService *service.ResumeService
}

func NewResumeHandler(resumeService *service.ResumeService) *ResumeHandler {
	return &ResumeHandler{resumeService: resumeService}
}

// Create 创建简历
// @Summary 创建简历
// @Description 创建用户简历信息
// @Tags 简历管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param resume body request.CreateResumeRequest true "简历信息"
// @Success 200 {object} response.Response{data=model.Resume}
// @Failure 400 {object} response.Response
// @Router /api/v1/resumes [post]
func (h *ResumeHandler) Create(c *gin.Context) {
	var req request.CreateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	userID := c.GetUint("userId")
	resume, err := h.resumeService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(resume))
}

// GetByUser 获取当前用户的简历
// @Summary 获取用户简历
// @Description 获取当前登录用户的简历信息
// @Tags 简历管理
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} response.Response{data=model.Resume}
// @Failure 404 {object} response.Response
// @Router /api/v1/resumes/my [get]
func (h *ResumeHandler) GetByUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	resume, err := h.resumeService.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.NotFound))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(resume))
}

// UploadResume 上传简历文件
// @Summary 上传简历文件
// @Description 上传用户简历文件
// @Tags 简历管理
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param resume formData file true "简历文件"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/resumes/upload [post]
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetUint("user_id")

	// 获取上传的文件
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, "无法读取文件"))
		return
	}
	defer src.Close()

	// 调用service处理上传
	fileURL, err := h.resumeService.UploadResumeFile(userID, src, file.Filename)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(gin.H{
		"file_url": fileURL,
	}))
}

// UpdateAccessStatus 更新简历访问状态
// @Summary 更新简历访问状态
// @Description 更新用户简历的访问状态（公开/隐藏）
// @Tags 简历管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param status body request.UpdateResumeStatusRequest true "访问状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/resumes/access-status [put]
func (h *ResumeHandler) UpdateAccessStatus(c *gin.Context) {
	var req request.UpdateResumeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	userID := c.GetUint("userId")
	if err := h.resumeService.UpdateAccessStatus(userID, req.Status); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// UpdateWorkingStatus 更新简历工作状态
// @Summary 更新简历工作状态
// @Description 更新用户简历的工作状态（在职/离职）
// @Tags 简历管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param status body request.UpdateResumeStatusRequest true "工作状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/resumes/working-status [put]
func (h *ResumeHandler) UpdateWorkingStatus(c *gin.Context) {
	var req request.UpdateResumeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	userID := c.GetUint("userId")
	if err := h.resumeService.UpdateWorkingStatus(userID, req.Status); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}
