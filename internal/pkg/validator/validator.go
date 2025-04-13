package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
)

// 全局验证器实例
var validate *validator.Validate

// Init 初始化验证器
func Init() {
	validate = validator.New()

	// 注册自定义验证标签
	// 可以在这里添加自定义验证规则

	// 使用JSON标签作为错误信息的字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate 验证结构体
func Validate(s interface{}) error {
	if validate == nil {
		Init()
	}
	return validate.Struct(s)
}

// ValidateRequest 验证请求并返回错误信息
func ValidateRequest(c *fiber.Ctx, req interface{}) error {
	// 解析请求体
	if err := c.BodyParser(req); err != nil {
		return response.ParamError(c, "请求参数解析失败")
	}

	// 验证请求参数
	if err := Validate(req); err != nil {
		// 处理验证错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				errors = append(errors, formatValidationError(e))
			}
			return response.ParamError(c, strings.Join(errors, "; "))
		}
		return response.ParamError(c, "请求参数验证失败")
	}

	return nil
}

// 格式化验证错误信息
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s不能为空", err.Field())
	case "email":
		return fmt.Sprintf("%s格式不正确", err.Field())
	case "min":
		return fmt.Sprintf("%s长度不能小于%s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s长度不能大于%s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s验证失败: %s", err.Field(), err.Tag())
	}
}
