package enums

import "fmt"

//简历隐私枚举
type ResumeAccessEnum int

const (
	Hide ResumeAccessEnum = 1 // 简历隐藏
	Open ResumeAccessEnum = 2 // 简历公开
)

func (e ResumeAccessEnum) String() string {
	switch e {
	case Hide:
		return "简历隐藏"
	case Open:
		return "简历公开"
	default:
		return "未知状态"
	}
}

//在职状态枚举
type WorkingStatusEnum int

const (
	OnTheJob  WorkingStatusEnum = 1 // 在职
	Dimission WorkingStatusEnum = 2 // 离职
)

func (e WorkingStatusEnum) String() string {
	switch e {
	case OnTheJob:
		return "在职"
	case Dimission:
		return "离职"
	default:
		return "未知状态"
	}
}

// ParseResumeAccess 解析简历隐私枚举
func ParseResumeAccess(v int) (ResumeAccessEnum, error) {
	e := ResumeAccessEnum(v)
	if !e.IsValid() {
		return 0, fmt.Errorf("无效的简历隐私状态值: %d", v)
	}
	return e, nil
}

// ParseWorkingStatus 解析在职状态枚举
func ParseWorkingStatus(v int) (WorkingStatusEnum, error) {
	e := WorkingStatusEnum(v)
	if !e.IsValid() {
		return 0, fmt.Errorf("无效的在职状态值: %d", v)
	}
	return e, nil
}
func (e ResumeAccessEnum) IsValid() bool {
	switch e {
	case Hide, Open:
		return true
	default:
		return false
	}
}

func (e WorkingStatusEnum) IsValid() bool {
	switch e {
	case OnTheJob, Dimission:
		return true
	default:
		return false
	}
}
