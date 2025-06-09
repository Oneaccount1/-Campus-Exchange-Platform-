package errors

import (
	"errors"
	"fmt"
)

// 预定义错误类型
var (
	// ErrNotFound 表示资源未找到
	ErrNotFound = errors.New("资源未找到")

	// ErrUnauthorized 表示未授权
	ErrUnauthorized = errors.New("未授权访问")

	// ErrForbidden 表示禁止访问
	ErrForbidden = errors.New("禁止访问")

	// ErrBadRequest 表示请求参数错误
	ErrBadRequest = errors.New("请求参数错误")

	// ErrInternalServer 表示服务器内部错误
	ErrInternalServer = errors.New("服务器内部错误")

	// ErrDuplicateRecord 表示记录已存在
	ErrDuplicateRecord = errors.New("记录已存在")

	// ErrValidation 表示验证错误
	ErrValidation = errors.New("验证错误")
)

// ErrorType 错误类型
type ErrorType string

// 错误类型常量
const (
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeInternalServer ErrorType = "INTERNAL_SERVER"
	ErrorTypeDuplicate      ErrorType = "DUPLICATE"
	ErrorTypeValidation     ErrorType = "VALIDATION"
)

// AppError 应用错误结构体
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用错误
func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError 创建资源未找到错误
func NewNotFoundError(resource string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s未找到", resource),
		Err:     err,
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, err error) *AppError {
	if message == "" {
		message = "未授权访问"
	}
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Err:     err,
	}
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, err error) *AppError {
	if message == "" {
		message = "禁止访问"
	}
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Err:     err,
	}
}

// NewBadRequestError 创建请求参数错误
func NewBadRequestError(message string, err error) *AppError {
	if message == "" {
		message = "请求参数错误"
	}
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewInternalServerError 创建服务器内部错误
func NewInternalServerError(message string, err error) *AppError {
	if message == "" {
		message = "服务器内部错误"
	}
	return &AppError{
		Type:    ErrorTypeInternalServer,
		Message: message,
		Err:     err,
	}
}

// NewDuplicateError 创建记录已存在错误
func NewDuplicateError(resource string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDuplicate,
		Message: fmt.Sprintf("%s已存在", resource),
		Err:     err,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string, err error) *AppError {
	if message == "" {
		message = "验证错误"
	}
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

// IsNotFound 判断是否为资源未找到错误
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeNotFound
	}
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized 判断是否为未授权错误
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeUnauthorized
	}
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden 判断是否为禁止访问错误
func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeForbidden
	}
	return errors.Is(err, ErrForbidden)
}

// IsBadRequest 判断是否为请求参数错误
func IsBadRequest(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeBadRequest
	}
	return errors.Is(err, ErrBadRequest)
}

// IsInternalServer 判断是否为服务器内部错误
func IsInternalServer(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeInternalServer
	}
	return errors.Is(err, ErrInternalServer)
}

// IsDuplicate 判断是否为记录已存在错误
func IsDuplicate(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeDuplicate
	}
	return errors.Is(err, ErrDuplicateRecord)
}

// IsValidation 判断是否为验证错误
func IsValidation(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeValidation
	}
	return errors.Is(err, ErrValidation)
}
