package enums

// JobApplyEnum defines the enumeration for job application statuses.

type JobApplyEnum int

const (
	JobApplyPending    JobApplyEnum = 1 // 待处理
	JobApplyInProgress JobApplyEnum = 2 // 进行中
	JobApplyAccepted   JobApplyEnum = 3 // 已接受
	JobApplyRejected   JobApplyEnum = 4 // 已拒绝
	JobApplyWithdrawn  JobApplyEnum = 5 // 已撤回

	// 新增面试相关状态
	JobApplyWaitInterview JobApplyEnum = 6  // 待面试
	JobApplyInterviewPass JobApplyEnum = 7  // 面试通过
	JobApplyInterviewFail JobApplyEnum = 8  // 面试不通过
	JobApplyOfferSent     JobApplyEnum = 9  // 已发送Offer
	JobApplyOfferAccept   JobApplyEnum = 10 // Offer已接受
	JobApplyOfferReject   JobApplyEnum = 11 // Offer已拒绝
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
	case JobApplyWaitInterview:
		return "待面试"
	case JobApplyInterviewPass:
		return "面试通过"
	case JobApplyInterviewFail:
		return "面试不通过"
	case JobApplyOfferSent:
		return "已发送Offer"
	case JobApplyOfferAccept:
		return "Offer已接受"
	case JobApplyOfferReject:
		return "Offer已拒绝"
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
		JobApplyRejected, JobApplyWithdrawn,
		JobApplyWaitInterview, JobApplyInterviewPass, JobApplyInterviewFail,
		JobApplyOfferSent, JobApplyOfferAccept, JobApplyOfferReject:
		return true
	}
	return false
}
