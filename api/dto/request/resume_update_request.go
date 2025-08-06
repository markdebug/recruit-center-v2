package request

import "time"

// UpdateResumeModuleType 定义更新模块类型
type UpdateResumeModuleType string

const (
	ModuleBasic          UpdateResumeModuleType = "basic"     // 基本信息
	ModuleEducation      UpdateResumeModuleType = "education" // 教育经历
	ModuleWorkExperience UpdateResumeModuleType = "work"      // 工作经历
	ModuleProject        UpdateResumeModuleType = "project"   // 项目经历
	ModuleSkills         UpdateResumeModuleType = "skills"    // 技能特长
)

// UpdateResumeRequest 更新简历请求
type UpdateResumeRequest struct {
	Module UpdateResumeModuleType `json:"module" binding:"required"` // 更新模块
	Data   interface{}            `json:"data" binding:"required"`   // 更新数据
}

// UpdateResumeBasicRequest 更新基本信息
type UpdateResumeBasicRequest struct {
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	Gender       int       `json:"gender"`
	Birthday     time.Time `json:"birthday"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Location     string    `json:"location"`
	Introduction string    `json:"introduction"`
}

// UpdateResumeEducationRequest 更新教育经历
type UpdateResumeEducationRequest struct {
	ID        uint      `json:"id"` // 0表示新增
	School    string    `json:"school"`
	Major     string    `json:"major"`
	Degree    string    `json:"degree"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// UpdateResumeWorkRequest 更新工作经历
type UpdateResumeWorkRequest struct {
	ID          uint      `json:"id"` // 0表示新增
	CompanyName string    `json:"companyName"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Achievement string    `json:"achievement"`
}
