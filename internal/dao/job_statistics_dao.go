package dao

import (
	"time"

	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

type JobStatisticsDAO struct {
	db *gorm.DB
}

func NewJobStatisticsDAO(db *gorm.DB) *JobStatisticsDAO {
	return &JobStatisticsDAO{db: db}
}

func (dao *JobStatisticsDAO) GetByJobID(jobID uint) (*model.JobStatistics, error) {
	var stats model.JobStatistics
	err := dao.db.Where("job_id = ?", jobID).First(&stats).Error
	return &stats, err
}

// SaveOrUpdate 保存或更新职位统计信息
func (dao *JobStatisticsDAO) SaveOrUpdate(stats *model.JobStatistics) error {
	// 检查记录是否存在
	var count int64
	err := dao.db.Model(&model.JobStatistics{}).
		Where("job_id = ?", stats.JobID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		// 记录不存在,执行插入
		return dao.db.Create(stats).Error
	}

	// 记录存在,执行原子更新
	result := dao.db.Model(&model.JobStatistics{}).
		Where("job_id = ?", stats.JobID).
		Updates(map[string]interface{}{
			"view_count":      gorm.Expr("view_count + ?", stats.ViewCount),
			"apply_count":     gorm.Expr("apply_count + ?", stats.ApplyCount),
			"conversion_rate": gorm.Expr("CAST(apply_count AS FLOAT) / NULLIF(view_count, 0)"),
			"last_view_time":  stats.LastViewTime,
			"last_apply_time": stats.LastApplyTime,
			"update_time":     time.Now(),
		})

	return result.Error
}

// GetCompanyStats 获取公司所有职位的统计信息
func (dao *JobStatisticsDAO) GetCompanyStats(companyID uint) ([]model.JobStatistics, error) {
	var stats []model.JobStatistics
	err := dao.db.Where("company_id = ?", companyID).Find(&stats).Error
	return stats, err
}

// GetByCompanyID 获取公司的职位统计列表
func (dao *JobStatisticsDAO) GetByCompanyID(companyID uint, offset, limit int) ([]model.JobStatistics, int64, error) {
	var stats []model.JobStatistics
	var total int64

	// 查询总数
	if err := dao.db.Model(&model.JobStatistics{}).Where("company_id = ?", companyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	err := dao.db.Where("company_id = ?", companyID).
		Order("view_count desc").
		Offset(offset).
		Limit(limit).
		Find(&stats).Error

	return stats, total, err
}
