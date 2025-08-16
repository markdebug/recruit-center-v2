package model

import "time"

// JobApply 职位申请记录
type JobApply struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	JobID         uint      `gorm:"not null;index:idx_job_company_status,priority:1" json:"jobId"`
	CompanyID     uint      `gorm:"not null;index:idx_job_company_status,priority:2" json:"companyId"` // 企业ID
	UserID        uint      `gorm:"not null;index:idx_user_resume,priority:1" json:"userId"`
	ResumeID      uint      `gorm:"not null;index:idx_user_resume,priority:2" json:"resumeId"`
	ApplyTime     time.Time `gorm:"not null;index" json:"applyTime"`
	ApplyProgress string    `gorm:"size:50" json:"applyProgress"`
	Reason        string    `gorm:"size:255" json:"reason"`                                          //拒绝原因
	Status        int       `gorm:"default:1;index:idx_job_company_status,priority:3" json:"status"` // 状态 1: 正常 0: 删除
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
