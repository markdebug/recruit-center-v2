package enums

// JobApplyEnum defines the enumeration for job application statuses.

type JobApplyEnum int

const (
	JobApplyPending    JobApplyEnum = 1 // 待处理
	JobApplyInProgress JobApplyEnum = 2 // 进行中
	JobApplyAccepted   JobApplyEnum = 3 // 已接受
	JobApplyRejected   JobApplyEnum = 4 // 已拒绝
	JobApplyWithdrawn  JobApplyEnum = 5 // 已撤回
)

func (e JobApplyEnum) String() string {
	switch e {
	case JobApplyPending:
		return "待处理"
	case JobApplyInProgress:
		return "进行中"
	case JobApplyAccepted:
		return "已接受"
	case JobApplyRejected:
		return "已拒绝"
	case JobApplyWithdrawn:
		return "已撤回"
	default:
		return "未知状态"
	}
}

// FromInt 将int类型转换为JobApplyEnum
func FromInt(status int) JobApplyEnum {
	return JobApplyEnum(status)
}

// GetStatusText 获取状态文本
func GetStatusText(status int) string {
	return FromInt(status).String()
}

// IsValid 检查状态值是否有效
func (e JobApplyEnum) IsValid() bool {
	switch e {
	case JobApplyPending, JobApplyInProgress, JobApplyAccepted,
		JobApplyRejected, JobApplyWithdrawn:
		return true
	}
	return false
}
