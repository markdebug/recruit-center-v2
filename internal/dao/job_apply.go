package dao

import (
	"time"

	"gorm.io/gorm"
)

// JobApply 职位申请记录
type JobApply struct {
	ID            uint      `gorm:"primarykey" json:"id"`              // 申请ID
	JobID         uint      `gorm:"not null;index" json:"job_id"`      // 职位ID
	UserID        uint      `gorm:"not null;index" json:"user_id"`     // 用户ID
	ApplyTime     time.Time `gorm:"not null" json:"apply_time"`        // 申请时间
	ApplyProgress string    `gorm:"size:50" json:"apply_progress"`     // 申请进度
	Status        int       `gorm:"default:1" json:"status"`           // 申请状态
	CreateTime    time.Time `gorm:"autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"update_time"` // 更新时间
}

// TableName 指定表名
func (JobApply) TableName() string {
	return "t_rc_job_apply"
}

// JobApplyDao 职位申请数据访问对象
type JobApplyDao struct {
	db *gorm.DB
}

// NewJobApplyDao 创建职位申请DAO实例
func NewJobApplyDao(db *gorm.DB) *JobApplyDao {
	return &JobApplyDao{db: db}
}

// Create 创建职位申请记录
func (d *JobApplyDao) Create(apply *JobApply) error {
	apply.ApplyTime = time.Now()
	return d.db.Create(apply).Error
}

// Update 更新职位申请记录
func (d *JobApplyDao) Update(apply *JobApply) error {
	return d.db.Save(apply).Error
}

// Delete 删除职位申请记录
func (d *JobApplyDao) Delete(id uint) error {
	return d.db.Delete(&JobApply{}, id).Error
}

// GetByID 根据ID获取申请记录
func (d *JobApplyDao) GetByID(id uint) (*JobApply, error) {
	var apply JobApply
	err := d.db.First(&apply, id).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

// GetByUserAndJob 获取用户对特定职位的申请记录
func (d *JobApplyDao) GetByUserAndJob(userID, jobID uint) (*JobApply, error) {
	var apply JobApply
	err := d.db.Where("user_id = ? AND job_id = ?", userID, jobID).First(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

// ListByUser 获取用户的所有申请记录
func (d *JobApplyDao) ListByUser(userID uint, page, pageSize int) ([]JobApply, int64, error) {
	var applies []JobApply
	var total int64

	if err := d.db.Model(&JobApply{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := d.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("apply_time DESC").
		Find(&applies).Error; err != nil {
		return nil, 0, err
	}

	return applies, total, nil
}

// ListByJob 获取职位的所有申请记录
func (d *JobApplyDao) ListByJob(jobID uint, page, pageSize int) ([]JobApply, int64, error) {
	var applies []JobApply
	var total int64

	if err := d.db.Model(&JobApply{}).Where("job_id = ?", jobID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := d.db.Where("job_id = ?", jobID).
		Offset(offset).
		Limit(pageSize).
		Order("apply_time DESC").
		Find(&applies).Error; err != nil {
		return nil, 0, err
	}

	return applies, total, nil
}

// UpdateProgress 更新申请进度
func (d *JobApplyDao) UpdateProgress(id uint, progress string) error {
	return d.db.Model(&JobApply{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"apply_progress": progress,
			"update_time":    time.Now(),
		}).Error
}

// UpdateStatus 更新申请状态
func (d *JobApplyDao) UpdateStatus(id uint, status int) error {
	return d.db.Model(&JobApply{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}
