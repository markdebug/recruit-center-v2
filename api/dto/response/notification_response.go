package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
)

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Type       int       `json:"type"`
	IsRead     bool      `json:"isRead"`
	CreateTime time.Time `json:"createTime"`
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
