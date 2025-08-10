package request

// NotificationMarkReadRequest 标记通知已读请求
type NotificationMarkReadRequest struct {
	NotificationIDs []uint `json:"notificationIds" binding:"required,min=1"`
}

// NotificationSendRequest 发送通知请求
type NotificationSendRequest struct {
	UserID     uint                   `json:"userId" binding:"required"`     // 接收者ID
	UserType   int                    `json:"userType" binding:"required"`   // 用户类型
	TemplateID string                 `json:"templateId" binding:"required"` // 模板ID
	Variables  map[string]interface{} `json:"variables"`                     // 模板变量
}
