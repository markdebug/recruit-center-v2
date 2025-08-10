package request

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
	"org.thinkinai.com/recruit-center/pkg/enums"
)

// JobApply 职位申请实体
type JobApply struct {
	ID            uint               `json:"id"`
	JobID         uint               `json:"jobId"`         // 改为 jobId
	UserID        uint               `json:"userId"`        // 改为 userId
	ApplyTime     time.Time          `json:"applyTime"`     // 改为 applyTime
	ApplyProgress string             `json:"applyProgress"` // 改为 applyProgress
	Reason        string             `json:"reason"`        //拒绝原因
	Status        enums.JobApplyEnum `json:"status"`
	CreateTime    time.Time          `json:"createTime"` // 改为 createTime
	UpdateTime    time.Time          `json:"updateTime"` // 改为 updateTime
}

type JobApplyUpdateStatus struct {
	Status enums.JobApplyEnum `json:"status" binding:"required"`
	JobID  uint               `json:"jobId" binding:"required"`
	UsID   uint               `json:"userId" binding:"required"`
	Reason string             `json:"reason"`
}

// NewJobApply 创建新的职位申请
func NewJobApply(jobID, userID uint) *JobApply {
	return &JobApply{
		JobID:         jobID,
		UserID:        userID,
		ApplyTime:     time.Now(),
		Status:        enums.JobApplyPending,
		ApplyProgress: enums.JobApplyPending.String(),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
}

// canTransitionTo 检查是否可以转换到目标状态
func (ja *JobApply) canTransitionTo(targetStatus enums.JobApplyEnum) bool {
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

// ToModel 将 DTO 转换为数据模型
func (ja *JobApply) ToModel() *model.JobApply {
	return &model.JobApply{
		ID:            ja.ID,
		JobID:         ja.JobID,
		UserID:        ja.UserID,
		ApplyTime:     ja.ApplyTime,
		ApplyProgress: ja.ApplyProgress,
		Status:        int(ja.Status),
		CreateTime:    ja.CreateTime,
		UpdateTime:    ja.UpdateTime,
	}
}
