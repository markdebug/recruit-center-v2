package request

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
)

// CreateResumeRequest 创建简历请求
type CreateResumeRequest struct {
	Name            string                  `json:"name" binding:"required"`
	Avatar          string                  `json:"avatar"`
	Gender          int                     `json:"gender"`
	Birthday        time.Time               `json:"birthday"`
	Phone           string                  `json:"phone" binding:"required"`
	Email           string                  `json:"email" binding:"required,email"`
	Location        string                  `json:"location"`
	Experience      int                     `json:"experience"`
	JobStatus       int                     `json:"jobStatus"`
	ExpectedJob     string                  `json:"expectedJob"`
	ExpectedCity    string                  `json:"expectedCity"`
	Introduction    string                  `json:"introduction"`
	Skills          string                  `json:"skills"`
	Educations      []EducationRequest      `json:"educations"`
	WorkExperiences []WorkExperienceRequest `json:"workExperiences"`
	Projects        []ProjectRequest        `json:"projects"`
}

// 更新简历状态请求
type UpdateResumeStatusRequest struct {
	Status   int  `json:"status" binding:"required,oneof=1 2"`
	ResumeID uint `json:"resumeId" binding:"required"`
	UserID   uint `json:"userId" binding:"required"`
}

// ToModel 将请求转换为模型
func (r *CreateResumeRequest) ToModel(userID uint) *model.Resume {
	resume := &model.Resume{
		UserID:       userID,
		Name:         r.Name,
		Avatar:       r.Avatar,
		Gender:       r.Gender,
		Birthday:     r.Birthday,
		Phone:        r.Phone,
		Email:        r.Email,
		Location:     r.Location,
		Experience:   r.Experience,
		JobStatus:    r.JobStatus,
		ExpectedJob:  r.ExpectedJob,
		ExpectedCity: r.ExpectedCity,
		Introduction: r.Introduction,
		Skills:       r.Skills,
	}

	// 转换教育经历
	for _, edu := range r.Educations {
		resume.Educations = append(resume.Educations, *edu.ToModel())
	}

	// 转换工作经历
	for _, work := range r.WorkExperiences {
		resume.WorkExperiences = append(resume.WorkExperiences, *work.ToModel())
	}

	// 转换项目经历
	for _, proj := range r.Projects {
		resume.Projects = append(resume.Projects, *proj.ToModel())
	}

	return resume
}

type EducationRequest struct {
	School    string    `json:"school" binding:"required"`
	Major     string    `json:"major"`
	Degree    string    `json:"degree"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// ToModel 转换为模型
func (r *EducationRequest) ToModel() *model.Education {
	return &model.Education{
		School:    r.School,
		Major:     r.Major,
		Degree:    r.Degree,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	}
}

type WorkExperienceRequest struct {
	CompanyName string    `json:"companyName" binding:"required"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Achievement string    `json:"achievement"`
}

// ToModel 转换为模型
func (r *WorkExperienceRequest) ToModel() *model.WorkExperience {
	return &model.WorkExperience{
		CompanyName: r.CompanyName,
		Position:    r.Position,
		Department:  r.Department,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		Description: r.Description,
		Achievement: r.Achievement,
	}
}

type ProjectRequest struct {
	Name        string    `json:"name" binding:"required"`
	Role        string    `json:"role"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Technology  string    `json:"technology"`
	Achievement string    `json:"achievement"`
}

// ToModel 转换为模型
func (r *ProjectRequest) ToModel() *model.Project {
	return &model.Project{
		Name:        r.Name,
		Role:        r.Role,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		Description: r.Description,
		Technology:  r.Technology,
		Achievement: r.Achievement,
	}
}
