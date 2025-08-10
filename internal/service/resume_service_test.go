package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/testutil"
)

func TestResumeService_Create_Success(t *testing.T) {
	// Arrange
	db := testutil.SetupTestDB(t)
	mockDao := dao.NewResumeDAO(db)
	service := &ResumeService{resumeDao: mockDao}
	birthday, err := time.Parse("2006-01-02", "1990-01-01")
	assert.NoError(t, err)
	eduStart, err := time.Parse("2006", "2008")
	assert.NoError(t, err)
	eduEnd, err := time.Parse("2006", "2012")
	assert.NoError(t, err)
	workStart, err := time.Parse("2006", "2012")
	assert.NoError(t, err)
	workEnd, err := time.Parse("2006", "2017")
	assert.NoError(t, err)
	projStart, err := time.Parse("2006", "2016")
	assert.NoError(t, err)
	projEnd, err := time.Parse("2006", "2017")
	assert.NoError(t, err)

	req := &request.CreateResumeRequest{
		Name:         "John Doe",
		Avatar:       "avatar.png",
		Gender:       1,
		Birthday:     birthday,
		Phone:        "1234567890",
		Email:        "john@example.com",
		Location:     "Beijing",
		Experience:   5,
		JobStatus:    1,
		ExpectedJob:  "Golang Developer",
		ExpectedCity: "Shanghai",
		Introduction: "Hello",
		Skills:       "Go, Docker",
		Educations: []request.EducationRequest{
			{School: "Tsinghua", Major: "CS", Degree: "Bachelor", StartTime: eduStart, EndTime: eduEnd},
		},
		WorkExperiences: []request.WorkExperienceRequest{
			{CompanyName: "ABC", Position: "Dev", Department: "IT", StartTime: workStart, EndTime: workEnd, Description: "Worked", Achievement: "Award"},
		},
		Projects: []request.ProjectRequest{
			{Name: "Proj1", Role: "Lead", StartTime: projStart, EndTime: projEnd, Description: "Desc", Technology: "Go", Achievement: "Success"},
		},
	}
	userID := uint(10012)
	// Act
	createResume, err := service.Create(userID, req)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createResume)
	assert.Equal(t, "John Doe", createResume.Name)
	assert.Equal(t, "mocktoken", createResume.ShareToken)
}

func TestResumeService_Create_AlreadyExists(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDao := dao.NewResumeDAO(db)
	service := &ResumeService{resumeDao: mockDao}
	userID := uint(10012)
	req := &request.CreateResumeRequest{}

	resume, err := service.Create(userID, req)

	assert.Nil(t, resume)
	assert.Error(t, err)
}

func TestResumeService_Create_GenerateNanoIDError(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDao := dao.NewResumeDAO(db)
	service := &ResumeService{resumeDao: mockDao}
	userID := uint(10012)
	req := &request.CreateResumeRequest{}
	resume, err := service.Create(userID, req)
	assert.Nil(t, resume)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "生成分享令牌失败")
}

func TestResumeService_Create_DaoCreateError(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDao := dao.NewResumeDAO(db)
	service := &ResumeService{resumeDao: mockDao}
	userID := uint(10012)
	req := &request.CreateResumeRequest{}
	resume, err := service.Create(userID, req)
	assert.Nil(t, resume)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestResumeService_GetByShareToken(t *testing.T) {
	db := testutil.SetupTestDB(t)
	mockDao := dao.NewResumeDAO(db)
	service := &ResumeService{resumeDao: mockDao}
	token := "aM6tR3KU2W"
	// 获取分享令牌
	resume, err := service.GetByShareToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, resume)

}
