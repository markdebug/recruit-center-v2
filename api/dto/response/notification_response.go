package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
)

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID         uint      `json:"id"`         // 通知ID
	UserID     uint      `json:"userId"`     // 用户ID
	Title      string    `json:"title"`      // 通知标题
	Content    string    `json:"content"`    // 通知内容
	Type       int       `json:"type"`       // 通知类型
	IsRead     bool      `json:"isRead"`     // 是否已读
	CreateTime time.Time `json:"createTime"` // 创建时间
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	Total   int64                  `json:"total"`
	Records []NotificationResponse `json:"records"`
}

// FromModel 从模型转换为响应
func (r *NotificationResponse) FromModel(n *model.Notification) {
	r.ID = n.ID
	r.UserID = n.UserID
	r.Title = n.Title
	r.Content = n.Content
	r.Type = int(n.Type)
	r.IsRead = n.IsRead
	r.CreateTime = n.CreateTime
}
