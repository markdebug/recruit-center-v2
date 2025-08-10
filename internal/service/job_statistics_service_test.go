package service

import (
	"testing"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestJobStatisticsService_IncrementViewCount(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobStatisticsDAO(db)
	service := NewJobStatisticsService(mockDAO)
	//生成测试数据
	err := service.IncrementViewCount(1, 1)
	if err != nil {
		t.Fatalf("IncrementViewCount failed: %v", err)
	}
}
