package response

import "time"

// JobApplyResponse 职位申请响应对象
type JobApplyResponse struct {
	ID            uint      `json:"id"`
	JobID         uint      `json:"jobId"`
	UserID        uint      `json:"userId"`
	ResumeID      uint      `json:"resumeId"`
	Status        int       `json:"status"`
	ApplyProgress string    `json:"applyProgress"`
	ApplyTime     time.Time `json:"applyTime"`

	// 关联数据
	Job    *JobResponse    `json:"job,omitempty"`
	Resume *ResumeResponse `json:"resume,omitempty"`
}

// JobApplyListResponse 职位申请列表响应
type JobApplyListResponse struct {
	Total   int64              `json:"total"`
	Records []JobApplyResponse `json:"records"`
}
