package model

import (
	"time"

	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobType 职位类型
type JobType int

const (
	FullTime   JobType = 1 // 全职
	PartTime   JobType = 2 // 兼职
	Internship JobType = 3 // 实习
)

// Job 职位信息
type Job struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	CompanyID     uint      `gorm:"not null" json:"companyId"`
	JobSkill      string    `gorm:"size:500" json:"jobSkill"`
	JobSalary     string    `gorm:"size:50" json:"jobSalary"`
	JobDescribe   string    `gorm:"type:text" json:"jobDescribe"`
	JobLocation   string    `gorm:"size:200" json:"jobLocation"`
	JobExpireTime time.Time `gorm:"index" json:"jobExpireTime"`
	Status        int       `gorm:"default:1" json:"status"`
	JobType       int       `gorm:"size:50" json:"jobType"`
	JobCategory   string    `gorm:"size:50" json:"jobCategory"`
	JobExperience string    `gorm:"size:50" json:"jobExperience"`
	JobEducation  string    `gorm:"size:50" json:"jobEducation"`
	JobBenefit    string    `gorm:"size:500" json:"jobBenefit"`
	JobContact    string    `gorm:"size:100" json:"jobContact"`
	DeleteStatus  int       `gorm:"default:0" json:"deleteStatus"`
	JobSource     string    `gorm:"size:100" json:"jobSource"`
	CreateTime    time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"updateTime"`

	// 添加新的字段和方法
	ViewCount  int      `gorm:"default:0" json:"viewCount"`  // 浏览次数
	ApplyCount int      `gorm:"default:0" json:"applyCount"` // 申请次数
	Priority   int      `gorm:"default:0" json:"priority"`   // 优先级
	Tags       []string `gorm:"type:json" json:"tags"`       // 职位标签

	Applications []JobApply `gorm:"foreignKey:JobID" json:"-"`
}

// TableName 指定表名
func (Job) TableName() string {
	return "t_rc_job"
}

func NewJob(name string, companyID uint) *Job {
	return &Job{
		Name:      name,
		CompanyID: companyID,
	}
}

// IsExpired 检查职位是否过期
func (j *Job) IsExpired() bool {
	return j.JobExpireTime.Before(time.Now())
}

// IsActive 检查职位是否有效
func (j *Job) IsActive() bool {
	return j.Status == int(enums.JobStatusNormal) && !j.IsExpired()
}

// ValidateJobType 验证职位类型
func (j *Job) ValidateJobType() bool {
	return JobType(j.JobType) >= FullTime && JobType(j.JobType) <= Internship
}
