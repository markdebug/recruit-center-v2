package request

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobApplyRequest 职位申请请求
type JobApplyRequest struct {
	JobID    uint `json:"jobId" binding:"required"`
	UserID   uint `json:"userId" binding:"required"`
	ResumeID uint `json:"resumeId" binding:"required"`
}

type JobApplyUpdateStatus struct {
	Status enums.JobApplyEnum `json:"status" binding:"required"`
	JobID  uint               `json:"jobId" binding:"required"`
	UsID   uint               `json:"userId" binding:"required"`
	Reason string             `json:"reason"`
}

// NewJobApply 创建新的职位申请
func (ja *JobApplyRequest) NewJobApply() *model.JobApply {
	return &model.JobApply{
		JobID:         ja.JobID,
		ResumeID:      ja.ResumeID,
		UserID:        ja.UserID,
		ApplyTime:     time.Now(),
		Status:        int(enums.DeleteStatusNormal),
		ApplyProgress: enums.JobApplyPending.String(),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
}

// canTransitionTo 检查是否可以转换到目标状态
func (ja *JobApplyUpdateStatus) canTransitionTo(targetStatus enums.JobApplyEnum) bool {
	// 定义状态转换规则
	transitions := map[enums.JobApplyEnum][]enums.JobApplyEnum{
		enums.JobApplyPending: {
			enums.JobApplyInProgress,
			enums.JobApplyRejected,
			enums.JobApplyWithdrawn,
		},
		enums.JobApplyInProgress: {
			enums.JobApplyAccepted,
			enums.JobApplyRejected,
		},
		// 已接受、已拒绝、已撤回状态为终态，不能再转换
	}

	allowedStatuses, exists := transitions[ja.Status]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowedStatuses {
		if allowedStatus == targetStatus {
			return true
		}
	}
	return false
}

// 校验创建参数是否正确
func (ja *JobApplyRequest) Validate() bool {
	if ja.JobID == 0 || ja.UserID == 0 || ja.ResumeID == 0 {
		return false
	}
	return true
}
