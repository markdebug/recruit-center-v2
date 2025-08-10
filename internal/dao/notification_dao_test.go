package dao

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestNotificationDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	dao := NewNotificationDAO(db)

	tests := []struct {
		name         string
		notification *model.Notification
		wantErr      bool
	}{
		{
			name: "valid notification",
			notification: &model.Notification{
				UserID:     1,
				Title:      "Test Title",
				Content:    "Test Content",
				Type:       model.NotificationTypeJobApply,
				IsRead:     false,
				Channels:   model.ChannelInApp,
				CreateTime: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dao.Create(tt.notification)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.notification.ID)
			}
		})
	}
}

func TestNotificationDAO_ListByUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	dao := NewNotificationDAO(db)

	// 创建测试数据
	notifications := []model.Notification{
		{UserID: 1, Title: "Test 1"},
		{UserID: 1, Title: "Test 2"},
		{UserID: 2, Title: "Test 3"},
	}

	for _, n := range notifications {
		err := dao.Create(&n)
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		userID    uint
		page      int
		size      int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "user 1 notifications",
			userID:    1,
			page:      1,
			size:      10,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "user 2 notifications",
			userID:    2,
			page:      1,
			size:      10,
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := dao.ListByUser(tt.userID, tt.page, tt.size)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(got))
				assert.Equal(t, int64(tt.wantCount), total)
			}
		})
	}
}
