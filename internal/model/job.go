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

// RemoteType 远程办公类型
type RemoteType int

const (
	OnSite     RemoteType = 1 // 纯办公室办公
	Hybrid     RemoteType = 2 // 混合办公(居家+办公室)
	FullRemote RemoteType = 3 // 全远程
	Flexible   RemoteType = 4 // 灵活办公(员工自主选择)
)

// JobBenefitType 职位福利类型
type JobBenefitType int

const (
	Insurance    JobBenefitType = 1  // 五险一金
	Bonus        JobBenefitType = 2  // 年终奖
	Leave        JobBenefitType = 3  // 带薪休假
	Training     JobBenefitType = 4  // 培训发展
	Meals        JobBenefitType = 5  // 餐补
	Transport    JobBenefitType = 6  // 交通补助
	Stock        JobBenefitType = 7  // 股票期权
	FlexibleTime JobBenefitType = 8  // 弹性工作
	Healthcare   JobBenefitType = 9  // 医疗保险
	Gym          JobBenefitType = 10 // 健身设施
)

// Job 职位信息
type Job struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	CompanyID     uint      `gorm:"not null" json:"companyId"`
	JobSkill      string    `gorm:"size:500" json:"jobSkill"`
	JobSalary     int       `gorm:"size:50" json:"jobSalary"`
	JobSalaryMax  int       `gorm:"size:50" json:"jobSalaryMax"`
	JobDescribe   string    `gorm:"type:text" json:"jobDescribe"`
	JobLocation   string    `gorm:"size:200" json:"jobLocation"`
	JobExpireTime time.Time `gorm:"index" json:"jobExpireTime"`
	Status        int       `gorm:"default:1" json:"status"` // 职位状态 1: 待发布 2: 已发布 3: 暂停发布
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

	// 远程办公相关字段
	RemoteType  RemoteType `gorm:"default:1" json:"remoteType"`  // 远程办公类型
	RemoteDesc  string     `gorm:"size:200" json:"remoteDesc"`   // 远程办公补充说明
	RemoteRatio int        `gorm:"default:0" json:"remoteRatio"` // 远程办公比例(0-100)

	// 福利相关字段
	Benefits    []JobBenefitType `gorm:"type:json" json:"benefits"`   // 福利列表
	BenefitDesc string           `gorm:"size:500" json:"benefitDesc"` // 福利补充说明

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

// IsRemoteEnabled 是否支持远程办公
func (j *Job) IsRemoteEnabled() bool {
	return j.RemoteType != OnSite
}

// GetRemoteTypeText 获取远程办公类型文本
func (j *Job) GetRemoteTypeText() string {
	switch j.RemoteType {
	case OnSite:
		return "办公室办公"
	case Hybrid:
		return "混合办公"
	case FullRemote:
		return "全远程"
	case Flexible:
		return "灵活办公"
	default:
		return "未知"
	}
}

// GetBenefitText 获取福利类型文本
func (j *Job) GetBenefitText(benefit JobBenefitType) string {
	benefitTexts := map[JobBenefitType]string{
		Insurance:    "五险一金",
		Bonus:        "年终奖",
		Leave:        "带薪休假",
		Training:     "培训发展",
		Meals:        "餐补",
		Transport:    "交通补助",
		Stock:        "股票期权",
		FlexibleTime: "弹性工作",
		Healthcare:   "医疗保险",
		Gym:          "健身设施",
	}
	return benefitTexts[benefit]
}
