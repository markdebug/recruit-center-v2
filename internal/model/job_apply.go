package model

import "time"

// JobApply 职位申请记录
type JobApply struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	JobID         uint      `gorm:"not null;index" json:"jobId"`
	UserID        uint      `gorm:"not null;index" json:"userId"`
	ApplyTime     time.Time `gorm:"not null" json:"applyTime"`
	ApplyProgress string    `gorm:"size:50" json:"applyProgress"`
	Status        int       `gorm:"default:1" json:"status"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// TableName 指定表名
func (JobApply) TableName() string {
	return "t_rc_job_apply"
}

func NewJobApply(jobID, userID uint) *JobApply {
	return &JobApply{
		JobID:     jobID,
		UserID:    userID,
		Status:    1,
		ApplyTime: time.Now(),
	}
}
