package enums

// ErrorCode 错误码
type ErrorCode int

const (
	// 成功
	Success ErrorCode = 200

	// 客户端错误 (4xx)
	BadRequest       ErrorCode = 400 // 请求参数错误
	Unauthorized     ErrorCode = 401 // 未授权
	Forbidden        ErrorCode = 403 // 禁止访问
	NotFound         ErrorCode = 404 // 资源不存在
	MethodNotAllowed ErrorCode = 405 // 方法不允许
	Conflict         ErrorCode = 409 // 资源冲突
	TooManyRequests  ErrorCode = 429 // 请求过多

	// 服务器错误 (5xx)
	InternalServerError ErrorCode = 500 // 服务器内部错误
	ServiceUnavailable  ErrorCode = 503 // 服务不可用

	// 业务错误码 (1001-9999)
	// 用户模块 (1001-1999)
	UserNotFound      ErrorCode = 1001 // 用户不存在
	UserAlreadyExists ErrorCode = 1002 // 用户已存在
	PasswordIncorrect ErrorCode = 1003 // 密码错误

	// 职位模块 (2001-2999)
	JobNotFound       ErrorCode = 2001 // 职位不存在
	JobExpired        ErrorCode = 2002 // 职位已过期
	JobAlreadyApplied ErrorCode = 2003 // 已经申请过该职位

	// 公司模块 (3001-3999)
	CompanyNotFound ErrorCode = 3001 // 公司不存在
	CompanyInactive ErrorCode = 3002 // 公司未激活
)

// Error 定义错误信息结构
type Error struct {
	Code    ErrorCode `json:"code"`    // 错误码
	Message string    `json:"message"` // 错误信息
}

// String 获取错误码对应的文本描述
func (e ErrorCode) String() string {
	switch e {
	case Success:
		return "成功"
	case BadRequest:
		return "请求参数错误"
	case Unauthorized:
		return "未授权"
	case Forbidden:
		return "禁止访问"
	case NotFound:
		return "资源不存在"
	case InternalServerError:
		return "服务器内部错误"
	case UserNotFound:
		return "用户不存在"
	case UserAlreadyExists:
		return "用户已存在"
	case PasswordIncorrect:
		return "密码错误"
	case JobNotFound:
		return "职位不存在"
	case JobExpired:
		return "职位已过期"
	case JobAlreadyApplied:
		return "已经申请过该职位"
	case CompanyNotFound:
		return "公司不存在"
	case CompanyInactive:
		return "公司未激活"
	default:
		return "未知错误"
	}
}

// NewError 创建新的错误
func NewError(code ErrorCode) *Error {
	return &Error{
		Code:    code,
		Message: code.String(),
	}
}

// WithMessage 自定义错误信息
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}
