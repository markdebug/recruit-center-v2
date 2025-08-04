package service

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// JobApplyService 职位申请服务
type JobApplyService struct {
	jobApplyDao *dao.JobApplyDao
	jobDao      *dao.JobDao // 用于验证职位信息
}

// NewJobApplyService 创建职位申请服务实例
func NewJobApplyService(jobApplyDao *dao.JobApplyDao, jobDao *dao.JobDao) *JobApplyService {
	return &JobApplyService{
		jobApplyDao: jobApplyDao,
		jobDao:      jobDao,
	}
}

// Create 创建职位申请
func (s *JobApplyService) Create(apply *dao.JobApply) error {
	// 1. 验证职位是否存在且有效
	job, err := s.jobDao.GetByID(apply.JobID)
	if err != nil {
		return fmt.Errorf("职位不存在")
	}
	if job.JobExpireTime.Before(time.Now()) {
		return fmt.Errorf("职位已过期")
	}

	// 2. 检查是否重复申请
	exist, err := s.jobApplyDao.GetByUserAndJob(apply.UserID, apply.JobID)
	if err == nil && exist != nil {
		return fmt.Errorf("您已经申请过该职位")
	}

	// 3. 设置初始状态
	apply.Status = int(enums.JobApplyPending)
	apply.ApplyProgress = enums.JobApplyPending.String()

	// 4. 创建申请记录
	if err := s.jobApplyDao.Create(apply); err != nil {
		logger.L.Error("创建职位申请失败",
			zap.Error(err),
			zap.Uint("user_id", apply.UserID),
			zap.Uint("job_id", apply.JobID))
		return err
	}

	return nil
}

// GetByID 获取申请详情
func (s *JobApplyService) GetByID(id uint) (*dao.JobApply, error) {
	return s.jobApplyDao.GetByID(id)
}

// ListByUser 获取用户的申请列表
func (s *JobApplyService) ListByUser(userID uint, page, size int) ([]dao.JobApply, int64, error) {
	return s.jobApplyDao.ListByUser(userID, page, size)
}

// ListByJob 获取职位的申请列表
func (s *JobApplyService) ListByJob(jobID uint, page, size int) ([]dao.JobApply, int64, error) {
	return s.jobApplyDao.ListByJob(jobID, page, size)
}

// UpdateStatus 更新申请状态
func (s *JobApplyService) UpdateStatus(id uint, status int) error {
	// 1. 验证状态是否有效
	if !enums.JobApplyEnum(status).IsValid() {
		return fmt.Errorf("无效的状态值")
	}

	// 2. 检查申请是否存在
	apply, err := s.jobApplyDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("申请记录不存在")
	}

	// 3. 验证状态流转是否合法
	if !s.isValidStatusTransition(apply.Status, status) {
		return fmt.Errorf("非法的状态变更")
	}

	// 4. 更新状态
	if err := s.jobApplyDao.UpdateStatus(id, status); err != nil {
		logger.L.Error("更新申请状态失败",
			zap.Error(err),
			zap.Uint("id", id),
			zap.Int("status", status))
		return err
	}

	return nil
}

// isValidStatusTransition 验证状态流转是否合法
func (s *JobApplyService) isValidStatusTransition(from, to int) bool {
	// 定义状态流转规则
	transitions := map[int][]int{
		int(enums.JobApplyPending): {
			int(enums.JobApplyInProgress),
			int(enums.JobApplyRejected),
			int(enums.JobApplyWithdrawn),
		},
		int(enums.JobApplyInProgress): {
			int(enums.JobApplyAccepted),
			int(enums.JobApplyRejected),
		},
	}

	// 检查状态流转是否允许
	if allowedStatus, exists := transitions[from]; exists {
		for _, status := range allowedStatus {
			if status == to {
				return true
			}
		}
	}
	return false
}
