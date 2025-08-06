package service

import (
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
)

type ResumeInteractionService struct {
	interactionDAO *dao.ResumeInteractionDAO
}

func NewResumeInteractionService(interactionDAO *dao.ResumeInteractionDAO) *ResumeInteractionService {
	return &ResumeInteractionService{interactionDAO: interactionDAO}
}

// RecordView 记录查看
func (s *ResumeInteractionService) RecordView(resumeID, userID uint) error {
	return s.interactionDAO.AddInteraction(resumeID, userID, model.InteractionView)
}

// AddFavorite 添加收藏
func (s *ResumeInteractionService) AddFavorite(resumeID, userID uint) error {
	return s.interactionDAO.AddInteraction(resumeID, userID, model.InteractionFavorite)
}

// RemoveFavorite 取消收藏
func (s *ResumeInteractionService) RemoveFavorite(resumeID, userID uint) error {
	return s.interactionDAO.RemoveInteraction(resumeID, userID, model.InteractionFavorite)
}

// GetAllInteractions 查询简历的所有信息，并进行分组统计返回
func (s *ResumeInteractionService) GetAllInteractions(resumeID uint) (map[string]int, error) {
	interactions, err := s.interactionDAO.GetAllByResumeID(resumeID)
	if err != nil {
		return nil, err
	}
	mapStats := make(map[string]int)
	//使用lambda函数进行分组统计
	for _, interaction := range interactions {
		mapStats[string(interaction.Type)]++
	}
	return mapStats, nil
}

// GetInteractionStats 获取交互统计
func (s *ResumeInteractionService) GetInteractionStats(resumeID uint) (map[string]int64, error) {
	viewCount, err := s.interactionDAO.GetStats(resumeID, model.InteractionView)
	if err != nil {
		return nil, err
	}

	favoriteCount, err := s.interactionDAO.GetStats(resumeID, model.InteractionFavorite)
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"view_count":     viewCount,
		"favorite_count": favoriteCount,
	}, nil
}

// IsFavorited 检查是否已收藏
func (s *ResumeInteractionService) IsFavorited(resumeID, userID uint) (bool, error) {
	return s.interactionDAO.HasInteraction(resumeID, userID, model.InteractionFavorite)
}
