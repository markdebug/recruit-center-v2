package response

import "time"

// JobFavoriteStatistics 收藏统计响应
type JobFavoriteStatistics struct {
	TotalFavorites  int64   `json:"totalFavorites"`  // 收藏总数
	AverageSalary   float64 `json:"averageSalary"`   // 平均工资
	ActiveJobsCount int64   `json:"activeJobsCount"` // 活跃职位数
}

// FavoriteJobDetail 收藏职位详情
type FavoriteJobDetail struct {
	JobID      uint      `json:"jobId"`
	JobSalary  int       `json:"jobSalary"`
	SalaryMax  int       `json:"salaryMax"`
	UpdateTime time.Time `json:"updateTime"`
}
