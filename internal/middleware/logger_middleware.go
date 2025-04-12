package middleware

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"go.uber.org/zap"
)

// LoggerConfig 定义日志中间件的配置
type LoggerConfig struct {
	// 是否记录请求体（可能会影响性能）
	LogRequestBody bool
	// 是否记录响应体（可能会影响性能）
	LogResponseBody bool
	// 跳过日志记录的路径
	SkipPaths []string
}

// DefaultLoggerConfig 默认日志中间件配置
var DefaultLoggerConfig = LoggerConfig{
	LogRequestBody:  false,
	LogResponseBody: false,
	SkipPaths:       []string{},
}

// Logger 创建一个日志中间件
func Logger(config ...LoggerConfig) fiber.Handler {
	// 使用默认配置或提供的配置
	cfg := DefaultLoggerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		// 检查是否在跳过路径列表中
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		// 开始计时
		start := time.Now()

		// 记录请求信息
		method := c.Method()
		ip := c.IP()
		userAgent := c.Get("User-Agent")
		referer := c.Get("Referer")
		contentType := c.Get("Content-Type")

		// 如果配置了记录请求体，则获取并记录
		var requestBody string
		if cfg.LogRequestBody && c.Body() != nil && len(c.Body()) > 0 {
			requestBody = string(c.Body())
		}

		// 存储原始响应写入器
		responseBody := new(bytes.Buffer)
		if cfg.LogResponseBody {
			// 使用一个自定义的响应写入器来捕获响应体
			c.Response().SetBodyStream(responseBody, -1)
		}

		// 处理请求
		err := c.Next()

		// 计算处理时间
		duration := time.Since(start).Milliseconds()

		// 获取响应状态码
		statusCode := c.Response().StatusCode()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.Int("status", statusCode),
			zap.Int64("duration_ms", duration),
			zap.String("user_agent", userAgent),
		}

		// 添加可选字段
		if referer != "" {
			fields = append(fields, zap.String("referer", referer))
		}

		if contentType != "" {
			fields = append(fields, zap.String("content_type", contentType))
		}

		// 如果有请求体且配置为记录，则添加到日志字段
		if requestBody != "" && cfg.LogRequestBody {
			// 截断过长的请求体
			if len(requestBody) > 1000 {
				requestBody = requestBody[:1000] + "... (truncated)"
			}
			fields = append(fields, zap.String("request_body", requestBody))
		}

		// 如果配置为记录响应体，则获取并添加到日志字段
		if cfg.LogResponseBody {
			respBody := responseBody.String()
			if len(respBody) > 1000 {
				respBody = respBody[:1000] + "... (truncated)"
			}
			fields = append(fields, zap.String("response_body", respBody))
		}

		// 根据状态码选择适当的日志级别
		if statusCode >= 500 {
			utils.Error(fmt.Sprintf("HTTP %d %s %s", statusCode, method, path), fields...)
		} else if statusCode >= 400 {
			utils.Warn(fmt.Sprintf("HTTP %d %s %s", statusCode, method, path), fields...)
		} else {
			utils.Info(fmt.Sprintf("HTTP %d %s %s", statusCode, method, path), fields...)
		}

		return err
	}
}
