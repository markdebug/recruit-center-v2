package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/internal/model"
)

func TestResumeDAO_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{
		UserID: 1,
		Name:   "Test Resume",
	}
	err := dao.Create(resume)
	assert.NoError(t, err)
	assert.NotZero(t, resume.ID)

	got, err := dao.GetByID(resume.ID)
	assert.NoError(t, err)
	assert.Equal(t, resume.Name, got.Name)
}

func TestResumeDAO_Update(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 2, Name: "Old"}
	_ = dao.Create(resume)

	resume.Name = "New"
	err := dao.Update(resume)
	assert.NoError(t, err)

	got, _ := dao.GetByID(resume.ID)
	assert.Equal(t, "New", got.Name)
}

func TestResumeDAO_GetByUser(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 3, Name: "UserResume"}
	_ = dao.Create(resume)

	got, err := dao.GetByUser(3)
	assert.NoError(t, err)
	assert.Equal(t, "UserResume", got.Name)
}

func TestResumeDAO_AddAndDeleteEducation(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 4, Name: "EduResume"}
	_ = dao.Create(resume)

	edu := &model.Education{ResumeID: resume.ID, School: "Test School"}
	err := dao.AddEducation(edu)
	assert.NoError(t, err)
	assert.NotZero(t, edu.ID)

	err = dao.DeleteEducation(edu.ID)
	assert.NoError(t, err)
}

func TestResumeDAO_AddAndDeleteWorkExperience(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 5, Name: "WorkResume"}
	_ = dao.Create(resume)

	exp := &model.WorkExperience{ResumeID: resume.ID}
	err := dao.AddWorkExperience(exp)
	assert.NoError(t, err)
	assert.NotZero(t, exp.ID)

	err = dao.DeleteWorkExperience(exp.ID)
	assert.NoError(t, err)
}

func TestResumeDAO_AddAndDeleteProject(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 6, Name: "ProjResume"}
	_ = dao.Create(resume)

	proj := &model.Project{ResumeID: resume.ID, Name: "Test Project"}
	err := dao.AddProject(proj)
	assert.NoError(t, err)
	assert.NotZero(t, proj.ID)

	err = dao.DeleteProject(proj.ID)
	assert.NoError(t, err)
}

func TestResumeDAO_UpdateEducation(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 7, Name: "EduUpdate"}
	_ = dao.Create(resume)

	edu := &model.Education{ResumeID: resume.ID, School: "Old"}
	_ = dao.AddEducation(edu)

	edu.School = "New"
	err := dao.UpdateEducation(edu)
	assert.NoError(t, err)
}

func TestResumeDAO_UpdateWorkExperience(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 8, Name: "WorkUpdate"}
	_ = dao.Create(resume)

	exp := &model.WorkExperience{ResumeID: resume.ID}
	_ = dao.AddWorkExperience(exp)

	err := dao.UpdateWorkExperience(exp)
	assert.NoError(t, err)
}

func TestResumeDAO_UpdateProject(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 9, Name: "ProjUpdate"}
	_ = dao.Create(resume)

	proj := &model.Project{ResumeID: resume.ID, Name: "Old"}
	_ = dao.AddProject(proj)

	proj.Name = "New"
	err := dao.UpdateProject(proj)
	assert.NoError(t, err)
}

func TestResumeDAO_AddUpdateDeleteAttachment(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 10, Name: "AttachResume"}
	_ = dao.Create(resume)

	att := &model.ResumeAttachment{ResumeID: resume.ID, FileName: "file.txt", FileURL: "url"}
	err := dao.AddAttachment(att)
	assert.NoError(t, err)
	assert.NotZero(t, att.ID)

	att.FileName = "file2.txt"
	err = dao.UpdateAttachment(att)
	assert.NoError(t, err)

	err = dao.DeleteAttachment(att.ID)
	assert.NoError(t, err)
}

func TestResumeDAO_GetAttachments(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 11, Name: "AttachList"}
	_ = dao.Create(resume)

	att := &model.ResumeAttachment{ResumeID: resume.ID, FileName: "file.txt", FileURL: "url"}
	_ = dao.AddAttachment(att)

	atts, err := dao.GetAttachments(resume.ID)
	assert.NoError(t, err)
	assert.Len(t, atts, 1)
}

func TestResumeDAO_UpdateAttachmentURL(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	resume := &model.Resume{UserID: 12, Name: "AttachURL"}
	_ = dao.Create(resume)

	att := &model.ResumeAttachment{ResumeID: resume.ID, FileName: "file.txt", FileURL: "url"}
	_ = dao.AddAttachment(att)

	err := dao.UpdateAttachmentURL(att.ID, "newurl", "newfile.txt")
	assert.NoError(t, err)

	var got model.ResumeAttachment
	_ = db.First(&got, att.ID).Error
	assert.Equal(t, "newurl", got.FileURL)
	assert.Equal(t, "newfile.txt", got.FileName)
}

func TestResumeDAO_List(t *testing.T) {
	db := setupTestDB(t)
	dao := NewResumeDAO(db)

	for i := 0; i < 5; i++ {
		resume := &model.Resume{UserID: uint(100 + i), Name: "Resume"}
		_ = dao.Create(resume)
	}

	resumes, total, err := dao.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, resumes, 2)
}
