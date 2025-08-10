package service

import (
	"encoding/json"
	"fmt"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
)

type NotificationService struct {
	notificationDAO *dao.NotificationDAO
	templateDAO     *dao.NotificationTemplateDAO
}

func NewNotificationService(notificationDAO *dao.NotificationDAO, templateDAO *dao.NotificationTemplateDAO) *NotificationService {
	return &NotificationService{
		notificationDAO: notificationDAO,
		templateDAO:     templateDAO,
	}
}

// Create 创建通知
func (s *NotificationService) Create(notification *model.Notification) error {
	return s.notificationDAO.Create(notification)
}

// MarkAsRead 标记通知为已读
func (s *NotificationService) MarkAsRead(id uint) error {
	return s.notificationDAO.UpdateReadStatus(id, true)
}

// ListUserNotifications 获取用户的通知列表
func (s *NotificationService) ListUserNotifications(userID uint, page, size int) ([]model.Notification, int64, error) {
	return s.notificationDAO.ListByUser(userID, page, size)
}

// GetUnreadCount 获取用户未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationDAO.CountUnread(userID)
}

// SendNotification 发送通知
func (s *NotificationService) SendNotification(userID uint, userType model.UserType, templateCode string, variables map[string]interface{}) error {
	// 获取模板
	tmpl, err := s.templateDAO.GetByCode(templateCode)
	if err != nil {
		return err
	}

	// 检查用户类型是否支持
	if !containsUserType(tmpl.UserTypes, userType) {
		return fmt.Errorf("template not supported for user type: %v", userType)
	}

	// 渲染模板
	varsJSON, _ := json.Marshal(variables)
	notification := &model.Notification{
		UserID:     userID,
		UserType:   userType,
		Type:       tmpl.Type,
		TemplateID: templateCode,
		Variables:  string(varsJSON),
		Channels:   tmpl.Channels,
	}

	// 渲染内容
	title, content, err := s.renderTemplate(tmpl, variables)
	if err != nil {
		return err
	}
	notification.Title = title
	notification.Content = content

	// 发送通知
	return s.sendToChannels(notification)
}

// sendToChannels 发送到各个通知渠道
func (s *NotificationService) sendToChannels(notification *model.Notification) error {
	if notification.IsChannelEnabled(model.ChannelInApp) {
		s.notificationDAO.Create(notification)
	}
	if notification.IsChannelEnabled(model.ChannelEmail) {
		// 发送邮件
	}
	if notification.IsChannelEnabled(model.ChannelSMS) {
		// 发送短信
	}
	// ... 其他渠道
	return nil
}

// renderTemplate 渲染模板
func (s *NotificationService) renderTemplate(tmpl *model.NotificationTemplate, vars map[string]interface{}) (string, string, error) {
	// titleTmpl, err := template.New("title").Parse(tmpl.Title)
	// if err != nil {
	// 	return "", "", err
	// }
	// contentTmpl, err := template.New("content").Parse(tmpl.Content)
	// if err != nil {
	// 	return "", "", err
	// }

	// ... 渲染逻辑
	return "rendered title", "rendered content", nil
}

// containsUserType 检查用户类型是否在列表中
func containsUserType(userTypes []model.UserType, userType model.UserType) bool {
	for _, ut := range userTypes {
		if ut == userType {
			return true
		}
	}
	return false
}
