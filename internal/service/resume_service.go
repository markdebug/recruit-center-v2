package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/request"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/logger"
	"org.thinkinai.com/recruit-center/pkg/oss"
	"org.thinkinai.com/recruit-center/pkg/utils"
)

type ResumeService struct {
	resumeDao *dao.ResumeDAO
}

func NewResumeService(resumeDao *dao.ResumeDAO) *ResumeService {
	return &ResumeService{resumeDao: resumeDao}
}

// convertToResumeResponse 将 model.Resume 转换为 response.ResumeResponse
func (s *ResumeService) convertToResumeResponse(resume *model.Resume) *response.ResumeResponse {
	if resume == nil {
		return nil
	}

	resp := &response.ResumeResponse{
		ID:            resume.ID,
		UserID:        resume.UserID,
		Name:          resume.Name,
		Avatar:        resume.Avatar,
		Gender:        resume.Gender,
		Birthday:      resume.Birthday,
		Phone:         resume.Phone,
		Email:         resume.Email,
		Location:      resume.Location,
		Experience:    resume.Experience,
		JobStatus:     resume.JobStatus,
		ExpectedJob:   resume.ExpectedJob,
		ExpectedCity:  resume.ExpectedCity,
		Introduction:  resume.Introduction,
		Skills:        resume.Skills,
		ShareToken:    resume.ShareToken,
		AccessStatus:  resume.AccessStatus,
		WorkingStatus: resume.WorkingStatus,
		Status:        resume.Status,
	}

	// 转换教育经历
	resp.Educations = make([]response.EducationResponse, len(resume.Educations))
	for i, edu := range resume.Educations {
		resp.Educations[i] = response.EducationResponse{
			ID:        edu.ID,
			ResumeID:  edu.ResumeID,
			School:    edu.School,
			Major:     edu.Major,
			Degree:    edu.Degree,
			StartTime: edu.StartTime,
			EndTime:   edu.EndTime,
		}
	}

	// 转换工作经历
	resp.WorkExperiences = make([]response.WorkExperienceResponse, len(resume.WorkExperiences))
	for i, work := range resume.WorkExperiences {
		resp.WorkExperiences[i] = response.WorkExperienceResponse{
			ID:          work.ID,
			ResumeID:    work.ResumeID,
			CompanyName: work.CompanyName,
			Position:    work.Position,
			Department:  work.Department,
			StartTime:   work.StartTime,
			EndTime:     work.EndTime,
			Description: work.Description,
			Achievement: work.Achievement,
		}
	}

	// 转换项目经历
	resp.Projects = make([]response.ProjectResponse, len(resume.Projects))
	for i, proj := range resume.Projects {
		resp.Projects[i] = response.ProjectResponse{
			ID:          proj.ID,
			ResumeID:    proj.ResumeID,
			Name:        proj.Name,
			Role:        proj.Role,
			StartTime:   proj.StartTime,
			EndTime:     proj.EndTime,
			Description: proj.Description,
			Technology:  proj.Technology,
			Achievement: proj.Achievement,
		}
	}

	// 转换附件
	resp.Attachments = make([]response.AttachmentResponse, len(resume.Attachments))
	for i, att := range resume.Attachments {
		resp.Attachments[i] = response.AttachmentResponse{
			ID:       att.ID,
			ResumeID: att.ResumeID,
			FileName: att.FileName,
			FileURL:  att.FileURL,
			FileSize: att.FileSize,
			FileType: att.FileType,
			Status:   att.Status,
		}
	}

	return resp
}

// Create 创建简历
func (s *ResumeService) Create(userID uint, req *request.CreateResumeRequest) (*model.Resume, error) {
	// 检查用户是否已有简历
	if _, err := s.resumeDao.GetByUser(userID); err == nil {
		logger.L.Warn("用户已存在简历", zap.Uint("userID", userID))
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

// 根据分享token查询用户简历
func (s *ResumeService) GetByShareToken(token string) (*model.Resume, error) {
	if token == "" {
		return nil, errors.New(errors.InvalidParams)
	}
	resume, err := s.resumeDao.GetByShareToken(token)
	if err != nil {
		logger.L.Warn("获取简历失败", zap.String("token", token), zap.Error(err))
		return nil, fmt.Errorf("获取简历失败: %w", err)
	}
	if resume == nil {
		logger.L.Warn("简历不存在", zap.String("token", token))
		return nil, errors.New(errors.ResumeNotFound)
	}
	// 检查简历隐私设置
	if resume.AccessStatus == int(enums.Hide) {
		logger.L.Warn("简历访问被拒绝", zap.Uint("resumeID", resume.ID), zap.Uint("userID", resume.UserID))
		return nil, errors.New(errors.ResumeAccessDenied)
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

// GetResumeByID 获取简历详情
func (s *ResumeService) GetResumeByID(id uint) (*response.ResumeResponse, error) {
	resume, err := s.resumeDao.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResumeResponse(resume), nil
}

// ListResumes 获取简历列表
func (s *ResumeService) ListResumes(page, pageSize int) (*response.ResumeListResponse, error) {
	resumes, total, err := s.resumeDao.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := &response.ResumeListResponse{
		Total:   total,
		Records: make([]response.ResumeResponse, len(resumes)),
	}

	for i, resume := range resumes {
		resp.Records[i] = *s.convertToResumeResponse(&resume)
	}

	return resp, nil
}
