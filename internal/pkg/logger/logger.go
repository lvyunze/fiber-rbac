package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lvyunze/fiber-rbac/config"
)

// Setup 初始化日志配置
func Setup(cfg *config.LogConfig, env string) *slog.Logger {
	// 设置日志级别
	var level slog.Level

	// 根据环境设置日志级别
	switch env {
	case "prod":
		// 生产环境只输出警告和错误日志
		level = slog.LevelWarn
	case "uat":
		// UAT环境只输出警告和错误日志
		level = slog.LevelWarn
	case "qa":
		// QA环境输出所有日志
		level = slog.LevelDebug
	default:
		// 开发环境根据配置文件设置
		switch strings.ToLower(cfg.Level) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	}

	// 设置日志输出位置
	var output io.Writer = os.Stdout
	if cfg.Output == "file" && cfg.File != "" {
		// 确保日志目录存在
		dir := cfg.File[:strings.LastIndex(cfg.File, "/")]
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("创建日志目录失败", "error", err)
			// 失败时回退到标准输出
			output = os.Stdout
		} else {
			file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				slog.Error("打开日志文件失败", "error", err)
				// 失败时回退到标准输出
				output = os.Stdout
			} else {
				output = file
			}
		}
	}

	// 设置日志格式
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	// 创建并设置全局日志记录器
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// 输出日志设置信息
	logger.Info("日志初始化完成",
		"env", env,
		"level", level.String(),
		"format", cfg.Format,
		"output", cfg.Output)

	return logger
}
