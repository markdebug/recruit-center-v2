package service

import (
	"testing"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/testutil"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// TestJobApplyService_Create 测试创建职位申请
func TestJobApplyService_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobApplyDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao)
	notifyDao := dao.NewNotificationDAO(db)
	notificationDAO := dao.NewNotificationTemplateDAO(db)
	mockNotificationService := NewNotificationService(notifyDao, notificationDAO)
	service := NewJobApplyService(mockDAO, mockJobService, mockNotificationService)
	apply := &model.JobApply{
		JobID:         1,
		UserID:        1,
		ResumeID:      1,
		Status:        int(enums.JobApplyPending),
		ApplyProgress: enums.JobApplyPending.String(),
	}
	err := service.Create(apply)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}
}

// TestJobApplyService_UpdateStatus 测试更新申请状态
func TestJobApplyService_UpdateStatus(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewJobApplyDAO(db)
	jobDao := dao.NewJobDAO(db)
	mockJobService := NewJobService(jobDao)
	notifyDao := dao.NewNotificationDAO(db)
	notificationDAO := dao.NewNotificationTemplateDAO(db)
	mockNotificationService := NewNotificationService(notifyDao, notificationDAO)

	service := NewJobApplyService(mockDAO, mockJobService, mockNotificationService)

	err := service.UpdateStatus(1, 1, enums.JobApplyAccepted)
	if err != nil {
		t.Errorf("UpdateStatus() error = %v", err)
		return
	}
}
