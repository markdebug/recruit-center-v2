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
	jobApplyDAO         *dao.JobApplyDAO
	jobService          *JobService
	notificationService *NotificationService
}

// NewJobApplyService 创建职位申请服务实例
func NewJobApplyService(jobApplyDao *dao.JobApplyDAO, jobService *JobService, notificationService *NotificationService) *JobApplyService {
	return &JobApplyService{
		jobApplyDAO:         jobApplyDao,
		jobService:          jobService,
		notificationService: notificationService,
	}
}

// Create 创建职位申请
func (s *JobApplyService) Create(apply *model.JobApply) error {
	// 1. 验证职位是否存在且有效
	job, err := s.jobService.GetByID(apply.JobID, 0)
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
func (s *JobApplyService) ListByJob(jobID uint, page, size int) (*response.JobApplyListResponse, error) {
	applies, total, err := s.jobApplyDAO.ListByJob(jobID, page, size)
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

// VerifyApplyOwner 验证申请是否属于指定用户
func (s *JobApplyService) VerifyApplyOwner(applyID, userID uint) error {
	apply, err := s.jobApplyDAO.GetByID(applyID)
	if err != nil {
		return errors.Wrap(err, errors.NotFound)
	}
	if apply.UserID != userID {
		logger.L.Warn("非法操作：用户尝试操作非本人的申请记录",
			zap.Uint("applyID", applyID),
			zap.Uint("userID", userID),
			zap.Uint("ownerID", apply.UserID))
		return errors.New(errors.Forbidden)
	}
	return nil
}

// Delete 删除职位申请
func (s *JobApplyService) Delete(id, userID uint) error {
	// 验证操作权限
	if err := s.VerifyApplyOwner(id, userID); err != nil {
		return err
	}

	if err := s.jobApplyDAO.Delete(id); err != nil {
		logger.L.Error("删除职位申请失败",
			zap.Error(err),
			zap.Uint("id", id))
		return errors.Wrap(err, errors.InternalServerError)
	}
	return nil
}

// getStatusNotification 获取状态变更的通知内容
func (s *JobApplyService) getStatusNotification(apply *model.JobApply, status enums.JobApplyEnum) (userNotify, companyNotify *model.Notification) {
	job, _ := s.jobService.GetByID(apply.JobID, 0)
	jobName := "该职位"
	if job != nil {
		jobName = job.Name
	}

	// 根据不同状态生成不同的通知内容
	switch status {
	case enums.JobApplyWaitInterview:
		userNotify = &model.Notification{
			UserID:  apply.UserID,
			Title:   "面试通知",
			Content: fmt.Sprintf("您申请的 %s 已通过初筛，请等待面试安排", jobName),
			Type:    model.NotificationTypeInterview,
		}
		companyNotify = &model.Notification{
			UserID:  job.CompanyID,
			Title:   "候选人状态更新",
			Content: fmt.Sprintf("职位 %s 的候选人已进入面试环节，请及时安排面试", jobName),
			Type:    model.NotificationTypeStatusUpdate,
		}

	case enums.JobApplyInterviewPass:
		userNotify = &model.Notification{
			UserID:  apply.UserID,
			Title:   "面试结果通知",
			Content: fmt.Sprintf("恭喜！您在 %s 的面试已通过", jobName),
			Type:    model.NotificationTypeStatusUpdate,
		}
		companyNotify = &model.Notification{
			UserID:  job.CompanyID,
			Title:   "面试结果提醒",
			Content: fmt.Sprintf("职位 %s 的候选人面试已通过，请及时处理后续流程", jobName),
			Type:    model.NotificationTypeStatusUpdate,
		}

	case enums.JobApplyInterviewFail:
		userNotify = &model.Notification{
			UserID:  apply.UserID,
			Title:   "面试结果通知",
			Content: fmt.Sprintf("很遗憾，您在 %s 的面试未通过，欢迎继续投递其他职位", jobName),
			Type:    model.NotificationTypeStatusUpdate,
		}

	case enums.JobApplyOfferSent:
		userNotify = &model.Notification{
			UserID:  apply.UserID,
			Title:   "Offer通知",
			Content: fmt.Sprintf("恭喜！%s 向您发出了录用意向，请查看并确认", jobName),
			Type:    model.NotificationTypeStatusUpdate,
		}

	default:
		userNotify = &model.Notification{
			UserID:  apply.UserID,
			Title:   "申请状态更新",
			Content: fmt.Sprintf("您的职位申请状态已更新为：%s", status.String()),
			Type:    model.NotificationTypeStatusUpdate,
		}
	}

	return userNotify, companyNotify
}

// UpdateStatus 更新申请状态
func (s *JobApplyService) UpdateStatus(id uint, userID uint, status enums.JobApplyEnum) error {
	// 验证操作权限
	if err := s.VerifyApplyOwner(id, userID); err != nil {
		return err
	}

	// 1. 验证状态是否有效
	if !status.IsValid() {
		return fmt.Errorf("无效的状态值")
	}

	// 2. 检查申请是否存在
	apply, err := s.jobApplyDAO.GetByID(id)
	if err != nil {
		return errors.Wrap(err, errors.NotFound)
	}

	// 3. 验证状态流转是否合法
	if !s.isValidStatusTransition(apply.Status, int(status)) {
		return fmt.Errorf("无效的状态流转")
	}

	// 4. 更新状态
	if err := s.jobApplyDAO.UpdateStatus(id, int(status)); err != nil {
		logger.L.Error("更新申请状态失败",
			zap.Error(err),
			zap.Uint("id", id),
			zap.Int("status", int(status)))
		return err
	}

	// 5. 发送通知
	userNotify, companyNotify := s.getStatusNotification(apply, status)

	// 发送给求职者的通知
	if userNotify != nil {
		if err := s.notificationService.Create(userNotify); err != nil {
			logger.L.Error("发送求职者通知失败", zap.Error(err))
		}
	}

	// 发送给公司的通知
	if companyNotify != nil {
		if err := s.notificationService.Create(companyNotify); err != nil {
			logger.L.Error("发送公司通知失败", zap.Error(err))
		}
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

// List 获取职位申请列表
func (s *JobApplyService) List(page, size int) (*response.JobApplyListResponse, error) {
	applies, total, err := s.jobApplyDAO.List(page, size)
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

// ListByUserID 根据用户id，查询其全部的申请信息
func (s *JobApplyService) ListByUserID(userID uint, page, size int) (*response.JobApplyListResponse, error) {
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

// ListByCompanyID 根据公司信息，查询所有的职位申请记录
func (s *JobApplyService) ListByCompanyID(companyID uint, page, size int) (*response.JobApplyListResponse, error) {
	applies, total, err := s.jobApplyDAO.ListByCompany(companyID, page, size)
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
