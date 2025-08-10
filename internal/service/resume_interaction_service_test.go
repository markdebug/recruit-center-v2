package service

import (
	"testing"

	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestRecordView(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	err := service.RecordView(1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddFavorite(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	err := service.AddFavorite(1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRemoveFavorite(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	err := service.RemoveFavorite(1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetInteractionStats(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	stats, err := service.GetInteractionStats(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stats["view_count"] != 10 {
		t.Errorf("expected view_count 10, got %d", stats["view_count"])
	}
	if stats["favorite_count"] != 5 {
		t.Errorf("expected favorite_count 5, got %d", stats["favorite_count"])
	}
}

func TestGetInteractionStats_Error(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	_, err := service.GetInteractionStats(1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIsFavorited(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	favorited, err := service.IsFavorited(1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !favorited {
		t.Errorf("expected favorited true, got false")
	}
}

func TestIsFavorited_Error(t *testing.T) {
	db := testutil.SetupTestDB(t)
	interactionDAO := dao.NewResumeInteractionDAO(db)
	service := NewResumeInteractionService(interactionDAO)
	_, err := service.IsFavorited(1, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
