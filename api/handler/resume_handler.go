package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type ResumeHandler struct {
	resumeService      *service.ResumeService
	interactionService *service.ResumeInteractionService
}

func NewResumeHandler(resumeService *service.ResumeService, interactionService *service.ResumeInteractionService) *ResumeHandler {
	return &ResumeHandler{resumeService: resumeService, interactionService: interactionService}
}

// Create 创建简历
//
//	@Summary		创建简历
//	@Description	创建用户简历信息
//	@Tags			简历管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer 用户令牌"
//	@Param			resume			body		request.CreateResumeRequest	true	"简历信息"
//	@Success		0000			{object}	response.Response{data=model.Resume}
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes [post]
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

// Update 更新简历
//
//	@Summary		更新简历
//	@Description	更新用户简历信息
//	@Tags			简历管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer 用户令牌"
//	@Param			id				path		int							true	"简历ID"
//	@Param			resume			body		request.UpdateResumeRequest	true	"更新简历信息"
//	@Success		0000			{object}	response.Response
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/{id} [put]
func (h *ResumeHandler) Update(c *gin.Context) {
	var req request.UpdateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	resumeID := c.GetUint("id")
	userID := c.GetUint("userId")

	// 判断简历是否存在
	_, err := h.resumeService.GetByUserIDAndResumeID(userID, resumeID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	// 根据模块类型分别处理更新
	switch req.Module {
	case request.ModuleBasic:
		var basicData request.UpdateResumeBasicRequest
		if err := mapstructure.Decode(req.Data, &basicData); err != nil {
			c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
			return
		}
		err = h.resumeService.UpdateBasic(resumeID, &basicData)

	case request.ModuleEducation:
		var eduData request.UpdateResumeEducationRequest
		if err := mapstructure.Decode(req.Data, &eduData); err != nil {
			c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
			return
		}
		err = h.resumeService.UpdateEducation(resumeID, &eduData)

	case request.ModuleWorkExperience:
		var workData request.UpdateResumeWorkRequest
		if err := mapstructure.Decode(req.Data, &workData); err != nil {
			c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
			return
		}
		err = h.resumeService.UpdateWork(resumeID, &workData)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// GetByUser 获取当前用户的简历
//
//	@Summary		获取用户简历
//	@Description	获取当前登录用户的简历信息
//	@Tags			简历管理
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Success		0000			{object}	response.Response{data=model.Resume}
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/my [get]
func (h *ResumeHandler) GetByUser(c *gin.Context) {
	userID := c.GetUint("userId")
	resume, err := h.resumeService.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.NotFound))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(resume))
}

// GetByID 获取简历详情
//
//	@Summary		获取简历详情
//	@Description	获取指定ID的简历详细信息
//	@Tags			简历管理
//	@Produce		json
//	@Param			id		path		int	true	"简历ID"
//	@Success		0000	{object}	response.Response{data=model.Resume}
//	@Failure		5000	{object}	response.Response
//	@Router			/api/v1/resumes/{id} [get]
func (h *ResumeHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	resume, err := h.resumeService.GetResumeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.ResumeNotFound))
		return
	}
	c.JSON(http.StatusOK, response.NewSuccess(resume))
}

// 根据分享token获取用户简历
//
//	@Summary		根据分享token获取简历
//	@Description	根据分享token获取简历信息
//	@Tags			简历管理
//	@Produce		json
//	@Param			token	path		string	true	"分享Token"
//	@Success		0000	{object}	response.Response{data=model.Resume}
//	@Failure		5000	{object}	response.Response
//	@Router			/api/v1/resumes/share/{token} [get]
func (h *ResumeHandler) GetByShareToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	resume, err := h.resumeService.GetByShareToken(token)
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.NotFound))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(resume))
}

// UploadResume 上传简历文件
//
//	@Summary		上传简历文件
//	@Description	上传用户简历文件
//	@Tags			简历管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			resume			formData	file	true	"简历文件"
//	@Success		0000			{object}	response.Response{data=map[string]string}
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/upload [post]
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetUint("userId")

	// 获取上传的文件
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusOK, response.NewError(errors.InvalidParams))
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.FileUploadFailed, "无法读取文件"))
		return
	}
	defer src.Close()

	// 调用service处理上传
	fileURL, err := h.resumeService.UploadResumeFile(userID, src, file.Filename)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.FileUploadFailed, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(gin.H{
		"fileUrl": fileURL,
	}))
}

// UpdateAccessStatus 更新简历访问状态
//
//	@Summary		更新简历访问状态
//	@Description	更新用户简历的访问状态（公开/隐藏）
//	@Tags			简历管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer 用户令牌"
//	@Param			status			body		request.UpdateResumeStatusRequest	true	"访问状态"
//	@Success		0000			{object}	response.Response
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/access-status [put]
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
//
//	@Summary		更新简历工作状态
//	@Description	更新用户简历的工作状态（在职/离职）
//	@Tags			简历管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer 用户令牌"
//	@Param			status			body		request.UpdateResumeStatusRequest	true	"工作状态"
//	@Success		0000			{object}	response.Response
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/working-status [put]
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

// ViewResume 查看简历收藏相关信息
//
//	@Summary		查看简历收藏相关信息
//	@Description	记录简历查看行为
//	@Tags			简历管理
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			id				path		int		true	"简历ID"
//	@Success		0000			{object}	response.Response
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/{id}/view [post]
func (h *ResumeHandler) ViewResume(c *gin.Context) {
	resumeID := c.GetUint("id")
	userID := c.GetUint("userId")

	if err := h.interactionService.RecordView(resumeID, userID); err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	stats, err := h.interactionService.GetInteractionStats(resumeID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(stats))
}

// ToggleFavorite 切换简历收藏状态
//
//	@Summary		收藏/取消收藏简历
//	@Description	切换简历的收藏状态
//	@Tags			简历管理
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			id				path		int		true	"简历ID"
//	@Success		0000			{object}	response.Response
//	@Failure		5000			{object}	response.Response
//	@Router			/api/v1/resumes/{id}/favorite [post]
func (h *ResumeHandler) ToggleFavorite(c *gin.Context) {
	resumeID := c.GetUint("id")
	userID := c.GetUint("userId")

	isFavorited, err := h.interactionService.IsFavorited(resumeID, userID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	if isFavorited {
		err = h.interactionService.RemoveFavorite(resumeID, userID)
	} else {
		err = h.interactionService.AddFavorite(resumeID, userID)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	stats, err := h.interactionService.GetInteractionStats(resumeID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(stats))
}

// 获取简历的统计信息
//
//	@Summary		获取简历统计信息
//	@Description	获取简历的查看和收藏统计信息
//	@Tags			简历管理
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			id				path		int		true	"简历ID"
//	@Success		0000			{object}	response.Response{data=map[string]int64}
//	@Failure		5000			{object}	response.Response{}
//	@Router			/api/v1/resumes/{id}/stats [get]
func (h *ResumeHandler) GetStats(c *gin.Context) {
	resumeID := c.GetUint("id")

	stats, err := h.interactionService.GetInteractionStats(resumeID)
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorWithMsg(errors.InternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(stats))
}
