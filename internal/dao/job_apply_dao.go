package dao

import (
	"time"

	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

// JobApplyDAO 职位申请数据访问对象
type JobApplyDAO struct {
	db *gorm.DB
}

// NewJobApplyDAO 创建职位申请DAO实例
func NewJobApplyDAO(db *gorm.DB) *JobApplyDAO {
	return &JobApplyDAO{db: db}
}

// Create 创建职位申请记录
func (d *JobApplyDAO) Create(apply *model.JobApply) error {
	apply.ApplyTime = time.Now()
	return d.db.Create(apply).Error
}

// Update 更新职位申请记录
func (d *JobApplyDAO) Update(apply *model.JobApply) error {
	return d.db.Save(apply).Error
}

// Delete 删除职位申请记录
func (d *JobApplyDAO) Delete(id uint) error {
	return d.db.Delete(&model.JobApply{}, id).Error
}

// GetByID 根据ID获取申请记录
func (d *JobApplyDAO) GetByID(id uint) (*model.JobApply, error) {
	var apply model.JobApply
	err := d.db.First(&apply, id).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

// List 获取申请列表
func (d *JobApplyDAO) List(page, size int) ([]model.JobApply, int64, error) {
	var applies []model.JobApply
	var total int64

	if err := d.db.Model(&model.JobApply{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := d.db.Offset(offset).Limit(size).Find(&applies).Error
	return applies, total, err
}

// 只保留特定的业务方法
func (d *JobApplyDAO) GetByUserAndJob(userID, jobID uint) (*model.JobApply, error) {
	var apply model.JobApply
	err := d.db.Where("user_id = ? AND job_id = ?", userID, jobID).First(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

// ListByUser 获取用户的所有申请记录
func (d *JobApplyDAO) ListByUser(userID uint, page, size int) ([]model.JobApply, int64, error) {
	var applies []model.JobApply
	var total int64

	if err := d.db.Model(&model.JobApply{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := d.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(size).
		Order("apply_time DESC").
		Find(&applies).Error; err != nil {
		return nil, 0, err
	}

	return applies, total, nil
}

// ListByCompany 获取公司所有的职位申请记录
func (d *JobApplyDAO) ListByCompany(companyID uint, page, size int) ([]model.JobApply, int64, error) {
	var applies []model.JobApply
	var total int64

	if err := d.db.Model(&model.JobApply{}).Where("company_id = ?", companyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := d.db.Where("company_id = ?", companyID).
		Offset(offset).
		Limit(size).
		Order("apply_time DESC").
		Find(&applies).Error; err != nil {
		return nil, 0, err
	}

	return applies, total, nil
}

// ListByJob 获取职位的所有申请记录
func (d *JobApplyDAO) ListByJob(jobID uint, page, size int) ([]model.JobApply, int64, error) {
	var applies []model.JobApply
	var total int64

	if err := d.db.Model(&model.JobApply{}).Where("job_id = ?", jobID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := d.db.Where("job_id = ?", jobID).
		Offset(offset).
		Limit(size).
		Order("apply_time DESC").
		Find(&applies).Error; err != nil {
		return nil, 0, err
	}

	return applies, total, nil
}

// UpdateProgress 更新申请进度
func (d *JobApplyDAO) UpdateProgress(id uint, progress string) error {
	return d.db.Model(&model.JobApply{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"apply_progress": progress,
			"update_time":    time.Now(),
		}).Error
}

// UpdateStatus 更新申请状态
func (d *JobApplyDAO) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.JobApply{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}
