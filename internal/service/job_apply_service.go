package service

import (
	"fmt"

	"go.uber.org/zap"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/dao"
	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
	"org.thinkinai.com/recruit-center/pkg/errors"
	"org.thinkinai.com/recruit-center/pkg/logger"
)

// JobApplyService 职位申请服务
type JobApplyService struct {
	jobApplyDAO *dao.JobApplyDAO
	jobService  *JobService
}

// NewJobApplyService 创建职位申请服务实例
func NewJobApplyService(jobApplyDao *dao.JobApplyDAO, jobService *JobService) *JobApplyService {
	return &JobApplyService{
		jobApplyDAO: jobApplyDao,
		jobService:  jobService,
	}
}

// Create 创建职位申请
func (s *JobApplyService) Create(apply *model.JobApply) error {
	// 1. 验证职位是否存在且有效
	job, err := s.jobService.GetByID(apply.JobID)
	if err != nil {
		return err
	}
	if !job.IsActive() {
		return errors.New(errors.JobExpired)
	}

	// 2. 检查是否重复申请
	exists, _ := s.jobApplyDAO.GetByUserAndJob(apply.UserID, apply.JobID)
	if exists != nil {
		return errors.New(errors.JobAlreadyApplied)
	}

	// 3. 设置初始状态
	apply.Status = int(enums.JobApplyPending)
	apply.ApplyProgress = enums.JobApplyPending.String()

	// 4. 创建申请记录
	if err := s.jobApplyDAO.Create(apply); err != nil {
		logger.L.Error("创建职位申请失败",
			zap.Error(err),
			zap.Uint("user_id", apply.UserID),
			zap.Uint("job_id", apply.JobID))
		return err
	}

	return nil
}

// convertToJobApplyResponse 将 model.JobApply 转换为 response.JobApplyResponse
func (s *JobApplyService) ConvertToJobApplyResponse(apply *model.JobApply) *response.JobApplyResponse {
	if apply == nil {
		return nil
	}

	resp := &response.JobApplyResponse{
		ID:            apply.ID,
		JobID:         apply.JobID,
		UserID:        apply.UserID,
		ResumeID:      apply.ResumeID,
		Status:        apply.Status,
		ApplyProgress: apply.ApplyProgress,
		ApplyTime:     apply.ApplyTime,
	}

	return resp
}

// GetByID 获取申请详情
func (s *JobApplyService) GetByID(id uint) (*response.JobApplyResponse, error) {
	apply, err := s.jobApplyDAO.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ConvertToJobApplyResponse(apply), nil
}

// ListByUser 获取用户的申请列表
func (s *JobApplyService) ListByUser(userID uint, page, size int) (*response.JobApplyListResponse, error) {
	applies, total, err := s.jobApplyDAO.ListByUser(userID, page, size)
	if err != nil {
		return nil, err
	}

	resp := &response.JobApplyListResponse{
		Total:   total,
		Records: make([]response.JobApplyResponse, len(applies)),
	}

	for i, apply := range applies {
		resp.Records[i] = *s.ConvertToJobApplyResponse(&apply)
	}

	return resp, nil
}

// ListByJob 获取职位的申请列表
func (s *JobApplyService) ListByJob(jobID uint, page, size int) ([]model.JobApply, int64, error) {
	return s.jobApplyDAO.ListByJob(jobID, page, size)
}

// UpdateStatus 更新申请状态
func (s *JobApplyService) UpdateStatus(id uint, status int) error {
	// 1. 验证状态是否有效
	if !enums.JobApplyEnum(status).IsValid() {
		return fmt.Errorf("无效的状态值")
	}

	// 2. 检查申请是否存在
	_, err := s.jobApplyDAO.GetByID(id)
	if err != nil {
		return errors.Wrap(err, errors.NotFound)
	}

	// 3. 验证状态流转是否合法
	// if !apply.CanTransitionTo(status) {
	// 	return errors.New(errors.BadRequest).WithMessage("无效的状态转换")
	// }

	// 4. 更新状态
	if err := s.jobApplyDAO.UpdateStatus(id, status); err != nil {
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

// Delete 删除职位申请
func (s *JobApplyService) Delete(id uint) error {
	if err := s.jobApplyDAO.Delete(id); err != nil {
		logger.L.Error("删除职位申请失败",
			zap.Error(err),
			zap.Uint("id", id))
		return errors.Wrap(err, errors.InternalServerError)
	}
	return nil
}

// List 获取职位申请列表
func (s *JobApplyService) List(page, size int) ([]model.JobApply, int64, error) {
	return s.jobApplyDAO.List(page, size)
}
