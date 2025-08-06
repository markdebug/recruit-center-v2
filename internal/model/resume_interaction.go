package model

import (
	"time"

	"gorm.io/gorm"
)

// InteractionType 定义交互类型
type InteractionType string

const (
	InteractionView     InteractionType = "view"     // 查看
	InteractionFavorite InteractionType = "favorite" // 收藏
)

// ResumeInteraction 简历交互记录
type ResumeInteraction struct {
	gorm.Model
	ResumeID uint            `gorm:"index:idx_resume_user_type"`
	UserID   uint            `gorm:"index:idx_resume_user_type"`
	Type     InteractionType `gorm:"index:idx_resume_user_type;type:varchar(20)"`
	LastTime time.Time       // 最后交互时间
}

// TableName 设置 ResumeInteraction 的表名
func (ResumeInteraction) TableName() string {
	return "t_rc_resume_interaction"
}
