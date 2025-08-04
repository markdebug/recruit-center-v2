package enums

// JobStatus 职位状态
type JobStatus int

const (
	JobStatusNormal  JobStatus = 1 // 正常
	JobStatusExpired JobStatus = 0 // 已过期
)

// String 实现 String 接口
func (s JobStatus) String() string {
	switch s {
	case JobStatusNormal:
		return "正常"
	case JobStatusExpired:
		return "已过期"
	default:
		return "未知状态"
	}
}

// JobType 职位类型
type JobType int

const (
	JobTypeFullTime JobType = 1 // 全职
	JobTypePartTime JobType = 2 // 兼职
	JobTypeIntern   JobType = 3 // 实习
)

func (t JobType) String() string {
	switch t {
	case JobTypeFullTime:
		return "全职"
	case JobTypePartTime:
		return "兼职"
	case JobTypeIntern:
		return "实习"
	default:
		return "未知类型"
	}
}

// JobCategory 职位分类
type JobCategory string

const (
	JobCategoryInternal JobCategory = "internal" // 内推
	JobCategoryDirect   JobCategory = "direct"   // 直招
	JobCategoryHunter   JobCategory = "hunter"   // 猎头
)

func (c JobCategory) String() string {
	switch c {
	case JobCategoryInternal:
		return "内推"
	case JobCategoryDirect:
		return "直招"
	case JobCategoryHunter:
		return "猎头"
	default:
		return "未知分类"
	}
}

func JobStatusFromInt(status int) JobStatus {
	return JobStatus(status)
}

func GetJobStatusText(status int) string {
	return JobStatusFromInt(status).String()
}

func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusNormal, JobStatusExpired:
		return true
	}
	return false
}

// JobType 相关方法
func JobTypeFromInt(typeStr int) JobType {
	return JobType(typeStr)
}

func GetJobTypeText(typeStr int) string {
	return JobTypeFromInt(typeStr).String()
}

func (t JobType) IsValid() bool {
	switch t {
	case JobTypeFullTime, JobTypePartTime, JobTypeIntern:
		return true
	}
	return false
}

// JobCategory 相关方法
func JobCategoryFromString(category string) JobCategory {
	return JobCategory(category)
}

func GetJobCategoryText(category string) string {
	return JobCategoryFromString(category).String()
}

func (c JobCategory) IsValid() bool {
	switch c {
	case JobCategoryInternal, JobCategoryDirect, JobCategoryHunter:
		return true
	}
	return false
}

// 获取所有职位状态
func ListJobStatus() []JobStatus {
	return []JobStatus{
		JobStatusNormal,
		JobStatusExpired,
	}
}

// 获取所有职位类型
func ListJobTypes() []JobType {
	return []JobType{
		JobTypeFullTime,
		JobTypePartTime,
		JobTypeIntern,
	}
}

// 获取所有职位分类
func ListJobCategories() []JobCategory {
	return []JobCategory{
		JobCategoryInternal,
		JobCategoryDirect,
		JobCategoryHunter,
	}
}
