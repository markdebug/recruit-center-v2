package service

import (
	"time"

	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
)

type JobFavoriteService struct {
	favoriteDAO *dao.JobFavoriteDAO
	jobService  *JobService
}

func NewJobFavoriteService(favoriteDAO *dao.JobFavoriteDAO, jobService *JobService) *JobFavoriteService {
	return &JobFavoriteService{
		favoriteDAO: favoriteDAO,
		jobService:  jobService,
	}
}

// AddFavorite 添加收藏
func (s *JobFavoriteService) AddFavorite(userID, jobID uint) error {
	favorite := &model.JobFavorite{
		UserID: userID,
		JobID:  jobID,
	}
	return s.favoriteDAO.Create(favorite)
}

// RemoveFavorite 取消收藏
func (s *JobFavoriteService) RemoveFavorite(userID, jobID uint) error {
	return s.favoriteDAO.Delete(userID, jobID)
}

// IsFavorited 检查是否已收藏
func (s *JobFavoriteService) IsFavorited(userID, jobID uint) (bool, error) {
	return s.favoriteDAO.IsFavorited(userID, jobID)
}

// ListFavorites 获取用户收藏的职位列表
func (s *JobFavoriteService) ListFavorites(userID uint, page, size int) (*response.JobListResponse, error) {
	favorites, total, err := s.favoriteDAO.ListByUser(userID, page, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   total,
		Records: make([]response.JobResponse, 0),
	}

	for _, fav := range favorites {
		if job, err := s.jobService.GetByID(fav.JobID, userID); err == nil {
			job.IsFavorited = true
			job.FavoriteTime = fav.CreateTime // 设置收藏时间
			resp.Records = append(resp.Records, *job)
		}
	}

	return resp, nil
}

// GetUserStatistics 获取用户收藏统计
func (s *JobFavoriteService) GetUserStatistics(userID uint) (*response.JobFavoriteStatistics, error) {
	// 获取收藏职位的原始数据
	jobs, err := s.favoriteDAO.GetUserFavoriteJobs(userID)
	if err != nil {
		return nil, err
	}

	stats := &response.JobFavoriteStatistics{
		TotalFavorites: int64(len(jobs)),
	}

	if len(jobs) > 0 {
		// 计算平均工资
		var totalSalary float64
		for _, job := range jobs {
			avgJobSalary := float64(job.JobSalary+job.SalaryMax) / 2
			totalSalary += avgJobSalary
		}
		stats.AverageSalary = totalSalary / float64(len(jobs))

		// 统计近7天更新的职位数
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		var activeCount int64
		for _, job := range jobs {
			if job.UpdateTime.After(sevenDaysAgo) {
				activeCount++
			}
		}
		stats.ActiveJobsCount = activeCount
	}

	return stats, nil
}
