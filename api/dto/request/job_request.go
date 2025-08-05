package request

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// CreateJobRequest 创建职位请求
type CreateJobRequest struct {
	CompanyID     uint          `json:"companyId" binding:"required" example:"1"`
	Name          string        `json:"name" binding:"required,min=1,max=100" example:"高级Go开发工程师"`
	JobSkill      string        `json:"jobSkill" binding:"max=500" example:"Go,Redis,MySQL,Docker"`
	JobSalary     string        `json:"jobSalary" binding:"max=50" example:"15k-25k"`
	JobDescribe   string        `json:"jobDescribe" binding:"max=2000" example:"负责后端服务开发..."`
	JobLocation   string        `json:"jobLocation" binding:"max=200" example:"北京市朝阳区"`
	JobExpireTime time.Time     `json:"jobExpireTime,omitempty" example:"2024-12-31T23:59:59Z"`
	JobType       enums.JobType `json:"jobType" binding:"required" example:"full_time"`
	JobCategory   string        `json:"jobCategory" binding:"max=50" example:"技术"`
	JobExperience string        `json:"jobExperience" binding:"max=50" example:"3-5年"`
	JobEducation  string        `json:"jobEducation" binding:"max=50" example:"本科及以上"`
	JobBenefit    string        `json:"jobBenefit" binding:"max=500" example:"五险一金,年终奖"`
	JobContact    string        `json:"jobContact" binding:"max=100" example:"hr@company.com"`
	JobSource     string        `json:"jobSource" binding:"max=100" example:"内推"`
}

// UpdateJobRequest 更新职位请求
type UpdateJobRequest struct {
	ID            uint          `json:"id" binding:"required"`
	Name          string        `json:"name" binding:"min=1,max=100"`
	JobSkill      string        `json:"jobSkill" binding:"max=500"`
	JobSalary     string        `json:"jobSalary" binding:"max=50"`
	JobDescribe   string        `json:"jobDescribe" binding:"max=2000"`
	JobLocation   string        `json:"jobLocation" binding:"max=200"`
	JobExpireTime time.Time     `json:"jobExpireTime,omitempty"`
	JobType       enums.JobType `json:"jobType"`
	JobCategory   string        `json:"jobCategory" binding:"max=50"`
	JobExperience string        `json:"jobExperience" binding:"max=50"`
	JobEducation  string        `json:"jobEducation" binding:"max=50"`
	JobBenefit    string        `json:"jobBenefit" binding:"max=500"`
	JobContact    string        `json:"jobContact" binding:"max=100"`
	JobSource     string        `json:"jobSource" binding:"max=100"`
}

// JobSearchRequest 职位搜索请求
type JobSearchRequest struct {
	Keyword     string        `json:"keyword" form:"keyword" example:"Go开发"`
	JobType     enums.JobType `json:"jobType" form:"jobType" example:"full_time"`
	JobCategory string        `json:"jobCategory" form:"jobCategory" example:"技术"`
	Location    string        `json:"location" form:"location" example:"北京"`
	SalaryMin   int           `json:"salaryMin" form:"salaryMin" example:"10000"`
	SalaryMax   int           `json:"salaryMax" form:"salaryMax" example:"30000"`
	CompanyID   uint          `json:"companyId" form:"companyId" example:"1"`
	Page        int           `json:"page" form:"page" binding:"min=1" example:"1"`
	PageSize    int           `json:"pageSize" form:"pageSize" binding:"min=1,max=100" example:"10"`
}

// JobListRequest 职位列表请求
type JobListRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1" example:"1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100" example:"10"`
}

// ToModel 将创建请求转换为模型
func (r *CreateJobRequest) ToModel() *model.Job {
	return &model.Job{
		CompanyID:     r.CompanyID,
		Name:          r.Name,
		JobSkill:      r.JobSkill,
		JobSalary:     r.JobSalary,
		JobDescribe:   r.JobDescribe,
		JobLocation:   r.JobLocation,
		JobExpireTime: r.JobExpireTime,
		JobType:       int(r.JobType),
		JobCategory:   r.JobCategory,
		JobExperience: r.JobExperience,
		JobEducation:  r.JobEducation,
		JobBenefit:    r.JobBenefit,
		JobContact:    r.JobContact,
		JobSource:     r.JobSource,
		Status:        int(enums.JobStatusNormal),
	}
}

// ToModel 将更新请求转换为模型
func (r *UpdateJobRequest) ToModel() *model.Job {
	return &model.Job{
		ID:            r.ID,
		Name:          r.Name,
		JobSkill:      r.JobSkill,
		JobSalary:     r.JobSalary,
		JobDescribe:   r.JobDescribe,
		JobLocation:   r.JobLocation,
		JobExpireTime: r.JobExpireTime,
		JobType:       int(r.JobType),
		JobCategory:   r.JobCategory,
		JobExperience: r.JobExperience,
		JobEducation:  r.JobEducation,
		JobBenefit:    r.JobBenefit,
		JobContact:    r.JobContact,
		JobSource:     r.JobSource,
	}
}

// ToConditions 将搜索请求转换为查询条件
func (r *JobSearchRequest) ToConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if r.Keyword != "" {
		conditions["keyword"] = r.Keyword
	}
	if r.JobType != 0 {
		conditions["job_type"] = r.JobType
	}
	if r.JobCategory != "" {
		conditions["job_category"] = r.JobCategory
	}
	if r.Location != "" {
		conditions["job_location"] = r.Location
	}
	if r.SalaryMin > 0 {
		conditions["salary_min"] = r.SalaryMin
	}
	if r.SalaryMax > 0 {
		conditions["salary_max"] = r.SalaryMax
	}
	if r.CompanyID > 0 {
		conditions["company_id"] = r.CompanyID
	}

	return conditions
}
