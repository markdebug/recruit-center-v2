package service

import (
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
			resp.Records = append(resp.Records, *job)
		}
	}

	return resp, nil
}
