package errors

// ErrorCode 错误码
type ErrorCode int

const (
	// 成功
	Success ErrorCode = 0000

	// 客户端错误 (1xxx)
	BadRequest       ErrorCode = 1001 // 请求参数错误
	Unauthorized     ErrorCode = 1002 // 未授权
	Forbidden        ErrorCode = 1003 // 禁止访问
	NotFound         ErrorCode = 1004 // 资源不存在
	MethodNotAllowed ErrorCode = 1005 // 方法不允许
	Conflict         ErrorCode = 1006 // 资源冲突
	TooManyRequests  ErrorCode = 1007 // 请求过多
	InvalidParams    ErrorCode = 1008 // 无效的请求参数

	// 服务器错误 (9xxx)
	InternalServerError ErrorCode = 9000 // 服务器内部错误
	ServiceUnavailable  ErrorCode = 9001 // 服务不可用

	// 业务错误码 (1001-8999)
	// 用户模块 (6001-6999)
	UserNotFound      ErrorCode = 6001 // 用户不存在
	UserAlreadyExists ErrorCode = 6002 // 用户已存在
	PasswordIncorrect ErrorCode = 6003 // 密码错误

	// 职位模块 (2001-2999)
	JobNotFound                   ErrorCode = 2001 // 职位不存在
	JobExpired                    ErrorCode = 2002 // 职位已过期
	JobAlreadyApplied             ErrorCode = 2003 // 已经申请过该职位
	JobApplicationLimitReached    ErrorCode = 2004 // 职位申请次数已达上限
	JobCreationLimitReached       ErrorCode = 2005 // 职位创建次数已达上限
	JobUpdateNotAllowed           ErrorCode = 2006 // 职位更新不允许
	JobDeleteNotAllowed           ErrorCode = 2007 // 职位删除不允许
	JobApplicationNotFound        ErrorCode = 2008 // 职位申请不存在
	JobApplicationAlreadyReviewed ErrorCode = 2009 // 职位申请已被审核
	JobNotBelongToCompany         ErrorCode = 2010 // 职位不属于该公司
	InvalidJobStatus              ErrorCode = 2011 // 无效的职位状态

	// 公司模块 (3001-3999)
	CompanyNotFound ErrorCode = 3001 // 公司不存在
	CompanyInactive ErrorCode = 3002 // 公司未激活

	//文件上传模块 (4001-4999)
	FileTooLarge       ErrorCode = 4001 // 文件过大
	FileTypeNotAllowed ErrorCode = 4002 // 文件类型不允许
	InvalidFileFormat  ErrorCode = 4003 // 无效的文件格式
	TooManyFiles       ErrorCode = 4004 // 文件数量超过限制
	FileUploadFailed   ErrorCode = 4005 // 文件上传失败

	//简历模块 (5001-5999)
	ResumeNotFound          ErrorCode = 5001 // 简历不存在
	ResumeInvalid           ErrorCode = 5002 // 简历无效
	ResumeExists            ErrorCode = 5003 // 简历已存在
	ResumeTooLarge          ErrorCode = 5004 // 简历文件过大
	ResumeFormat            ErrorCode = 5005 // 简历格式错误
	ResumeUpload            ErrorCode = 5006 // 简历上传失败
	ResumeDelete            ErrorCode = 5007 // 简历删除失败
	ResumeAccessDenied      ErrorCode = 5008 // 无权限访问简历
	ResumeParse             ErrorCode = 5009 // 简历解析失败
	ResumeUpdate            ErrorCode = 5010 // 简历更新失败
	ResumeUpdateTooFrequent ErrorCode = 5011 // 简历更新过于频繁
	ResumeUpdateStatus      ErrorCode = 5012 // 简历更新状态错误

	// 外部模块 (8001-8999)

	InvalidToken ErrorCode = 8001 // 无效的令牌

)

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
	case FileTooLarge:
		return "文件过大"
	case FileTypeNotAllowed:
		return "文件类型不允许"
	case InvalidFileFormat:
		return "无效的文件格式"
	case MethodNotAllowed:
		return "方法不允许"
	case Conflict:
		return "资源冲突"
	case TooManyRequests:
		return "请求过多"
	case ServiceUnavailable:
		return "服务不可用"
	case InvalidParams:
		return "无效的请求参数"
	case JobApplicationLimitReached:
		return "职位申请次数已达上限"
	case JobCreationLimitReached:
		return "职位创建次数已达上限"
	case JobUpdateNotAllowed:
		return "职位更新不允许"
	case JobDeleteNotAllowed:
		return "职位删除不允许"
	case JobApplicationNotFound:
		return "职位申请不存在"
	case JobApplicationAlreadyReviewed:
		return "职位申请已被审核"
	case ResumeNotFound:
		return "简历不存在"
	case ResumeInvalid:
		return "简历无效"
	case ResumeExists:
		return "简历已存在"
	case ResumeTooLarge:
		return "简历文件过大"
	case ResumeFormat:
		return "简历格式错误"
	case ResumeUpload:
		return "简历上传失败"
	case ResumeDelete:
		return "简历删除失败"
	case ResumeAccessDenied:
		return "无权限访问简历"
	case ResumeParse:
		return "简历解析失败"
	case ResumeUpdate:
		return "简历更新失败"
	case ResumeUpdateTooFrequent:
		return "简历更新过于频繁"
	case ResumeUpdateStatus:
		return "简历更新状态错误"

	default:
		return "未知错误"
	}
}
