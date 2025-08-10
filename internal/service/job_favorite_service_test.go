package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

// TestJobFavoriteService_AddFavorite 测试添加收藏
func TestJobFavoriteService_AddFavorite(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobFavoriteDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao, mockDAO, dao.NewJobApplyDAO(db))
	service := NewJobFavoriteService(mockDAO, mockJobService)

	err := service.AddFavorite(1, 3)
	if err != nil {
		t.Errorf("AddFavorite() error = %v", err)
		return
	}
}

// TestJobFavoriteService_RemoveFavorite 测试取消收藏
func TestJobFavoriteService_RemoveFavorite(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobFavoriteDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao, mockDAO, dao.NewJobApplyDAO(db))
	service := NewJobFavoriteService(mockDAO, mockJobService)

	err := service.RemoveFavorite(1, 3)
	if err != nil {
		t.Errorf("RemoveFavorite() error = %v", err)
		return
	}
}

// TestJobFavoriteService_ListFavorites 测试列出收藏
func TestJobFavoriteService_ListFavorites(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobFavoriteDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao, mockDAO, dao.NewJobApplyDAO(db))
	service := NewJobFavoriteService(mockDAO, mockJobService)

	favorites, err := service.ListFavorites(1, 1, 10)
	if err != nil {
		t.Errorf("ListFavorites() error = %v", err)
		return
	}
	assert.NotEmpty(t, favorites)
}

// TestJobFavoriteService_GetUserStatistics 测试获取统计信息
func TestJobFavoriteService_GetUserStatistics(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobFavoriteDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao, mockDAO, dao.NewJobApplyDAO(db))
	service := NewJobFavoriteService(mockDAO, mockJobService)
	stats, err := service.GetUserStatistics(1)
	if err != nil {
		t.Errorf("GetUserStatistics() error = %v", err)
		return

	}
	assert.NotNil(t, stats)
}
