package model

import "time"

// Resume 简历基本信息
type Resume struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"userId"`
	Name          string    `gorm:"size:50;not null" json:"name"`
	Avatar        string    `gorm:"size:255" json:"avatar"`
	Gender        int       `gorm:"default:0" json:"gender"`
	Birthday      time.Time `json:"birthday"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Email         string    `gorm:"size:100" json:"email"`
	Location      string    `gorm:"size:100" json:"location"`
	Experience    int       `json:"experience"`
	JobStatus     int       `gorm:"default:0" json:"jobStatus"`
	ExpectedJob   string    `gorm:"size:50" json:"expectedJob"`
	ExpectedCity  string    `gorm:"size:50" json:"expectedCity"`
	Introduction  string    `gorm:"type:text" json:"introduction"`
	Skills        string    `gorm:"type:text" json:"skills"`
	AccessStatus  int       `gorm:"default:2" json:"accessStatus"`  // 1: 隐藏, 2: 公开
	WorkingStatus int       `gorm:"default:1" json:"workingStatus"` // 1: 在职, 2: 离职
	Status        int       `gorm:"default:1" json:"status"`        // 1: 正常, 0: 删除
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`

	// 关联
	Educations      []Education        `json:"educations"`
	WorkExperiences []WorkExperience   `json:"workExperiences"`
	Projects        []Project          `json:"projects"`
	Attachments     []ResumeAttachment `json:"attachments"` // 新增附件关联
}

// Education 教育经历
type Education struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ResumeID  uint      `gorm:"not null;index" json:"resumeId"`
	School    string    `gorm:"size:100;not null" json:"school"`
	Major     string    `gorm:"size:100" json:"major"`
	Degree    string    `gorm:"size:50" json:"degree"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkExperience 工作经历
type WorkExperience struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ResumeID    uint      `gorm:"not null;index" json:"resumeId"`
	CompanyName string    `gorm:"size:100;not null" json:"companyName"`
	Position    string    `gorm:"size:100" json:"position"`
	Department  string    `gorm:"size:100" json:"department"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `gorm:"type:text" json:"description"`
	Achievement string    `gorm:"type:text" json:"achievement"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Project 项目经历
type Project struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ResumeID    uint      `gorm:"not null;index" json:"resumeId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Role        string    `gorm:"size:50" json:"role"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `gorm:"type:text" json:"description"`
	Technology  string    `gorm:"type:text" json:"technology"`
	Achievement string    `gorm:"type:text" json:"achievement"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ResumeAttachment 简历附件
type ResumeAttachment struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	ResumeID  uint   `gorm:"not null;index" json:"resumeId"`
	FileName  string `gorm:"size:255;not null" json:"fileName"` // 文件名
	FileURL   string `gorm:"size:1000;not null" json:"fileUrl"` // 文件URL
	FileSize  int64  `gorm:"not null" json:"fileSize"`          // 文件大小(字节)
	FileType  string `gorm:"size:50;not null" json:"fileType"`  // 文件类型(如: pdf, doc, docx等)
	Status    int    `gorm:"default:1" json:"status"`           // 状态(1: 正常, 0: 删除)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName 指定表名
func (Resume) TableName() string {
	return "t_rc_resume"
}

// TableName 指定表名
func (Education) TableName() string {
	return "t_rc_resume_education"
}

// TableName 指定表名
func (WorkExperience) TableName() string {
	return "t_rc_resume_work_experience"
}

// TableName 指定表名
func (Project) TableName() string {
	return "t_rc_resume_project"
}

// TableName 指定表名
func (ResumeAttachment) TableName() string {
	return "t_rc_resume_attachment"
}
