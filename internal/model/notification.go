package model

import "time"

// UserType 用户类型
type UserType int

const (
	UserTypeJobSeeker UserType = iota + 1 // 求职者
	UserTypeRecruiter                     // 招聘者
	UserTypeAdmin                         // 管理员
	UserTypeCompany                       // 企业用户
)

// NotificationType 通知类型
type NotificationType int

const (
	NotificationTypeJobApply     NotificationType = iota + 1 // 职位申请
	NotificationTypeStatusUpdate                             // 申请状态更新
	NotificationTypeInterview                                // 面试通知
	NotificationTypeSystem                                   // 系统通知
)

// NotificationChannel 通知渠道
type NotificationChannel int

const (
	ChannelInApp  NotificationChannel = 1 << iota // 站内信
	ChannelEmail                                  // 邮件
	ChannelSMS                                    // 短信
	ChannelWechat                                 // 微信
)

// Notification 通知模型
type Notification struct {
	ID         uint                `gorm:"primarykey" json:"id"`
	UserID     uint                `gorm:"not null;index:idx_user_type_read,priority:1" json:"userId"`                          // 接收者ID
	UserType   UserType            `gorm:"not null;index:idx_user_type_read,priority:2" json:"userType"`                        // 用户类型
	Type       NotificationType    `gorm:"not null;index:idx_type_read" json:"type"`                                            // 通知类型
	Title      string              `gorm:"size:100" json:"title"`                                                               // 通知标题
	Content    string              `gorm:"size:500" json:"content"`                                                             // 通知内容
	Channels   NotificationChannel `gorm:"not null" json:"channels"`                                                            // 通知渠道
	TemplateID string              `gorm:"size:50" json:"templateId"`                                                           // 模板ID
	Variables  string              `gorm:"type:json" json:"variables"`                                                          // 模板变量
	IsRead     bool                `gorm:"default:false;index:idx_user_type_read,priority:3;index:idx_type_read" json:"isRead"` // 是否已读
	CreateTime time.Time           `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time           `gorm:"autoUpdateTime" json:"updateTime"`
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID         uint                `gorm:"primarykey" json:"id"`
	Code       string              `gorm:"size:50;uniqueIndex:idx_code_active" json:"code"`    // 模板代码
	Title      string              `gorm:"size:100" json:"title"`                              // 模板标题
	Content    string              `gorm:"size:1000" json:"content"`                           // 模板内容
	Type       NotificationType    `gorm:"not null" json:"type"`                               // 通知类型
	UserTypes  []UserType          `gorm:"type:json" json:"userTypes"`                         // 适用的用户类型
	Channels   NotificationChannel `gorm:"not null" json:"channels"`                           // 支持的通知渠道
	IsActive   bool                `gorm:"default:true;index:idx_code_active" json:"isActive"` // 是否启用
	CreateTime time.Time           `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time           `gorm:"autoUpdateTime" json:"updateTime"`
}

// IsChannelEnabled 检查通知渠道是否启用
func (n *Notification) IsChannelEnabled(channel NotificationChannel) bool {
	return n.Channels&channel != 0
}

func (n *Notification) TableName() string {
	return "t_rc_notification"
}

func (t *NotificationTemplate) TableName() string {
	return "t_rc_notification_template"
}

// AddChannel 添加通知渠道
func (n *Notification) AddChannel(channel NotificationChannel) {
	n.Channels |= channel
}
