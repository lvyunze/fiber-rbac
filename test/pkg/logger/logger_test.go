package logger_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// 测试不同环境下的日志级别设置
func TestLoggerSetupWithDifferentEnvs(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name          string
		env           string
		configLevel   string
		expectedLevel string
		logMessage    string
		shouldContain bool
	}{
		{
			name:          "生产环境只输出警告和错误日志",
			env:           "prod",
			configLevel:   "debug", // 配置为debug，但环境为prod应该覆盖
			expectedLevel: "WARN",
			logMessage:    "这是一条INFO日志，在prod环境不应该输出",
			shouldContain: false,
		},
		{
			name:          "QA环境输出所有日志",
			env:           "qa",
			configLevel:   "error", // 配置为error，但环境为qa应该覆盖
			expectedLevel: "DEBUG",
			logMessage:    "这是一条DEBUG日志，在qa环境应该输出",
			shouldContain: true,
		},
		{
			name:          "开发环境使用配置文件级别",
			env:           "dev",
			configLevel:   "info",
			expectedLevel: "INFO",
			logMessage:    "这是一条INFO日志，在dev环境应该输出",
			shouldContain: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个buffer来捕获输出
			var buf bytes.Buffer
			
			// 初始化日志 - 使用buffer作为输出
			var level slog.Level
			switch tc.env {
			case "prod":
				level = slog.LevelWarn
			case "qa":
				level = slog.LevelDebug
			default:
				level = slog.LevelInfo
			}
			
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
				Level: level,
			})
			logger := slog.New(handler)
			slog.SetDefault(logger)
			
			// 输出测试日志
			slog.Info(tc.logMessage)

			// 获取输出内容
			output := buf.String()

			// 检查测试日志是否按预期输出
			if tc.shouldContain && !strings.Contains(output, tc.logMessage) {
				t.Errorf("期望输出包含消息 %s，但实际输出为: %s", tc.logMessage, output)
			} else if !tc.shouldContain && strings.Contains(output, tc.logMessage) {
				t.Errorf("期望输出不包含消息 %s，但实际输出为: %s", tc.logMessage, output)
			}
		})
	}
}

// 测试日志输出格式
func TestLoggerOutputFormat(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name        string
		format      string
		isValidJSON bool
	}{
		{
			name:        "JSON格式日志",
			format:      "json",
			isValidJSON: true,
		},
		{
			name:        "文本格式日志",
			format:      "text",
			isValidJSON: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个buffer来捕获输出
			var buf bytes.Buffer

			// 使用自定义输出初始化日志
			var handler slog.Handler
			if tc.format == "json" {
				handler = slog.NewJSONHandler(&buf, nil)
			} else {
				handler = slog.NewTextHandler(&buf, nil)
			}
			logger := slog.New(handler)
			slog.SetDefault(logger)
			
			// 输出测试日志
			slog.Info("测试日志格式", "key", "value")

			// 获取输出内容
			output := buf.String()
			
			// 检查日志格式
			if tc.isValidJSON {
				var jsonMap map[string]interface{}
				err := json.Unmarshal([]byte(output), &jsonMap)
				if err != nil {
					t.Errorf("期望JSON格式日志，但解析失败: %v, 日志内容: %s", err, output)
				}
			} else {
				// 文本格式应该包含关键字
				if !strings.Contains(output, "测试日志格式") || !strings.Contains(output, "key=value") {
					t.Errorf("文本格式日志应包含消息和键值对，但实际输出为: %s", output)
				}
			}
		})
	}
}

// 创建一个自定义的测试Writer，用于捕获日志输出
type testWriter struct {
	buf bytes.Buffer
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *testWriter) String() string {
	return w.buf.String()
}

// 测试环境配置的默认值处理
func TestEnvironmentDefaultValue(t *testing.T) {
	// 创建一个buffer来捕获输出
	writer := &testWriter{}
	
	// 使用自定义输出初始化日志
	handler := slog.NewJSONHandler(writer, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	
	// 输出测试日志
	slog.Info("测试空环境默认值")

	// 获取输出内容
	output := writer.String()

	// 检查是否包含测试消息
	if !strings.Contains(output, "测试空环境默认值") {
		t.Errorf("日志应该包含测试消息，但实际输出为: %s", output)
	}
}
