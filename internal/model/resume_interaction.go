package model

import "time"

// InteractionType 定义交互类型
type InteractionType string

const (
	InteractionView     InteractionType = "view"     // 查看
	InteractionFavorite InteractionType = "favorite" // 收藏
)

// ResumeInteraction 简历交互记录
type ResumeInteraction struct {
	Id        uint            `gorm:"primarykey"`
	ResumeID  uint            `gorm:"index:idx_resume_user_type,priority:1"`
	UserID    uint            `gorm:"index:idx_resume_user_type,priority:2"`
	Type      InteractionType `gorm:"index:idx_resume_user_type,priority:3;type:varchar(20)"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time
}

// TableName 设置 ResumeInteraction 的表名
func (ResumeInteraction) TableName() string {
	return "t_rc_resume_interaction"
}
