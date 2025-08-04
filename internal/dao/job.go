package dao

import (
	"time"

	"gorm.io/gorm"
)

// Job 职位信息结构体
type Job struct {
	ID            uint      `gorm:"primarykey" json:"id"`              // 职位ID
	Name          string    `gorm:"size:100;not null" json:"name"`     // 职位名称
	CompanyID     uint      `gorm:"not null" json:"company_id"`        // 公司ID
	JobSkill      string    `gorm:"size:500" json:"job_skill"`         // 职位所需技能
	JobSalary     string    `gorm:"size:50" json:"job_salary"`         // 职位薪资范围
	JobDescribe   string    `gorm:"type:text" json:"job_describe"`     // 职位描述
	JobLocation   string    `gorm:"size:200" json:"job_location"`      // 工作地点
	JobExpireTime time.Time `gorm:"index" json:"job_expire_time"`      // 职位过期时间
	Status        int       `gorm:"default:1" json:"status"`           // 职位状态（1: 正常, 0: 已过期）
	JobType       string    `gorm:"size:50" json:"job_type"`           // 职位类型（全职、兼职、实习等）
	JobCategory   string    `gorm:"size:50" json:"job_category"`       // 职位分类 （行业分类）
	JobExperience string    `gorm:"size:50" json:"job_experience"`     // 职位经验要求
	JobEducation  string    `gorm:"size:50" json:"job_education"`      // 职位学历要求
	JobBenefit    string    `gorm:"size:500" json:"job_benefit"`       // 职位福利
	JobContact    string    `gorm:"size:100" json:"job_contact"`       // 联系方式
	DeleteStatus  int       `gorm:"default:0" json:"delete_status"`    // 删除状态（0: 正常, 1: 已删除）
	JobSource     string    `gorm:"size:100" json:"job_source"`        // 职位来源（内推、直招、猎头等）
	CreateTime    time.Time `gorm:"autoCreateTime" json:"create_time"` // 创建时间
	UpdateTime    time.Time `gorm:"autoUpdateTime" json:"update_time"` // 更新时间
}

// TableName 指定表名
func (Job) TableName() string {
	return "t_rc_job"
}

// JobDao 职位数据访问对象
type JobDao struct {
	db *gorm.DB
}

// NewJobDao 创建职位DAO实例
func NewJobDao(db *gorm.DB) *JobDao {
	return &JobDao{db: db}
}

// SearchByKeyword 根据关键词搜索职位
func (d *JobDao) SearchByKeyword(keyword string) ([]Job, error) {
	var jobs []Job
	err := d.db.Where("name LIKE ? OR job_skill LIKE ? OR job_describe LIKE ? AND job_expire_time > now() and delete_status=0",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&jobs).Error
	return jobs, err
}

// GetActiveJobs 获取未过期的职位
func (d *JobDao) GetActiveJobs() ([]Job, error) {
	var jobs []Job
	err := d.db.Where("job_expire_time > ? and delete_status=0", time.Now()).Find(&jobs).Error
	return jobs, err
}

// 获取指定类型且有效期内的职位
func (d *JobDao) GetJobsByType(jobType string) ([]Job, error) {
	var jobs []Job
	err := d.db.Where("job_type = ? AND job_expire_time > now() and delete_status=0", jobType).Find(&jobs).Error
	return jobs, err
}

// 根据公司id统计该公司已经发布的岗位总数
func (d *JobDao) CountJobsByCompany(companyID uint) (int64, error) {
	var count int64
	err := d.db.Model(&Job{}).Where("company_id = ? AND job_expire_time > now() and delete_status=0", companyID).Count(&count).Error
	return count, err
}

// Create 创建职位
func (d *JobDao) Create(job *Job) error {
	// 设置默认值
	if job.Status == 0 {
		job.Status = 1 // 默认状态为正常
	}
	if job.JobExpireTime.IsZero() {
		job.JobExpireTime = time.Now().AddDate(0, 1, 0) // 默认过期时间为1个月
	}
	return d.db.Create(job).Error
}

// Update 更新职位信息
func (d *JobDao) Update(job *Job) error {
	return d.db.Model(&Job{}).Where("id = ? AND delete_status = 0", job.ID).Updates(job).Error
}

// Delete 软删除职位
func (d *JobDao) Delete(id uint) error {
	return d.db.Model(&Job{}).Where("id = ?", id).Update("delete_status", 1).Error
}

// GetByID 根据ID获取未删除的职位
func (d *JobDao) GetByID(id uint) (*Job, error) {
	var job Job
	err := d.db.Where("id = ? AND delete_status = 0", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// List 获取职位列表(分页)
func (d *JobDao) List(page, pageSize int) ([]Job, int64, error) {
	var jobs []Job
	var total int64

	// 获取未删除职位的总数
	if err := d.db.Model(&Job{}).Where("delete_status = 0").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := d.db.Where("delete_status = 0").
		Order("create_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobs).Error

	return jobs, total, err
}

// SearchByCompany 根据公司ID查询有效职位
func (d *JobDao) SearchByCompany(companyID uint, page, pageSize int) ([]Job, int64, error) {
	var jobs []Job
	var total int64

	query := d.db.Model(&Job{}).Where("company_id = ? AND delete_status = 0 AND job_expire_time > ?",
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
func (d *JobDao) SearchByCondition(conditions map[string]interface{}, page, pageSize int) ([]Job, int64, error) {
	var jobs []Job
	var total int64

	// 添加基础条件
	query := d.db.Model(&Job{}).Where("delete_status = 0 AND job_expire_time > ?", time.Now())

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
func (d *JobDao) BatchUpdate(ids []uint, updates map[string]interface{}) error {
	return d.db.Model(&Job{}).Where("id IN ?", ids).Updates(updates).Error
}

// UpdateStatus 更新职位状态
func (d *JobDao) UpdateStatus(id uint, status int) error {
	return d.db.Model(&Job{}).Where("id = ?", id).Update("status", status).Error
}

// GetExpiredJobs 获取已过期职位
func (d *JobDao) GetExpiredJobs() ([]Job, error) {
	var jobs []Job
	err := d.db.Where("job_expire_time <= ? AND status = 1 AND delete_status = 0",
		time.Now()).Find(&jobs).Error
	return jobs, err
}
