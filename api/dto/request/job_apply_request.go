package request

import (
	"time"

	"org.thinkinai.com/recruit-center/internal/model"
)

// JobApply 职位申请实体
type JobApply struct {
	ID            uint           `json:"id"`
	JobID         uint           `json:"jobId"`         // 改为 jobId
	UserID        uint           `json:"userId"`        // 改为 userId
	ApplyTime     time.Time      `json:"applyTime"`     // 改为 applyTime
	ApplyProgress string         `json:"applyProgress"` // 改为 applyProgress
	Status        JobApplyStatus `json:"status"`
	CreateTime    time.Time      `json:"createTime"` // 改为 createTime
	UpdateTime    time.Time      `json:"updateTime"` // 改为 updateTime
}

// JobApplyStatus 申请状态枚举
type JobApplyStatus int

const (
	JobApplyStatusPending    JobApplyStatus = 1 // 待处理
	JobApplyStatusInProgress JobApplyStatus = 2 // 进行中
	JobApplyStatusAccepted   JobApplyStatus = 3 // 已接受
	JobApplyStatusRejected   JobApplyStatus = 4 // 已拒绝
	JobApplyStatusWithdrawn  JobApplyStatus = 5 // 已撤回
)

// String 返回状态的字符串表示
func (s JobApplyStatus) String() string {
	switch s {
	case JobApplyStatusPending:
		return "待处理"
	case JobApplyStatusInProgress:
		return "进行中"
	case JobApplyStatusAccepted:
		return "已接受"
	case JobApplyStatusRejected:
		return "已拒绝"
	case JobApplyStatusWithdrawn:
		return "已撤回"
	default:
		return "未知状态"
	}
}

// NewJobApply 创建新的职位申请
func NewJobApply(jobID, userID uint) *JobApply {
	return &JobApply{
		JobID:         jobID,
		UserID:        userID,
		ApplyTime:     time.Now(),
		Status:        JobApplyStatusPending,
		ApplyProgress: JobApplyStatusPending.String(),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
}

// canTransitionTo 检查是否可以转换到目标状态
func (ja *JobApply) canTransitionTo(targetStatus JobApplyStatus) bool {
	// 定义状态转换规则
	transitions := map[JobApplyStatus][]JobApplyStatus{
		JobApplyStatusPending: {
			JobApplyStatusInProgress,
			JobApplyStatusRejected,
			JobApplyStatusWithdrawn,
		},
		JobApplyStatusInProgress: {
			JobApplyStatusAccepted,
			JobApplyStatusRejected,
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
