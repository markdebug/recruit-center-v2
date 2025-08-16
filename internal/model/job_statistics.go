package model

import "time"

// JobStatistics 职位统计信息
type JobStatistics struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	JobID          uint      `gorm:"not null;index:idx_job_company,priority:1;uniqueIndex" json:"jobId"`
	CompanyID      uint      `gorm:"not null;index:idx_job_company,priority:2" json:"companyId"`
	ViewCount      int64     `gorm:"default:0" json:"viewCount"`      // 浏览数
	ApplyCount     int64     `gorm:"default:0" json:"applyCount"`     // 申请数
	ConversionRate float64   `gorm:"default:0" json:"conversionRate"` // 转化率 (申请数/浏览数)
	LastViewTime   time.Time `json:"lastViewTime"`                    // 最后浏览时间
	LastApplyTime  time.Time `json:"lastApplyTime"`                   // 最后申请时间
	CreateTime     time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (JobStatistics) TableName() string {
	return "t_rc_job_statistics"
}
