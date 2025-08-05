package dao

import (
	"time"

	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

// JobDAO 职位数据访问对象
type JobDAO struct {
	db *gorm.DB
}

// NewJobDAO 创建职位DAO实例
func NewJobDAO(db *gorm.DB) *JobDAO {
	return &JobDAO{db: db}
}

// Create 创建职位
func (d *JobDAO) Create(job *model.Job) error {
	return d.db.Create(job).Error
}

// Update 更新职位
func (d *JobDAO) Update(job *model.Job) error {
	return d.db.Save(job).Error
}

// Delete 删除职位
func (d *JobDAO) Delete(id uint) error {
	return d.db.Model(&model.Job{}).Where("id = ?", id).Update("delete_status", 1).Error
}

// GetByID 根据ID获取职位
func (d *JobDAO) GetByID(id uint) (*model.Job, error) {
	var job model.Job
	err := d.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// List 获取职位列表
func (d *JobDAO) List(page, size int) ([]model.Job, int64, error) {
	var jobs []model.Job
	var total int64

	if err := d.db.Model(&model.Job{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := d.db.Offset(offset).Limit(size).Find(&jobs).Error
	return jobs, total, err
}

// 只保留特定的业务方法，通用的CRUD方法继承自baseDAO
func (d *JobDAO) SearchByKeyword(keyword string) ([]model.Job, error) {
	var jobs []model.Job
	err := d.db.Where("name LIKE ? OR job_skill LIKE ? OR job_describe LIKE ? AND job_expire_time > now() AND delete_status = 0",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&jobs).Error
	return jobs, err
}

// GetActiveJobs 获取未过期的职位
func (d *JobDAO) GetActiveJobs() ([]model.Job, error) {
	var jobs []model.Job
	err := d.db.Where("job_expire_time > ? and delete_status=0", time.Now()).Find(&jobs).Error
	return jobs, err
}

// 获取指定类型且有效期内的职位
func (d *JobDAO) GetJobsByType(jobType string) ([]model.Job, error) {
	var jobs []model.Job
	err := d.db.Where("job_type = ? AND job_expire_time > now() and delete_status=0", jobType).Find(&jobs).Error
	return jobs, err
}

// 根据公司id统计该公司已经发布的岗位总数
func (d *JobDAO) CountJobsByCompany(companyID uint) (int64, error) {
	var count int64
	err := d.db.Model(&model.Job{}).Where("company_id = ? AND job_expire_time > now() and delete_status=0", companyID).Count(&count).Error
	return count, err
}

// SearchByCompany 根据公司ID查询有效职位
func (d *JobDAO) SearchByCompany(companyID uint, page, pageSize int) ([]model.Job, int64, error) {
	var jobs []model.Job
	var total int64

	query := d.db.Model(&model.Job{}).Where("company_id = ? AND delete_status = 0 AND job_expire_time > ?",
		companyID, time.Now())

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobs).Error

	return jobs, total, err
}

// SearchByCondition 多条件搜索职位
func (d *JobDAO) SearchByCondition(conditions map[string]interface{}, page, pageSize int) ([]model.Job, int64, error) {
	var jobs []model.Job
	var total int64

	// 添加基础条件
	query := d.db.Model(&model.Job{}).Where("delete_status = 0 AND job_expire_time > ?", time.Now())

	// 添加搜索条件
	for key, value := range conditions {
		switch key {
		case "keyword":
			query = query.Where("name LIKE ? OR job_skill LIKE ? OR job_describe LIKE ?",
				"%"+value.(string)+"%", "%"+value.(string)+"%", "%"+value.(string)+"%")
		case "job_type":
			query = query.Where("job_type = ?", value)
		case "job_category":
			query = query.Where("job_category = ?", value)
		case "job_location":
			query = query.Where("job_location LIKE ?", "%"+value.(string)+"%")
		case "salary_min":
			query = query.Where("CAST(SUBSTRING_INDEX(job_salary, '-', 1) AS SIGNED) >= ?", value)
		case "salary_max":
			query = query.Where("CAST(SUBSTRING_INDEX(job_salary, '-', -1) AS SIGNED) <= ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.Order("create_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobs).Error

	return jobs, total, err
}

// BatchUpdate 批量更新职位状态
func (d *JobDAO) BatchUpdate(ids []uint, updates map[string]interface{}) error {
	return d.db.Model(&model.Job{}).Where("id IN ?", ids).Updates(updates).Error
}

// UpdateStatus 更新职位状态
func (d *JobDAO) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.Job{}).Where("id = ?", id).Update("status", status).Error
}

// GetExpiredJobs 获取已过期职位
func (d *JobDAO) GetExpiredJobs() ([]model.Job, error) {
	var jobs []model.Job
	err := d.db.Where("job_expire_time <= ? AND status = 1 AND delete_status = 0",
		time.Now()).Find(&jobs).Error
	return jobs, err
}
