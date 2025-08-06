package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=192.168.10.153 user=pgtest password=123456 dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(
		&model.Resume{},
		&model.Education{},
		&model.WorkExperience{},
		&model.Project{},
		&model.ResumeAttachment{},
		&model.ResumeInteraction{},
	)
	assert.NoError(t, err)
	return db
}
