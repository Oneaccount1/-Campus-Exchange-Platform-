package response

import (
	"errors" // 标准库errors包
	"net/http"

	appErrors "campus/internal/utils/errors" // 自定义errors包，使用别名
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PagedResponse 分页响应结构体
type PagedResponse struct {
	Response
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// 响应状态码常量
const (
	StatusSuccess = http.StatusOK
	StatusError   = http.StatusInternalServerError
)

// JSON 发送JSON响应
func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	JSON(c, StatusSuccess, "success", data)
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	JSON(c, StatusSuccess, message, data)
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	JSON(c, code, message, nil)
}

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	// 处理应用自定义错误
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		handleAppError(c, appErr)
		return
	}

	// 处理标准错误
	handleStandardError(c, err)
}

// 处理应用自定义错误
func handleAppError(c *gin.Context, err *appErrors.AppError) {
	switch err.Type {
	case appErrors.ErrorTypeNotFound:
		Fail(c, http.StatusNotFound, err.Message)
	case appErrors.ErrorTypeUnauthorized:
		Fail(c, http.StatusUnauthorized, err.Message)
	case appErrors.ErrorTypeForbidden:
		Fail(c, http.StatusForbidden, err.Message)
	case appErrors.ErrorTypeBadRequest, appErrors.ErrorTypeDuplicate, appErrors.ErrorTypeValidation:
		Fail(c, http.StatusBadRequest, err.Message)
	default:
		if err.Err != nil {
			Fail(c, http.StatusInternalServerError, err.Message+": "+err.Err.Error())
		} else {
			Fail(c, http.StatusInternalServerError, err.Message)
		}
	}
}

// 处理标准错误
func handleStandardError(c *gin.Context, err error) {
	switch {
	case appErrors.IsNotFound(err):
		Fail(c, http.StatusNotFound, err.Error())
	case appErrors.IsUnauthorized(err):
		Fail(c, http.StatusUnauthorized, err.Error())
	case appErrors.IsForbidden(err):
		Fail(c, http.StatusForbidden, err.Error())
	case appErrors.IsBadRequest(err), appErrors.IsDuplicate(err), appErrors.IsValidation(err):
		Fail(c, http.StatusBadRequest, err.Error())
	default:
		Fail(c, http.StatusInternalServerError, "服务器内部错误: "+err.Error())
	}
}

// 以下是常用的HTTP状态响应函数，保留为便捷方法

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, message)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, message)
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, message)
}

// ServerError 500错误响应
func ServerError(c *gin.Context, message string, err error) {
	if err != nil {
		Fail(c, http.StatusInternalServerError, message+": "+err.Error())
	} else {
		Fail(c, http.StatusInternalServerError, message)
	}
}
