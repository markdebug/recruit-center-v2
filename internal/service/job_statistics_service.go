package service

import (
	"time"

	"org.thinkinai.com/recruit-center/api/dto/response"
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

// GetJobStatisticsByJobID 获取职位统计信息
func (s *JobStatisticsService) GetJobStatisticsByJobID(jobID uint) (*response.JobStatisticsResponse, error) {
	stats, err := s.statsDAO.GetByJobID(jobID)
	if err != nil {
		return nil, err
	}
	return response.FromJobStatistics(stats), nil
}

// GetCompanyStats 获取公司的职位统计信息
func (s *JobStatisticsService) GetCompanyStats(companyID uint, page, size int) (*response.JobStatisticsListResponse, error) {
	offset := (page - 1) * size
	stats, total, err := s.statsDAO.GetByCompanyID(companyID, offset, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobStatisticsListResponse{
		Total:   total,
		Records: make([]response.JobStatisticsResponse, len(stats)),
	}

	for i, stat := range stats {
		resp.Records[i] = *response.FromJobStatistics(&stat)
	}

	return resp, nil
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
	if stats.ApplyCount > 0 {
		stats.ConversionRate = float64(stats.ApplyCount) / float64(stats.ViewCount)
	}

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
	if stats.ViewCount > 0 {
		stats.ConversionRate = float64(stats.ApplyCount) / float64(stats.ViewCount)
	}

	return s.statsDAO.SaveOrUpdate(stats)
}
