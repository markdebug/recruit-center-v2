package response

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
)

type JobStatisticsResponse struct {
	ID             uint      `json:"id"`
	JobID          uint      `json:"jobId"`
	CompanyID      uint      `json:"companyId"`
	ViewCount      int64     `json:"viewCount"`
	ApplyCount     int64     `json:"applyCount"`
	ConversionRate float64   `json:"conversionRate"`
	LastViewTime   time.Time `json:"lastViewTime"`
	LastApplyTime  time.Time `json:"lastApplyTime"`
}

type JobStatisticsListResponse struct {
	Total   int64                   `json:"total"`
	Records []JobStatisticsResponse `json:"records"`
}

func FromJobStatistics(stats *model.JobStatistics) *JobStatisticsResponse {
	if stats == nil {
		return nil
	}
	return &JobStatisticsResponse{
		ID:             stats.ID,
		JobID:          stats.JobID,
		CompanyID:      stats.CompanyID,
		ViewCount:      stats.ViewCount,
		ApplyCount:     stats.ApplyCount,
		ConversionRate: stats.ConversionRate,
		LastViewTime:   stats.LastViewTime,
		LastApplyTime:  stats.LastApplyTime,
	}
}
