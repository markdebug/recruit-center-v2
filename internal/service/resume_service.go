package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/oss"
	"org.thinkinai.com/recruit-center/pkg/utils"
)

type ResumeService struct {
	resumeDao *dao.ResumeDAO
}

func NewResumeService(resumeDao *dao.ResumeDAO) *ResumeService {
	return &ResumeService{resumeDao: resumeDao}
}

// Create 创建简历
func (s *ResumeService) Create(userID uint, req *request.CreateResumeRequest) (*model.Resume, error) {
	// 检查用户是否已有简历
	if _, err := s.resumeDao.GetByUser(userID); err == nil {
		return nil, errors.New(errors.ResumeExists)
	}

	resume := &model.Resume{
		UserID:       userID,
		Name:         req.Name,
		Avatar:       req.Avatar,
		Gender:       req.Gender,
		Birthday:     req.Birthday,
		Phone:        req.Phone,
		Email:        req.Email,
		Location:     req.Location,
		Experience:   req.Experience,
		JobStatus:    req.JobStatus,
		ExpectedJob:  req.ExpectedJob,
		ExpectedCity: req.ExpectedCity,
		Introduction: req.Introduction,
		Skills:       req.Skills,
	}

	// 添加教育经历
	for _, edu := range req.Educations {
		resume.Educations = append(resume.Educations, model.Education{
			School:    edu.School,
			Major:     edu.Major,
			Degree:    edu.Degree,
			StartTime: edu.StartTime,
			EndTime:   edu.EndTime,
		})
	}

	// 添加工作经历
	for _, work := range req.WorkExperiences {
		resume.WorkExperiences = append(resume.WorkExperiences, model.WorkExperience{
			CompanyName: work.CompanyName,
			Position:    work.Position,
			Department:  work.Department,
			StartTime:   work.StartTime,
			EndTime:     work.EndTime,
			Description: work.Description,
			Achievement: work.Achievement,
		})
	}

	// 添加项目经历
	for _, proj := range req.Projects {
		resume.Projects = append(resume.Projects, model.Project{
			Name:        proj.Name,
			Role:        proj.Role,
			StartTime:   proj.StartTime,
			EndTime:     proj.EndTime,
			Description: proj.Description,
			Technology:  proj.Technology,
			Achievement: proj.Achievement,
		})
	}
	// 生成分享令牌
	shareToken, err := utils.GenerateNanoID(10) // 生成10位的友好URL
	if err != nil {
		return nil, fmt.Errorf("生成分享令牌失败: %w", err)
	}
	resume.ShareToken = shareToken
	if err := s.resumeDao.Create(resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// UpdateBasic 更新基本信息
func (s *ResumeService) UpdateBasic(resumeID uint, req *request.UpdateResumeBasicRequest) error {
	resume := &model.Resume{
		ID:           resumeID,
		Name:         req.Name,
		Avatar:       req.Avatar,
		Gender:       req.Gender,
		Birthday:     req.Birthday,
		Phone:        req.Phone,
		Email:        req.Email,
		Location:     req.Location,
		Introduction: req.Introduction,
	}
	return s.resumeDao.UpdateBasic(resume)
}

// UpdateEducation 更新教育经历
func (s *ResumeService) UpdateEducation(resumeID uint, req *request.UpdateResumeEducationRequest) error {
	edu := &model.Education{
		ResumeID:  resumeID,
		School:    req.School,
		Major:     req.Major,
		Degree:    req.Degree,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if req.ID > 0 {
		edu.ID = req.ID
		return s.resumeDao.UpdateEducation(edu)
	}
	return s.resumeDao.AddEducation(edu)
}

// UpdateWork 更新工作经历
func (s *ResumeService) UpdateWork(resumeID uint, req *request.UpdateResumeWorkRequest) error {
	work := &model.WorkExperience{
		ResumeID:    resumeID,
		CompanyName: req.CompanyName,
		Position:    req.Position,
		Department:  req.Department,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: req.Description,
		Achievement: req.Achievement,
	}

	if req.ID > 0 {
		work.ID = req.ID
		return s.resumeDao.UpdateWorkExperience(work)
	}
	return s.resumeDao.AddWorkExperience(work)
}

// GetByUser 获取用户简历
func (s *ResumeService) GetByUser(userID uint) (*model.Resume, error) {
	return s.resumeDao.GetByUser(userID)
}

// 根据用户id和简历id获取简历
func (s *ResumeService) GetByUserIDAndResumeID(userID uint, resumeID uint) (*model.Resume, error) {
	resume, err := s.resumeDao.GetByID(resumeID)
	if err != nil {
		return nil, fmt.Errorf("获取简历失败: %w", err)
	}
	if resume.UserID != userID {
		return nil, errors.New(errors.ResumeAccessDenied)
	}
	return resume, nil
}

// 更新简历访问状态
func (s *ResumeService) UpdateAccessStatus(userID uint, status int) error {
	if _, err := enums.ParseResumeAccess(2); err != nil {
		return errors.New(errors.ResumeUpdateStatus)
	}

	resume, err := s.resumeDao.GetByUser(userID)
	if err != nil {
		return fmt.Errorf("获取简历失败: %w", err)
	}

	resume.AccessStatus = status
	return s.resumeDao.UpdateBasic(resume)
}

// 更新简历工作状态
func (s *ResumeService) UpdateWorkingStatus(userID uint, targetStatus int) error {
	if _, err := enums.ParseWorkingStatus(1); err != nil {
		return errors.New(errors.ResumeUpdateStatus)
	}

	resume, err := s.resumeDao.GetByUser(userID)
	if err != nil {
		return fmt.Errorf("获取简历失败: %w", err)
	}

	resume.WorkingStatus = targetStatus
	return s.resumeDao.UpdateBasic(resume)
}

// UploadResumeFile 上传简历文件到MinIO
func (s *ResumeService) UploadResumeFile(userID uint, file io.Reader, filename string) (string, error) {

	// 生成唯一的文件名
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("resumes/%d/%s%s", userID, time.Now().Format("20060102150405"), ext)

	// 上传文件到MinIO
	_, err := oss.MinioClient.PutObject(
		context.Background(),
		"resumes",  // bucket名称
		objectName, // 对象名称
		file,       // 文件读取器
		-1,         // 文件大小 (-1表示未知)
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	// 生成文件访问URL
	fileURL := fmt.Sprintf("%s/%s/%s", oss.MinioEndpoint, "resumes", objectName)

	// 更新简历记录
	err = s.resumeDao.UpdateAttachmentURL(userID, fileURL, objectName)
	if err != nil {
		return "", fmt.Errorf("更新简历记录失败: %w", err)
	}

	return fileURL, nil
}
