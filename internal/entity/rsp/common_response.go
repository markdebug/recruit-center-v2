package rsp

import "org.thinkinai.com/recruit-center/pkg/enums"

// Response 通用响应结构
type Response struct {
	Code    enums.ErrorCode `json:"code"`           // 状态码
	Message string          `json:"message"`        // 响应信息
	Data    interface{}     `json:"data,omitempty"` // 响应数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    enums.ErrorCode `json:"code"`           // 状态码
	Message string          `json:"message"`        // 响应信息
	Data    interface{}     `json:"data,omitempty"` // 响应数据
	Total   int64           `json:"total"`          // 总记录数
	Page    int             `json:"page"`           // 当前页码
	Size    int             `json:"size"`           // 每页大小
}

// NewSuccess 创建成功响应
func NewSuccess(data interface{}) *Response {
	return &Response{
		Code:    enums.Success,
		Message: enums.Success.String(),
		Data:    data,
	}
}

// NewError 创建错误响应
func NewError(code enums.ErrorCode) *Response {
	return &Response{
		Code:    code,
		Message: code.String(),
	}
}

// NewErrorWithMsg 创建带自定义消息的错误响应
func NewErrorWithMsg(code enums.ErrorCode, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

// NewPage 创建分页响应
func NewPage(data interface{}, total int64, page, size int) *PageResponse {
	return &PageResponse{
		Code:    enums.Success,
		Message: enums.Success.String(),
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	}
}
