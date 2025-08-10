package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobResponse 职位响应对象
// @Description 职位响应对象
type JobResponse struct {
	ID            uint                   `json:"id"`                                                               // id
	Name          string                 `json:"name"`                                                             // 职位名称
	CompanyID     uint                   `json:"companyId"`                                                        // 公司ID
	JobSkill      string                 `json:"jobSkill"`                                                         // 职位技能要求
	JobSalary     int                    `json:"jobSalary"`                                                        // 职位薪资最低
	JobSalaryMax  int                    `json:"jobSalaryMax"`                                                     // 职位薪资最高
	JobDescribe   string                 `json:"jobDescribe"`                                                      // 职位描述
	JobLocation   string                 `json:"jobLocation"`                                                      // 工作地点
	JobExpireTime time.Time              `json:"jobExpireTime"`                                                    // 职位过期时间
	Status        int                    `json:"status" example:"1" enums:"0,1"`                                   // 职位状态 0: 禁用 1: 正常
	JobType       int                    `json:"jobType" example:"1" enums:"1,2,3"`                                // 职位类型 1 全职 2 兼职 3 实习
	RemoteType    int                    `json:"remoteType" example:"1" enums:"1,2,3,4"`                           // 远程办公类型
	RemoteDesc    string                 `json:"remoteDesc"`                                                       // 远程办公描述
	RemoteRatio   int                    `json:"remoteRatio" example:"60" minimum:"0" maximum:"100"`               // 远程办公比例(0-100)
	JobCategory   string                 `json:"jobCategory" example:"direct" enums:"internal,direct,headhunting"` // 职位分类 internal 内推 2 direct 直招 3 猎头
	JobExperience string                 `json:"jobExperience"`                                                    // 经验要求
	JobEducation  string                 `json:"jobEducation"`                                                     // 学历要求
	JobBenefit    string                 `json:"jobBenefit"`                                                       // 职位福利
	JobContact    string                 `json:"jobContact"`                                                       // 联系方式
	JobSource     string                 `json:"jobSource"`                                                        // 职位来源
	DeleteStatus  int                    `json:"deleteStatus" example:"0" enums:"0,1"`                             // 删除状态 0: 正常 1: 已删除
	ViewCount     int                    `json:"viewCount"`                                                        // 浏览次数
	ApplyCount    int                    `json:"applyCount"`                                                       // 申请次数
	Priority      int                    `json:"priority"`                                                         // 优先级
	Tags          []string               `json:"tags"`
	Benefits      []model.JobBenefitType `json:"benefits" `
	BenefitDesc   string                 `json:"benefitDesc"`
	BenefitTexts  []string               `json:"benefitTexts"`
	// 是否已收藏
	IsFavorited bool `json:"isFavorited"`
	// 收藏时间
	FavoriteTime time.Time `json:"favoriteTime,omitempty"` // 收藏时间，可能为空
	// 是否已投递
	IsApplied bool `json:"isApplied"`
}

// JobListResponse 职位列表响应
type JobListResponse struct {
	Total   int64         `json:"total"`   // 总数
	Records []JobResponse `json:"records"` // 公司信息
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
