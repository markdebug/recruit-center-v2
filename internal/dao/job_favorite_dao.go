package dao

import (
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/model"
)

type JobFavoriteDAO struct {
	db *gorm.DB
}

func NewJobFavoriteDAO(db *gorm.DB) *JobFavoriteDAO {
	return &JobFavoriteDAO{db: db}
}

func (dao *JobFavoriteDAO) Create(favorite *model.JobFavorite) error {
	return dao.db.Create(favorite).Error
}

func (dao *JobFavoriteDAO) Delete(userID, jobID uint) error {
	return dao.db.Where("user_id = ? AND job_id = ?", userID, jobID).Delete(&model.JobFavorite{}).Error
}

func (dao *JobFavoriteDAO) IsFavorited(userID, jobID uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.JobFavorite{}).
		Where("user_id = ? AND job_id = ?", userID, jobID).
		Count(&count).Error
	return count > 0, err
}

func (dao *JobFavoriteDAO) ListByUser(userID uint, page, size int) ([]model.JobFavorite, int64, error) {
	var favorites []model.JobFavorite
	var total int64

	offset := (page - 1) * size
	err := dao.db.Model(&model.JobFavorite{}).
		Where("user_id = ?", userID).
		Count(&total).
		Offset(offset).
		Limit(size).
		Order("create_time DESC").
		Find(&favorites).Error

	return favorites, total, err
}

// GetUserFavoriteJobs 获取用户收藏的职位详情
func (dao *JobFavoriteDAO) GetUserFavoriteJobs(userID uint) ([]response.FavoriteJobDetail, error) {
	var jobs []response.FavoriteJobDetail

	err := dao.db.Model(&model.JobFavorite{}).
		Select("j.id as job_id, j.job_salary, j.job_salary_max as salary_max, j.update_time").
		Joins("JOIN t_rc_job j ON j.id = job_id").
		Where("user_id = ?", userID).
		Scan(&jobs).Error

	return jobs, err
}
