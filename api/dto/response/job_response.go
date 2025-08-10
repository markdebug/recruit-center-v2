package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobResponse 职位响应对象
type JobResponse struct {
	ID            uint                   `json:"id"`
	Name          string                 `json:"name"`
	CompanyID     uint                   `json:"companyId"`
	JobSkill      string                 `json:"jobSkill"`
	JobSalary     int                    `json:"jobSalary"`
	JobSalaryMax  int                    `json:"jobSalaryMax"`
	JobDescribe   string                 `json:"jobDescribe"`
	JobLocation   string                 `json:"jobLocation"`
	JobExpireTime time.Time              `json:"jobExpireTime"`
	Status        int                    `json:"status"`
	JobType       int                    `json:"jobType"`
	JobCategory   string                 `json:"jobCategory"`
	JobExperience string                 `json:"jobExperience"`
	JobEducation  string                 `json:"jobEducation"`
	JobBenefit    string                 `json:"jobBenefit"`
	JobContact    string                 `json:"jobContact"`
	JobSource     string                 `json:"jobSource"`
	ViewCount     int                    `json:"viewCount"`
	ApplyCount    int                    `json:"applyCount"`
	Priority      int                    `json:"priority"`
	Tags          []string               `json:"tags"`
	Benefits      []model.JobBenefitType `json:"benefits"`
	BenefitDesc   string                 `json:"benefitDesc"`
	BenefitTexts  []string               `json:"benefitTexts"`
}

// JobListResponse 职位列表响应
type JobListResponse struct {
	Total   int64         `json:"total"`
	Records []JobResponse `json:"records"`
}

// 新增获取job是否激活的方法
func (j *JobResponse) IsActive() bool {
	return j.Status == int(enums.JobStatusNormal)
}

func FromJob(job *model.Job) *JobResponse {
	resp := &JobResponse{
		ID:            job.ID,
		Name:          job.Name,
		CompanyID:     job.CompanyID,
		JobSkill:      job.JobSkill,
		JobSalary:     job.JobSalary,
		JobSalaryMax:  job.JobSalaryMax,
		JobDescribe:   job.JobDescribe,
		JobLocation:   job.JobLocation,
		JobExpireTime: job.JobExpireTime,
		Status:        job.Status,
		JobType:       job.JobType,
		JobCategory:   job.JobCategory,
		JobExperience: job.JobExperience,
		JobEducation:  job.JobEducation,
		JobBenefit:    job.JobBenefit,
		JobContact:    job.JobContact,
		JobSource:     job.JobSource,
		ViewCount:     job.ViewCount,
		ApplyCount:    job.ApplyCount,
		Priority:      job.Priority,
		Tags:          job.Tags,
		Benefits:      job.Benefits,
		BenefitDesc:   job.BenefitDesc,
	}

	// 转换福利为文本描述
	resp.BenefitTexts = make([]string, len(job.Benefits))
	for i, benefit := range job.Benefits {
		resp.BenefitTexts[i] = job.GetBenefitText(benefit)
	}

	return resp
}
