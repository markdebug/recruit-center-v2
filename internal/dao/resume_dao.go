package dao

import (
	"time"

	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

// ResumeDAO 简历数据访问对象
type ResumeDAO struct {
	db *gorm.DB
}

func NewResumeDAO(db *gorm.DB) *ResumeDAO {
	return &ResumeDAO{db: db}
}

// Create 创建简历
func (d *ResumeDAO) Create(resume *model.Resume) error {
	// BeforeSave 钩子会自动处理加密
	return d.db.Create(resume).Error
}

// Update 更新简历
func (d *ResumeDAO) UpdateBasic(resume *model.Resume) error {
	// 获取原有数据，用于处理部分更新场景
	oldResume, err := d.GetByID(resume.ID)
	if err != nil {
		return err
	}

	// 如果敏感字段为空，保留原有加密数据
	if resume.Phone == "" {
		resume.Phone = oldResume.Phone
	}
	if resume.Email == "" {
		resume.Email = oldResume.Email
	}
	if resume.Location == "" {
		resume.Location = oldResume.Location
	}

	resume.UpdatedAt = time.Now()
	// BeforeSave 钩子会自动处理加密
	return d.db.Save(resume).Error
}

// GetByID 获取简历详情
func (d *ResumeDAO) GetByID(id uint) (*model.Resume, error) {
	var resume model.Resume
	err := d.db.Preload("Educations").
		Preload("WorkExperiences").
		Preload("Projects").
		Preload("Attachments").
		First(&resume, id).Error
	// AfterFind 钩子会自动处理解密
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

// GetByUser 获取用户的简历
func (d *ResumeDAO) GetByUser(userID uint) (*model.Resume, error) {
	var resume model.Resume
	err := d.db.Preload("Educations").
		Preload("WorkExperiences").
		Preload("Projects").
		Preload("Attachments").
		Where("user_id = ?", userID).
		First(&resume).Error
	// AfterFind 钩子会自动处理解密
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

// 通过分享token获取简历
func (d *ResumeDAO) GetByShareToken(token string) (*model.Resume, error) {
	var resume model.Resume
	err := d.db.Preload("Educations").
		Preload("WorkExperiences").
		Preload("Projects").
		Preload("Attachments").
		Where("share_token = ? AND status = 1", token).
		First(&resume).Error
	// AfterFind 钩子会自动处理解密
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

// AddEducation 添加教育经历
func (d *ResumeDAO) AddEducation(education *model.Education) error {
	return d.db.Create(education).Error
}

// AddWorkExperience 添加工作经历
func (d *ResumeDAO) AddWorkExperience(experience *model.WorkExperience) error {
	return d.db.Create(experience).Error
}

// AddProject 添加项目经历
func (d *ResumeDAO) AddProject(project *model.Project) error {
	return d.db.Create(project).Error
}

// DeleteEducation 删除教育经历
func (d *ResumeDAO) DeleteEducation(id uint) error {
	return d.db.Delete(&model.Education{}, id).Error
}

// DeleteWorkExperience 删除工作经历
func (d *ResumeDAO) DeleteWorkExperience(id uint) error {
	return d.db.Delete(&model.WorkExperience{}, id).Error
}

// DeleteProject 删除项目经历
func (d *ResumeDAO) DeleteProject(id uint) error {
	return d.db.Delete(&model.Project{}, id).Error
}

// UpdateEducation 更新教育经历
func (d *ResumeDAO) UpdateEducation(education *model.Education) error {
	return d.db.Save(education).Error
}

// UpdateWorkExperience 更新工作经历
func (d *ResumeDAO) UpdateWorkExperience(experience *model.WorkExperience) error {
	return d.db.Save(experience).Error
}

// UpdateProject 更新项目经历
func (d *ResumeDAO) UpdateProject(project *model.Project) error {
	return d.db.Save(project).Error
}

// AddAttachment 添加简历附件
func (d *ResumeDAO) AddAttachment(attachment *model.ResumeAttachment) error {
	return d.db.Create(attachment).Error
}

// UpdateAttachment 更新简历附件
func (d *ResumeDAO) UpdateAttachment(attachment *model.ResumeAttachment) error {
	return d.db.Save(attachment).Error
}

// DeleteAttachment 删除简历附件
func (d *ResumeDAO) DeleteAttachment(id uint) error {
	return d.db.Delete(&model.ResumeAttachment{}, id).Error
}

// GetAttachments 获取简历的所有附件
func (d *ResumeDAO) GetAttachments(resumeID uint) ([]model.ResumeAttachment, error) {
	var attachments []model.ResumeAttachment
	err := d.db.Where("resume_id = ?", resumeID).Find(&attachments).Error
	return attachments, err
}

// UpdateAttachmentURL 更新简历附件的URL信息
func (d *ResumeDAO) UpdateAttachmentURL(attachmentID uint, fileURL string, fileName string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		attachment := &model.ResumeAttachment{}
		if err := tx.First(attachment, attachmentID).Error; err != nil {
			return err
		}
		// 更新附件信息
		attachment.FileURL = fileURL
		attachment.FileName = fileName
		return tx.Save(attachment).Error
	})
}

// List 获取简历列表
func (d *ResumeDAO) List(page, size int) ([]model.Resume, int64, error) {
	var resumes []model.Resume
	var total int64

	if err := d.db.Model(&model.Resume{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := d.db.Preload("Educations").
		Preload("WorkExperiences").
		Preload("Projects").
		Preload("Attachments").
		Offset(offset).
		Limit(size).
		Find(&resumes).Error
	// AfterFind 钩子会自动处理解密

	return resumes, total, err
}
