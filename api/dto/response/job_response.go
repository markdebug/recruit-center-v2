package response

import (
	"time"

	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobResponse 职位响应对象
type JobResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	CompanyID     uint      `json:"companyId"`
	JobSkill      string    `json:"jobSkill"`
	JobSalary     string    `json:"jobSalary"`
	JobDescribe   string    `json:"jobDescribe"`
	JobLocation   string    `json:"jobLocation"`
	JobExpireTime time.Time `json:"jobExpireTime"`
	Status        int       `json:"status"`
	JobType       int       `json:"jobType"`
	JobCategory   string    `json:"jobCategory"`
	JobExperience string    `json:"jobExperience"`
	JobEducation  string    `json:"jobEducation"`
	JobBenefit    string    `json:"jobBenefit"`
	JobContact    string    `json:"jobContact"`
	JobSource     string    `json:"jobSource"`
	ViewCount     int       `json:"viewCount"`
	ApplyCount    int       `json:"applyCount"`
	Priority      int       `json:"priority"`
	Tags          []string  `json:"tags"`
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
