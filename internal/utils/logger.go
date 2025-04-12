package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger 全局日志对象
	Logger *zap.Logger
	// Sugar 全局Sugar日志对象，提供更便捷的API
	Sugar *zap.SugaredLogger
	once  sync.Once
)

// InitLogger 初始化日志系统
func InitLogger(logLevel string, logPath string, isDevelopment bool) {
	once.Do(func() {
		// 创建日志目录
		if logPath != "" {
			if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
				fmt.Printf("创建日志目录失败: %v\n", err)
				os.Exit(1)
			}
		}

		// 设置日志级别
		var level zapcore.Level
		switch logLevel {
		case "debug":
			level = zapcore.DebugLevel
		case "info":
			level = zapcore.InfoLevel
		case "warn":
			level = zapcore.WarnLevel
		case "error":
			level = zapcore.ErrorLevel
		default:
			level = zapcore.InfoLevel
		}

		// 配置日志输出
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 开发环境使用彩色输出
		if isDevelopment {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// 设置日志写入目标
		var core zapcore.Core
		if logPath == "" || isDevelopment {
			// 开发环境或未指定日志路径，输出到控制台
			consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
			consoleCore := zapcore.NewCore(
				consoleEncoder,
				zapcore.AddSync(os.Stdout),
				level,
			)
			core = consoleCore
		} else {
			// 生产环境，同时输出到文件和控制台
			jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

			// 日志文件输出
			logFile, err := os.OpenFile(
				logPath,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0644,
			)
			if err != nil {
				fmt.Printf("打开日志文件失败: %v\n", err)
				os.Exit(1)
			}

			fileCore := zapcore.NewCore(
				jsonEncoder,
				zapcore.AddSync(logFile),
				level,
			)

			// 控制台输出
			consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
			consoleCore := zapcore.NewCore(
				consoleEncoder,
				zapcore.AddSync(os.Stdout),
				level,
			)

			// 使用Tee将日志输出到多个目标
			core = zapcore.NewTee(fileCore, consoleCore)
		}

		// 创建日志对象
		Logger = zap.New(core,
			zap.AddCaller(),                       // 添加调用者信息
			zap.AddCallerSkip(1),                  // 跳过包装函数
			zap.AddStacktrace(zapcore.ErrorLevel), // 错误日志记录堆栈
		)

		Sugar = Logger.Sugar()

		// 记录启动日志
		Sugar.Infof("日志系统初始化完成，级别: %s, 环境: %s",
			logLevel,
			map[bool]string{true: "development", false: "production"}[isDevelopment],
		)
	})
}

// 日志轮转功能
func RotateLogFile(logPath string) {
	if logPath == "" {
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return
	}

	// 生成新文件名，包含日期
	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(logPath)
	baseName := logPath[0 : len(logPath)-len(ext)]
	newPath := fmt.Sprintf("%s.%s%s", baseName, timestamp, ext)

	// 重命名当前日志文件
	if err := os.Rename(logPath, newPath); err != nil {
		fmt.Printf("日志轮转失败: %v\n", err)
	}
}

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// Debugf 记录格式化的调试级别日志
func Debugf(format string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Debugf(format, args...)
	}
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Infof 记录格式化的信息级别日志
func Infof(format string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Infof(format, args...)
	}
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Warnf 记录格式化的警告级别日志
func Warnf(format string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Warnf(format, args...)
	}
}

// Error 记录错误级别日志
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// Errorf 记录格式化的错误级别日志
func Errorf(format string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Errorf(format, args...)
	}
}

// Fatal 记录致命错误日志并结束程序
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Fatalf 记录格式化的致命错误日志并结束程序
func Fatalf(format string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Fatalf(format, args...)
	}
}

// WithContext 添加上下文信息到日志中
func WithContext(key string, value interface{}) *zap.Logger {
	if Logger == nil {
		return nil
	}
	return Logger.With(zap.Any(key, value))
}

// WithFields 添加多个字段到日志中
func WithFields(fields map[string]interface{}) *zap.Logger {
	if Logger == nil {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return Logger.With(zapFields...)
}

// Close 关闭日志，释放资源
func Close() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
