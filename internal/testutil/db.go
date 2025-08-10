package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

// SetupTestDB 创建测试数据库连接
func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=192.168.2.185 user=pgtest password=123456 dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(
		&model.Resume{},
		&model.Education{},
		&model.WorkExperience{},
		&model.Project{},
		&model.ResumeAttachment{},
		&model.ResumeInteraction{},
		&model.Job{},
		&model.JobApply{},
		&model.JobStatistics{},
		&model.Notification{},
		&model.NotificationTemplate{},
	)
	assert.NoError(t, err)
	return db
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("failed to get underlying sql.DB: %v", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}
