package service

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/entity/rsp"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// JobService 职位服务
type JobService struct {
	jobDao *dao.JobDao
}

// NewJobService 创建职位服务实例
func NewJobService(jobDao *dao.JobDao) *JobService {
	return &JobService{
		jobDao: jobDao,
	}
}

// Create 创建职位
func (s *JobService) Create(job *dao.Job) error {
	// 参数校验
	if job.Name == "" {
		return fmt.Errorf("职位名称不能为空")
	}
	if job.CompanyID == 0 {
		return fmt.Errorf("公司ID不能为空")
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

// Update 更新职位信息
func (s *JobService) Update(job *dao.Job) error {
	// 检查职位是否存在
	exist, err := s.jobDao.GetByID(job.ID)
	if err != nil {
		return err
	}
	if exist == nil {
		return fmt.Errorf("职位不存在")
	}

	return s.jobDao.Update(job)
}

// Delete 删除职位
func (s *JobService) Delete(id uint) error {
	return s.jobDao.Delete(id)
}

// GetByID 获取职位详情
func (s *JobService) GetByID(id uint) (*dao.Job, error) {
	return s.jobDao.GetByID(id)
}

// List 获取职位列表
func (s *JobService) List(page, size int) ([]dao.Job, int64, error) {
	return s.jobDao.List(page, size)
}

// SearchByKeyword 关键词搜索职位
func (s *JobService) SearchByKeyword(keyword string) ([]dao.Job, error) {
	return s.jobDao.SearchByKeyword(keyword)
}

// SearchByCondition 多条件搜索职位
func (s *JobService) SearchByCondition(conditions map[string]interface{}, page, size int) (*rsp.PageResponse, error) {
	jobs, total, err := s.jobDao.SearchByCondition(conditions, page, size)
	if err != nil {
		return nil, err
	}

	return rsp.NewPage(jobs, total, page, size), nil
}

// UpdateStatus 更新职位状态
func (s *JobService) UpdateStatus(id uint, status int) error {
	// 检查状态是否有效
	if !enums.JobStatus(status).IsValid() {
		return fmt.Errorf("无效的职位状态")
	}

	return s.jobDao.UpdateStatus(id, status)
}

// GetExpiredJobs 获取已过期职位
func (s *JobService) GetExpiredJobs() ([]dao.Job, error) {
	return s.jobDao.GetExpiredJobs()
}

// SearchByCompany 获取公司发布的职位
func (s *JobService) SearchByCompany(companyID uint, page, size int) (*rsp.PageResponse, error) {
	jobs, total, err := s.jobDao.SearchByCompany(companyID, page, size)
	if err != nil {
		return nil, err
	}

	return rsp.NewPage(jobs, total, page, size), nil
}
