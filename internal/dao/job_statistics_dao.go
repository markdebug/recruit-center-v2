package dao

import (
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

func (dao *JobStatisticsDAO) SaveOrUpdate(stats *model.JobStatistics) error {
	return dao.db.Save(stats).Error
}

// GetCompanyStats 获取公司所有职位的统计信息
func (dao *JobStatisticsDAO) GetCompanyStats(companyID uint) ([]model.JobStatistics, error) {
	var stats []model.JobStatistics
	err := dao.db.Where("company_id = ?", companyID).Find(&stats).Error
	return stats, err
}
