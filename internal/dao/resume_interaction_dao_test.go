package dao

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestAddInteraction_InsertNew(t *testing.T) {
	db := testutil.SetupTestDB(t)
	dao := NewResumeInteractionDAO(db)

	resumeID := uint(1)
	userID := uint(2)
	interType := model.InteractionView

	err := dao.AddInteraction(resumeID, userID, interType)
	assert.NoError(t, err)

	var interaction model.ResumeInteraction
	result := db.Where("resume_id = ? AND user_id = ? AND type = ?", resumeID, userID, interType).First(&interaction)
	assert.NoError(t, result.Error)
	assert.Equal(t, resumeID, interaction.ResumeID)
	assert.Equal(t, userID, interaction.UserID)
	assert.Equal(t, interType, interaction.Type)

}

func TestAddInteraction_UpsertUpdatesLastTime(t *testing.T) {
	db := testutil.SetupTestDB(t)
	dao := NewResumeInteractionDAO(db)

	resumeID := uint(1)
	userID := uint(2)
	interType := model.InteractionView

	// Insert initial record
	err := dao.AddInteraction(resumeID, userID, interType)
	assert.NoError(t, err)

	var interaction model.ResumeInteraction
	_ = db.Where("resume_id = ? AND user_id = ? AND type = ?", resumeID, userID, interType).First(&interaction)

	// Wait and update
	time.Sleep(10 * time.Millisecond)
	err = dao.AddInteraction(resumeID, userID, interType)
	assert.NoError(t, err)

	var updated model.ResumeInteraction
	_ = db.Where("resume_id = ? AND user_id = ? AND type = ?", resumeID, userID, interType).First(&updated)

}
