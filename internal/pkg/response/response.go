package response

import (
	"github.com/gofiber/fiber/v2"
)

// 响应状态码
const (
	CodeSuccess      = 1000 // 成功
	CodeParamError   = 1001 // 参数错误
	CodeUnauthorized = 1002 // 未授权
	CodeForbidden    = 1003 // 禁止访问
	CodeNotFound     = 1004 // 资源不存在
	CodeServerError  = 1005 // 服务器错误
)

// Response 统一响应结构
type Response struct {
	Status  string      `json:"status"`  // 状态：success 或 error
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}, msg string) error {
	if msg == "" {
		msg = "操作成功"
	}
	return c.JSON(Response{
		Status:  "success",
		Code:    CodeSuccess,
		Message: msg,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *fiber.Ctx, code int, msg string) error {
	return c.JSON(Response{
		Status:  "error",
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// ParamError 参数错误
func ParamError(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "参数错误"
	}
	return Fail(c, CodeParamError, msg)
}

// Unauthorized 未授权
func Unauthorized(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "未授权或授权已过期"
	}
	return Fail(c, CodeUnauthorized, msg)
}

// Forbidden 禁止访问
func Forbidden(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "无权限访问"
	}
	return Fail(c, CodeForbidden, msg)
}

// NotFound 资源不存在
func NotFound(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "资源不存在"
	}
	return Fail(c, CodeNotFound, msg)
}

// ServerError 服务器错误
func ServerError(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "服务器内部错误"
	}
	return Fail(c, CodeServerError, msg)
}
