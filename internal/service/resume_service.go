package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/oss"
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
		return nil, errors.New("用户已有简历")
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

	if err := s.resumeDao.Create(resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// GetByUser 获取用户简历
func (s *ResumeService) GetByUser(userID uint) (*model.Resume, error) {
	return s.resumeDao.GetByUser(userID)
}

// Update 更新简历
func (s *ResumeService) Update(resume *model.Resume) error {
	return s.resumeDao.Update(resume)
}

// 更新简历访问状态
func (s *ResumeService) UpdateAccessStatus(userID uint, status int) error {
	if _, err := enums.ParseResumeAccess(2); err != nil {
		return errors.New("无效的简历访问状态")
	}

	resume, err := s.resumeDao.GetByUser(userID)
	if err != nil {
		return fmt.Errorf("获取简历失败: %w", err)
	}

	resume.AccessStatus = status
	return s.resumeDao.Update(resume)
}

// 更新简历工作状态
func (s *ResumeService) UpdateWorkingStatus(userID uint, targetStatus int) error {
	if _, err := enums.ParseWorkingStatus(1); err != nil {
		return errors.New("无效的在职状态")
	}

	resume, err := s.resumeDao.GetByUser(userID)
	if err != nil {
		return fmt.Errorf("获取简历失败: %w", err)
	}

	resume.WorkingStatus = targetStatus
	return s.resumeDao.Update(resume)
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
