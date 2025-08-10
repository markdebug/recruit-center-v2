package service

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// JobService 职位服务
type JobService struct {
	jobDao      *dao.JobDAO
	favoriteDAO *dao.JobFavoriteDAO
	jobApplyDAO *dao.JobApplyDAO
}

// NewJobService 创建职位服务实例
func NewJobService(jobDao *dao.JobDAO, favoriteDAO *dao.JobFavoriteDAO, jobApplyDAO *dao.JobApplyDAO) *JobService {
	return &JobService{
		jobDao:      jobDao,
		favoriteDAO: favoriteDAO,
		jobApplyDAO: jobApplyDAO,
	}
}

// Create 创建职位
func (s *JobService) Create(job *model.Job) error {
	// 参数校验
	if job.Name == "" {
		return fmt.Errorf("职位名称不能为空")
	}
	if job.CompanyID == 0 {
		return fmt.Errorf("公司ID不能为空")
	}

	// 验证公司是否存在且有效
	// company, err := s.companyDAO.GetByID(job.CompanyID)
	// if err != nil {
	// 	return errors.Wrap(err, enums.CompanyNotFound)
	// }
	// if !company.IsActive() {
	// 	return errors.New(enums.CompanyInactive)
	// }

	// 验证职位类型
	if !job.ValidateJobType() {
		return errors.New(errors.BadRequest).WithMessage("无效的职位类型")
	}

	// 设置默认值
	if job.Status == 0 {
		job.Status = int(enums.JobStatusNormal)
	}
	if job.JobExpireTime.IsZero() {
		job.JobExpireTime = time.Now().AddDate(0, 1, 0) // 默认一个月后过期
	}

	// 创建职位
	err := s.jobDao.Create(job)
	if err != nil {
		logger.L.Error("创建职位失败",
			zap.Error(err),
			zap.String("job_name", job.Name),
			zap.Uint("company_id", job.CompanyID))
		return err
	}
	return nil
}

// VerifyCompanyOwner 验证职位是否属于指定公司
func (s *JobService) VerifyCompanyOwner(jobID, companyID uint) error {
	job, err := s.jobDao.GetByID(jobID)
	if err != nil {
		logger.L.Error("获取职位失败",
			zap.Error(err),
			zap.Uint("jobId", jobID))
		return errors.Wrap(err, errors.JobNotFound)
	}

	if job.CompanyID != companyID {
		logger.L.Warn("职位不属于该公司",
			zap.Uint("jobId", jobID),
			zap.Uint("companyId", companyID),
			zap.Uint("ownerCompanyId", job.CompanyID))
		return errors.New(errors.JobNotBelongToCompany)
	}
	return nil
}

// Update 更新职位信息
func (s *JobService) Update(job *model.Job) error {
	if err := s.VerifyCompanyOwner(job.ID, job.CompanyID); err != nil {
		return err
	}
	return s.jobDao.Update(job)
}

// Delete 删除职位
func (s *JobService) Delete(id uint, companyID uint) error {
	if err := s.VerifyCompanyOwner(id, companyID); err != nil {
		return err
	}
	return s.jobDao.Delete(id)
}

// UpdateStatus 更新职位状态
func (s *JobService) UpdateStatus(id uint, status int, companyID uint) error {
	logger.L.Info("更新职位状态",
		zap.Uint("job_id", id),
		zap.Int("status", status),
		zap.Uint("company_id", companyID))

	if err := s.VerifyCompanyOwner(id, companyID); err != nil {
		return err
	}

	// 检查状态是否有效
	if !enums.JobStatus(status).IsValid() {
		return errors.New(errors.InvalidJobStatus)
	}

	return s.jobDao.UpdateStatus(id, status)
}

// ConvertToJobResponse 将 model.Job 转换为 response.JobResponse
func (s *JobService) ConvertToJobResponse(job *model.Job, userID uint) *response.JobResponse {
	if job == nil {
		return nil
	}

	resp := &response.JobResponse{
		ID:            job.ID,
		Name:          job.Name,
		CompanyID:     job.CompanyID,
		JobSkill:      job.JobSkill,
		JobSalary:     job.JobSalary,
		JobDescribe:   job.JobDescribe,
		JobLocation:   job.JobLocation,
		JobExpireTime: job.JobExpireTime,
		Status:        job.Status,
		JobType:       job.JobType,
		JobCategory:   job.JobCategory,
		JobExperience: job.JobExperience,
		JobEducation:  job.JobEducation,
		JobBenefit:    job.JobBenefit,
		JobContact:    job.JobContact,
		JobSource:     job.JobSource,
		ViewCount:     job.ViewCount,
		ApplyCount:    job.ApplyCount,
		Priority:      job.Priority,
		Tags:          job.Tags,
	}

	// 如果提供了用户ID，查询用户状态
	if userID > 0 {
		// 查询收藏状态
		isFavorited, _ := s.favoriteDAO.IsFavorited(userID, job.ID)
		resp.IsFavorited = isFavorited

		// 查询投递状态
		hasApplied, _ := s.jobApplyDAO.GetByUserAndJob(userID, job.ID)
		resp.IsApplied = hasApplied != nil
	}

	return resp
}

// GetByID 获取职位详情
func (s *JobService) GetByID(id uint, userID uint) (*response.JobResponse, error) {
	// 从数据库获取
	job, err := s.jobDao.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.JobNotFound)
	}

	return s.ConvertToJobResponse(job, userID), nil
}

// List 获取职位列表
func (s *JobService) List(page, size int, userID uint) (*response.JobListResponse, error) {
	jobs, total, err := s.jobDao.List(page, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   total,
		Records: make([]response.JobResponse, len(jobs)),
	}

	for i, job := range jobs {
		resp.Records[i] = *s.ConvertToJobResponse(&job, userID)
	}

	return resp, nil
}

// SearchByKeyword 关键词搜索职位
func (s *JobService) SearchByKeyword(keyword string) (*response.JobListResponse, error) {
	jobs, err := s.jobDao.SearchByKeyword(keyword)
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   int64(len(jobs)),
		Records: make([]response.JobResponse, len(jobs)),
	}

	for i, job := range jobs {
		resp.Records[i] = *s.ConvertToJobResponse(&job, 0)
	}

	return resp, nil
}

// SearchByCondition 多条件搜索职位
func (s *JobService) SearchByCondition(conditions map[string]interface{}, page, size int) (*response.JobListResponse, error) {
	jobs, total, err := s.jobDao.SearchByCondition(conditions, page, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   total,
		Records: make([]response.JobResponse, len(jobs)),
	}

	for i, job := range jobs {
		resp.Records[i] = *s.ConvertToJobResponse(&job, 0)
	}

	return resp, nil
}

// GetExpiredJobs 获取已过期职位
func (s *JobService) GetExpiredJobs() (*response.JobListResponse, error) {
	jobs, err := s.jobDao.GetExpiredJobs()
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   int64(len(jobs)),
		Records: make([]response.JobResponse, len(jobs)),
	}

	for i, job := range jobs {
		resp.Records[i] = *s.ConvertToJobResponse(&job, 0)
	}

	return resp, nil
}

// SearchByCompany 获取公司发布的职位
func (s *JobService) SearchByCompany(companyID uint, page, size int) (*response.JobListResponse, error) {
	jobs, total, err := s.jobDao.SearchByCompany(companyID, page, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobListResponse{
		Total:   total,
		Records: make([]response.JobResponse, len(jobs)),
	}

	for i, job := range jobs {
		resp.Records[i] = *s.ConvertToJobResponse(&job, 0)
	}

	return resp, nil
}

// 其他辅助方法...
func getCacheKey(jobID uint) string {
	return fmt.Sprintf("job:%d", jobID)
}
