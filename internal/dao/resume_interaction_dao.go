package dao

import (
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

type ResumeInteractionDAO struct {
	db *gorm.DB
}

func NewResumeInteractionDAO(db *gorm.DB) *ResumeInteractionDAO {
	return &ResumeInteractionDAO{db: db}
}

// AddInteraction 添加交互记录
func (d *ResumeInteractionDAO) AddInteraction(resumeID, userID uint, interType model.InteractionType) error {
	interaction := &model.ResumeInteraction{
		ResumeID: resumeID,
		UserID:   userID,
		Type:     interType,
	}

	// 使用upsert确保记录唯一性
	// return d.db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "resume_id"}, {Name: "user_id"}, {Name: "type"}},
	// 	DoUpdates: clause.AssignmentColumns([]string{"last_time"}),
	// }).Create(interaction).Error
	return d.db.Save(interaction).Error
}

// RemoveInteraction 移除交互记录
func (d *ResumeInteractionDAO) RemoveInteraction(resumeID, userID uint, interType model.InteractionType) error {

	return d.db.Where("resume_id = ? AND user_id = ? AND type = ?",
		resumeID, userID, interType).Delete(&model.ResumeInteraction{}).Error
}

// GetStats 获取统计数据
func (d *ResumeInteractionDAO) GetStats(resumeID uint, interType model.InteractionType) (int64, error) {
	var count int64
	err := d.db.Model(&model.ResumeInteraction{}).
		Where("resume_id = ? AND type = ?", resumeID, interType).
		Count(&count).Error
	return count, err
}

// GetAllByResumeID 获取指定简历的所有信息
func (d *ResumeInteractionDAO) GetAllByResumeID(resumeID uint) ([]model.ResumeInteraction, error) {
	var interactions []model.ResumeInteraction
	err := d.db.Where("resume_id = ?", resumeID).Find(&interactions).Error
	return interactions, err
}

// HasInteraction 检查是否有交互记录
func (d *ResumeInteractionDAO) HasInteraction(resumeID, userID uint, interType model.InteractionType) (bool, error) {
	var count int64
	err := d.db.Model(&model.ResumeInteraction{}).
		Where("resume_id = ? AND user_id = ? AND type = ?",
			resumeID, userID, interType).
		Count(&count).Error
	return count > 0, err
}

// GetRecentViewers 获取最近查看者
func (d *ResumeInteractionDAO) GetRecentViewers(resumeID uint, limit int) ([]uint, error) {
	var userIDs []uint
	err := d.db.Model(&model.ResumeInteraction{}).
		Where("resume_id = ? AND type = ?", resumeID, model.InteractionView).
		Order("last_time DESC").
		Limit(limit).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
