package utils

import "github.com/gofiber/fiber/v2"

// 统一状态码定义
const (
	// 成功状态码
	StatusSuccess = "success"

	// 错误状态码
	StatusError = "error"
)

// 统一错误代码
const (
	// 通用错误
	ErrBadRequest       = "BAD_REQUEST"        // 请求参数错误
	ErrUnauthorized     = "UNAUTHORIZED"       // 未授权
	ErrForbidden        = "FORBIDDEN"          // 禁止访问
	ErrNotFound         = "NOT_FOUND"          // 资源不存在
	ErrMethodNotAllowed = "METHOD_NOT_ALLOWED" // 方法不允许
	ErrServerError      = "SERVER_ERROR"       // 服务器错误

	// 业务错误
	ErrUserNotFound       = "USER_NOT_FOUND"       // 用户不存在
	ErrRoleNotFound       = "ROLE_NOT_FOUND"       // 角色不存在
	ErrPermissionNotFound = "PERMISSION_NOT_FOUND" // 权限不存在
	ErrInvalidCredentials = "INVALID_CREDENTIALS"  // 无效的凭证
	ErrDuplicateUsername  = "DUPLICATE_USERNAME"   // 用户名已存在
)

// 响应结构体
type Response struct {
	Status  string      `json:"status"`         // 状态: success 或 error
	Code    string      `json:"code,omitempty"` // 错误代码，仅当 status=error 时有值
	Message string      `json:"message"`        // 消息
	Data    interface{} `json:"data"`           // 数据，可以为 nil
}

// SuccessResponse 返回成功响应，HTTP状态码统一为200
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 返回错误响应，HTTP状态码统一为200
func ErrorResponse(c *fiber.Ctx, errCode string, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  StatusError,
		Code:    errCode,
		Message: message,
		Data:    nil,
	})
}

// BadRequestError 返回请求参数错误
func BadRequestError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrBadRequest, message)
}

// UnauthorizedError 返回未授权错误
func UnauthorizedError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrUnauthorized, message)
}

// ForbiddenError 返回禁止访问错误
func ForbiddenError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrForbidden, message)
}

// NotFoundError 返回资源不存在错误
func NotFoundError(c *fiber.Ctx, message string, code ...string) error {
	errCode := ErrNotFound
	if len(code) > 0 {
		errCode = code[0]
	}
	return ErrorResponse(c, errCode, message)
}

// ServerError 返回服务器错误
func ServerError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrServerError, message)
}
