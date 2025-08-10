package response

import "time"

// JobApplyResponse 职位申请响应对象
// @Description 职位申请响应对象
type JobApplyResponse struct {
	ID            uint   `json:"id"`
	JobID         uint   `json:"jobId"`
	UserID        uint   `json:"userId"`
	ResumeID      uint   `json:"resumeId"`
	Status        int    `json:"status"`
	ApplyProgress string `json:"applyProgress"` // 申请进度状态
	// enum: 待处理,进行中,已接受,已拒绝,已撤回,待面试,面试通过,面试不通过,已发送Offer,Offer已接受,Offer已拒绝
	// example: 待面试
	ApplyTime time.Time `json:"applyTime"`
}

// JobApplyListResponse 职位申请列表响应
type JobApplyListResponse struct {
	Total   int64              `json:"total"`
	Records []JobApplyResponse `json:"records"`
}

// ToResponse 将 JobApply 转换为 JobApplyResponse
func (apply *JobApplyResponse) ToResponse() *JobApplyResponse {
	return &JobApplyResponse{
		ID:            apply.ID,
		JobID:         apply.JobID,
		UserID:        apply.UserID,
		ResumeID:      apply.ResumeID,
		Status:        apply.Status,
		ApplyProgress: apply.ApplyProgress,
		ApplyTime:     apply.ApplyTime,
	}
}
