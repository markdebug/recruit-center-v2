package model

import "time"

// JobFavorite 职位收藏
type JobFavorite struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"not null;index:idx_user_job,priority:1" json:"userId"` // 用户ID
	JobID      uint      `gorm:"not null;index:idx_user_job,priority:2" json:"jobId"`  // 职位ID
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`                     // 收藏时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`                     // 更新时间
}

func (JobFavorite) TableName() string {
	return "t_rc_job_favorite"
}
