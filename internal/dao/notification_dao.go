package dao

import (
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

type NotificationDAO struct {
	db *gorm.DB
}

type NotificationTemplateDAO struct {
	db *gorm.DB
}

func NewNotificationDAO(db *gorm.DB) *NotificationDAO {
	return &NotificationDAO{db: db}
}

func NewNotificationTemplateDAO(db *gorm.DB) *NotificationTemplateDAO {
	return &NotificationTemplateDAO{db: db}
}

func (dao *NotificationDAO) Create(notification *model.Notification) error {
	return dao.db.Create(notification).Error
}

func (dao *NotificationDAO) UpdateReadStatus(id uint, isRead bool) error {
	return dao.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", isRead).Error
}

func (dao *NotificationDAO) ListByUser(userID uint, page, size int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	offset := (page - 1) * size
	db := dao.db.Model(&model.Notification{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(size).Order("create_time desc").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (dao *NotificationDAO) CountUnread(userID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (dao *NotificationTemplateDAO) GetByCode(code string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	err := dao.db.Where("code = ? AND is_active = ?", code, true).First(&template).Error
	return &template, err
}
