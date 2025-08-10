package dao

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestJobApplyDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	dao := NewJobApplyDAO(db)

	tests := []struct {
		name    string
		apply   *model.JobApply
		wantErr bool
	}{
		{
			name: "valid apply",
			apply: &model.JobApply{
				JobID:     1,
				UserID:    1,
				ResumeID:  1,
				Status:    1,
				ApplyTime: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dao.Create(tt.apply)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.apply.ID)
			}
		})
	}
}

func TestJobApplyDAO_GetByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	dao := NewJobApplyDAO(db)

	// 创建测试数据
	apply := &model.JobApply{
		JobID:     1,
		UserID:    1,
		ResumeID:  1,
		Status:    1,
		ApplyTime: time.Now(),
	}
	err := dao.Create(apply)
	assert.NoError(t, err)

	// 测试获取
	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "existing apply",
			id:      apply.ID,
			wantErr: false,
		},
		{
			name:    "non-existing apply",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.GetByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, apply.ID, got.ID)
			}
		})
	}
}
