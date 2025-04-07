package utils

import "github.com/gofiber/fiber/v2"

// 统一状态码定义
const (
	// 成功状态码
	StatusSuccess = "success"

	// 错误状态码
	StatusError = "error"
)

// 统一错误代码 - 使用数值
const (
	// 通用错误 (1000-1099)
	CodeSuccess          = 1000 // 成功
	CodeBadRequest       = 1001 // 请求参数错误
	CodeUnauthorized     = 1002 // 未授权
	CodeForbidden        = 1003 // 禁止访问
	CodeNotFound         = 1004 // 资源不存在
	CodeMethodNotAllowed = 1005 // 方法不允许
	CodeServerError      = 1006 // 服务器错误

	// 用户相关错误 (1100-1199)
	CodeUserNotFound       = 1100 // 用户不存在
	CodeInvalidCredentials = 1101 // 无效的凭证
	CodeDuplicateUsername  = 1102 // 用户名已存在

	// 角色相关错误 (1200-1299)
	CodeRoleNotFound = 1200 // 角色不存在

	// 权限相关错误 (1300-1399)
	CodePermissionNotFound = 1300 // 权限不存在
)

// 响应结构体
type Response struct {
	Status  string      `json:"status"`  // 状态: success 或 error
	Code    int         `json:"code"`    // 状态码: 成功为1000，错误时使用对应错误码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据，可以为 nil
}

// SuccessResponse 返回成功响应，HTTP状态码统一为200
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  StatusSuccess,
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 返回错误响应，HTTP状态码统一为200
func ErrorResponse(c *fiber.Ctx, errCode int, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  StatusError,
		Code:    errCode,
		Message: message,
		Data:    nil,
	})
}

// BadRequestError 返回请求参数错误
func BadRequestError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, CodeBadRequest, message)
}

// UnauthorizedError 返回未授权错误
func UnauthorizedError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, CodeUnauthorized, message)
}

// ForbiddenError 返回禁止访问错误
func ForbiddenError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, CodeForbidden, message)
}

// NotFoundError 返回资源不存在错误
func NotFoundError(c *fiber.Ctx, message string, code ...int) error {
	errCode := CodeNotFound
	if len(code) > 0 {
		errCode = code[0]
	}
	return ErrorResponse(c, errCode, message)
}

// ServerError 返回服务器错误
func ServerError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, CodeServerError, message)
}
