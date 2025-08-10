package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: service}
}

// List 获取通知列表
//	@Summary		获取通知列表
//	@Description	分页获取当前用户的通知列表
//	@Tags			通知管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Param			page			query		int		false	"页码 (默认值: 1)"		minimum(1)	default(1)
//	@Param			size			query		int		false	"每页数量 (默认值: 10)"	minimum(1)	maximum(100)	default(10)
//	@Success		200				{object}	response.Response{data=response.NotificationListResponse}
//	@Router			/api/v1/notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetUint("userId")
	page, size := parsePageSize(c)

	notifications, total, err := h.notificationService.ListUserNotifications(userID, page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	resp := &response.NotificationListResponse{
		Total:   total,
		Records: make([]response.NotificationResponse, len(notifications)),
	}

	for i, n := range notifications {
		resp.Records[i].FromModel(&n)
	}

	c.JSON(http.StatusOK, response.NewSuccess(resp))
}

// GetUnreadCount 获取未读通知数量
//	@Summary		获取未读通知数量
//	@Description	获取当前用户的未读通知数量
//	@Tags			通知管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer 用户令牌"
//	@Success		200				{object}	response.Response{data=int64}
//	@Router			/api/v1/notifications/unread/count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userId")
	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(count))
}

// MarkAsRead 标记通知为已读
//	@Summary		标记通知为已读
//	@Description	标记指定的通知为已读状态
//	@Tags			通知管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Bearer 用户令牌"
//	@Param			request			body		request.NotificationMarkReadRequest	true	"通知ID列表"
//	@Success		200				{object}	response.Response
//	@Router			/api/v1/notifications/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	var req request.NotificationMarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	for _, id := range req.NotificationIDs {
		if err := h.notificationService.MarkAsRead(id); err != nil {
			c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
			return
		}
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// Send 发送通知
//	@Summary		发送通知
//	@Description	发送通知到指定用户
//	@Tags			通知管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer 用户令牌"
//	@Param			request			body		request.NotificationSendRequest	true	"发送通知请求"
//	@Success		200				{object}	response.Response
//	@Router			/api/v1/notifications/send [post]
func (h *NotificationHandler) Send(c *gin.Context) {
	var req request.NotificationSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errors.BadRequest)
		return
	}

	err := h.notificationService.SendNotification(
		req.UserID,
		model.UserType(req.UserType),
		req.TemplateID,
		req.Variables,
	)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}
