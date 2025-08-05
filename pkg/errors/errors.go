package errors

import (
	"fmt"
)

// Error 自定义错误结构
type Error struct {
	Code    ErrorCode   `json:"code"`              // 错误码
	Message string      `json:"message"`           // 错误信息
	Details interface{} `json:"details,omitempty"` // 错误详情
	Err     error       `json:"-"`                 // 原始错误
}

// New 创建新的错误
func New(code ErrorCode) *Error {
	return &Error{
		Code:    code,
		Message: code.String(),
	}
}

// Wrap 包装已有错误
func Wrap(err error, code ErrorCode) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:    code,
		Message: code.String(),
		Err:     err,
	}
}

// WithMessage 设置自定义错误信息
func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

// WithDetails 设置错误详情
func (e *Error) WithDetails(details interface{}) *Error {
	e.Details = details
	return e
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("错误码: %d, 错误信息: %s, 原始错误: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}
