package service

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
)

type JobStatisticsService struct {
	statsDAO *dao.JobStatisticsDAO
}

func NewJobStatisticsService(statsDAO *dao.JobStatisticsDAO) *JobStatisticsService {
	return &JobStatisticsService{
		statsDAO: statsDAO,
	}
}

// IncrementViewCount 增加浏览数
func (s *JobStatisticsService) IncrementViewCount(jobID, companyID uint) error {
	stats, err := s.statsDAO.GetByJobID(jobID)
	if err != nil {
		// 如果不存在则创建新的统计记录
		stats = &model.JobStatistics{
			JobID:     jobID,
			CompanyID: companyID,
		}
	}

	stats.ViewCount++
	stats.LastViewTime = time.Now()
	stats.ConversionRate = float64(stats.ApplyCount) / float64(stats.ViewCount)

	return s.statsDAO.SaveOrUpdate(stats)
}

// IncrementApplyCount 增加申请数
func (s *JobStatisticsService) IncrementApplyCount(jobID, companyID uint) error {
	stats, err := s.statsDAO.GetByJobID(jobID)
	if err != nil {
		stats = &model.JobStatistics{
			JobID:     jobID,
			CompanyID: companyID,
		}
	}

	stats.ApplyCount++
	stats.LastApplyTime = time.Now()
	stats.ConversionRate = float64(stats.ApplyCount) / float64(stats.ViewCount)

	return s.statsDAO.SaveOrUpdate(stats)
}
