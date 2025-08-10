package request

import (
	"errors"
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// CreateJobRequest 创建职位请求
// @Description 创建职位的请求参数
type CreateJobRequest struct {
	// 职位名称
	// required: true
	Name string `json:"name" binding:"required" example:"高级Go工程师"`

	// 公司ID
	// required: true
	CompanyID uint `json:"companyId" binding:"required"`

	// 职位技能要求
	// required: true
	JobSkill string `json:"jobSkill" binding:"required" example:"Go,Docker,Kubernetes"`

	// 职位薪资范围
	// required: true
	JobSalary int `json:"jobSalary" binding:"required,min=1" example:"15000"`

	// 职位最高薪资
	// required: true
	JobSalaryMax int `json:"jobSalaryMax" binding:"required,gtefield=JobSalary" example:"25000"`

	// 职位描述
	// required: true
	JobDescribe string `json:"jobDescribe" binding:"required"`

	// 工作地点
	// required: true
	JobLocation string `json:"jobLocation" binding:"required" example:"上海"`

	// 职位过期时间
	// required: true
	// example: 2024-12-31T23:59:59Z
	JobExpireTime time.Time `json:"jobExpireTime" binding:"required"`

	// 职位类型（全职/兼职）
	// required: true
	JobType int `json:"jobType" binding:"required" example:"1"`

	// 职位类别
	// required: true
	JobCategory string `json:"jobCategory" binding:"required" example:"技术"`

	// 经验要求
	// required: true
	JobExperience string `json:"jobExperience" binding:"required" example:"3-5年"`

	// 学历要求
	// required: true
	JobEducation string `json:"jobEducation" binding:"required" example:"本科"`

	// 职位福利
	JobBenefit string `json:"jobBenefit" example:"五险一金,年终奖,带薪休假"`

	// 联系方式
	JobContact string `json:"jobContact"`

	// 职位来源
	JobSource string `json:"jobSource" example:"官网"`

	// 职位标签
	Tags []string `json:"tags" example:"['急招','环境好']"`

	// 远程办公类型(1:办公室办公 2:混合办公 3:全远程 4:灵活办公)
	// required: true
	RemoteType int `json:"remoteType" binding:"required,oneof=1 2 3 4" example:"2"`

	// 远程办公补充说明
	RemoteDesc string `json:"remoteDesc" example:"每周可远程办公3天,周四需到办公室参加团队会议"`

	// 远程办公比例(0-100)
	// required: true
	RemoteRatio int `json:"remoteRatio" binding:"required,min=0,max=100" example:"60"`

	// 职位福利列表
	// required: true
	Benefits []model.JobBenefitType `json:"benefits" binding:"required,dive,oneof=1 2 3 4 5 6 7 8 9 10""`

	// 福利补充说明
	BenefitDesc string `json:"benefitDesc" example:"额外提供商业医疗保险，每年体检一次"`
}

// UpdateJobRequest 更新职位请求
// @Description 更新职位的请求参数
type UpdateJobRequest struct {
	ID            uint          `json:"id" binding:"required"`
	Name          string        `json:"name" binding:"min=1,max=100"`
	JobSkill      string        `json:"jobSkill" binding:"max=500"`
	JobSalary     *int          `json:"jobSalary,omitempty" binding:"omitempty,min=1"`
	JobSalaryMax  *int          `json:"jobSalaryMax,omitempty" binding:"omitempty,gtefield=JobSalary"`
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

	// 远程办公类型(1:办公室办公 2:混合办公 3:全远程 4:灵活办公)
	RemoteType int `json:"remoteType,omitempty" binding:"omitempty,oneof=1 2 3 4"`

	// 远程办公补充说明
	RemoteDesc string `json:"remoteDesc,omitempty"`

	// 远程办公比例(0-100)
	RemoteRatio int `json:"remoteRatio,omitempty" binding:"omitempty,min=0,max=100"`

	// 职位福利列表
	Benefits []model.JobBenefitType `json:"benefits,omitempty" binding:"omitempty,dive,oneof=1 2 3 4 5 6 7 8 9 10"`

	// 福利补充说明
	BenefitDesc string `json:"benefitDesc,omitempty"`
}

// JobSearchRequest 职位搜索请求
// @Description 职位搜索的请求参数
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
// @Description 职位列表的请求参数
type JobListRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1" example:"1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100" example:"10"`
}

// Validate 请求验证
func (r *CreateJobRequest) Validate() error {
	if r.JobSalaryMax < r.JobSalary {
		return errors.New("最高薪资不能低于最低薪资")
	}
	return nil
}

// ToModel 将创建请求转换为模型
func (r *CreateJobRequest) ToModel() *model.Job {
	return &model.Job{
		CompanyID:     r.CompanyID,
		Name:          r.Name,
		JobSkill:      r.JobSkill,
		JobSalary:     r.JobSalary,
		JobSalaryMax:  r.JobSalaryMax,
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
		RemoteType:    model.RemoteType(r.RemoteType),
		RemoteDesc:    r.RemoteDesc,
		RemoteRatio:   r.RemoteRatio,
		Benefits:      r.Benefits,
		BenefitDesc:   r.BenefitDesc,
	}
}

// Validate 请求验证
func (r *UpdateJobRequest) Validate() error {
	if r.JobSalary != nil && r.JobSalaryMax != nil && *r.JobSalaryMax < *r.JobSalary {
		return errors.New("最高薪资不能低于最低薪资")
	}
	return nil
}

// ToModel 将更新请求转换为模型
func (r *UpdateJobRequest) ToModel() *model.Job {
	job := &model.Job{
		ID:            r.ID,
		Name:          r.Name,
		JobSkill:      r.JobSkill,
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
		RemoteType:    model.RemoteType(r.RemoteType),
		RemoteDesc:    r.RemoteDesc,
		RemoteRatio:   r.RemoteRatio,
		Benefits:      r.Benefits,
		BenefitDesc:   r.BenefitDesc,
	}

	// 处理可选的薪资字段
	if r.JobSalary != nil {
		job.JobSalary = *r.JobSalary
	}
	if r.JobSalaryMax != nil {
		job.JobSalaryMax = *r.JobSalaryMax
	}

	return job
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
