package enums

// CommonEnum 定义了一些通用的枚举类型
type DeleteStatus int

const (
	DeleteStatusNormal  DeleteStatus = iota // 正常状态
	DeleteStatusDeleted                     // 已删除状态
)

func (s DeleteStatus) String() string {
	switch s {
	case DeleteStatusNormal:
		return "正常"
	case DeleteStatusDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}
