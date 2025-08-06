package response

import "time"

// ResumeResponse 简历响应对象
type ResumeResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	Name          string    `json:"name"`
	Avatar        string    `json:"avatar"`
	Gender        int       `json:"gender"`
	Birthday      time.Time `json:"birthday"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Location      string    `json:"location"`
	Experience    int       `json:"experience"`
	JobStatus     int       `json:"jobStatus"`
	ExpectedJob   string    `json:"expectedJob"`
	ExpectedCity  string    `json:"expectedCity"`
	Introduction  string    `json:"introduction"`
	Skills        string    `json:"skills"`
	ShareToken    string    `json:"shareToken"`
	AccessStatus  int       `json:"accessStatus"`
	WorkingStatus int       `json:"workingStatus"`
	Status        int       `json:"status"`

	// 关联数据
	Educations      []EducationResponse      `json:"educations"`
	WorkExperiences []WorkExperienceResponse `json:"workExperiences"`
	Projects        []ProjectResponse        `json:"projects"`
	Attachments     []AttachmentResponse     `json:"attachments"`
}

// EducationResponse 教育经历响应对象
type EducationResponse struct {
	ID        uint      `json:"id"`
	ResumeID  uint      `json:"resumeId"`
	School    string    `json:"school"`
	Major     string    `json:"major"`
	Degree    string    `json:"degree"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// WorkExperienceResponse 工作经历响应对象
type WorkExperienceResponse struct {
	ID          uint      `json:"id"`
	ResumeID    uint      `json:"resumeId"`
	CompanyName string    `json:"companyName"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Achievement string    `json:"achievement"`
}

// ProjectResponse 项目经历响应对象
type ProjectResponse struct {
	ID          uint      `json:"id"`
	ResumeID    uint      `json:"resumeId"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
	Technology  string    `json:"technology"`
	Achievement string    `json:"achievement"`
}

// AttachmentResponse 简历附件响应对象
type AttachmentResponse struct {
	ID       uint   `json:"id"`
	ResumeID uint   `json:"resumeId"`
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
	Status   int    `json:"status"`
}

// ResumeListResponse 简历列表响应
type ResumeListResponse struct {
	Total   int64            `json:"total"`
	Records []ResumeResponse `json:"records"`
}
