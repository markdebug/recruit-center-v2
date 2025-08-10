package service

import (
	"testing"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestNotificationService_SendNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDAO := dao.NewNotificationDAO(db)
	mockTemplateDAO := dao.NewNotificationTemplateDAO(db)
	service := NewNotificationService(mockDAO, mockTemplateDAO)
	err := service.SendNotification(1, model.UserTypeAdmin, "Test Content", nil)
	if err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}

}
