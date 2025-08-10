package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobResponse 职位响应对象
// @Description 职位响应对象
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
	Status        int                    `json:"status" example:"1" enums:"0,1"`
	JobType       int                    `json:"jobType" example:"1" enums:"1,2,3"`                                // 职位类型 1 全职 2 兼职 3 实习
	RemoteType    int                    `json:"remoteType" example:"1" enums:"1,2,3,4"`                           // 远程办公类型
	RemoteDesc    string                 `json:"remoteDesc"`                                                       // 远程办公描述
	RemoteRatio   int                    `json:"remoteRatio" example:"60" minimum:"0" maximum:"100"`               // 远程办公比例(0-100)
	JobCategory   string                 `json:"jobCategory" example:"direct" enums:"internal,direct,headhunting"` // 职位分类 internal 内推 2 direct 直招 3 猎头
	JobExperience string                 `json:"jobExperience"`
	JobEducation  string                 `json:"jobEducation"`
	JobBenefit    string                 `json:"jobBenefit"`
	JobContact    string                 `json:"jobContact"`
	JobSource     string                 `json:"jobSource"`                            // 职位来源
	DeleteStatus  int                    `json:"deleteStatus" example:"0" enums:"0,1"` // 删除状态 0: 正常 1: 已删除
	ViewCount     int                    `json:"viewCount"`
	ApplyCount    int                    `json:"applyCount"`
	Priority      int                    `json:"priority"`
	Tags          []string               `json:"tags"`
	Benefits      []model.JobBenefitType `json:"benefits" `
	BenefitDesc   string                 `json:"benefitDesc"`
	BenefitTexts  []string               `json:"benefitTexts"`
	// 是否已收藏
	IsFavorited bool `json:"isFavorited"`
	// 是否已投递
	IsApplied bool `json:"isApplied"`
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
