package logger_test

import (
	"os"
	"testing"

	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/pkg/logger"
)

// 测试配置文件加载和日志初始化的集成测试
func TestLoggerWithConfigFile(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name          string
		configFile    string
		expectedEnv   string
		expectedLevel string
	}{
		{
			name:          "开发环境配置",
			configFile:    "test_config_dev.yaml",
			expectedEnv:   "dev",
			expectedLevel: "INFO",
		},
		{
			name:          "QA环境配置",
			configFile:    "test_config_qa.yaml",
			expectedEnv:   "qa",
			expectedLevel: "DEBUG",
		},
		{
			name:          "生产环境配置",
			configFile:    "test_config_prod.yaml",
			expectedEnv:   "prod",
			expectedLevel: "WARN",
		},
	}

	// 创建测试配置文件
	createTestConfigFile("test_config_dev.yaml", "dev")
	createTestConfigFile("test_config_qa.yaml", "qa")
	createTestConfigFile("test_config_prod.yaml", "prod")

	// 测试完成后删除测试文件
	defer func() {
		os.Remove("test_config_dev.yaml")
		os.Remove("test_config_qa.yaml")
		os.Remove("test_config_prod.yaml")
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 加载配置文件
			cfg, err := config.LoadConfig(tc.configFile)
			if err != nil {
				t.Fatalf("加载配置文件失败: %v", err)
			}

			// 初始化日志 - 这里我们不捕获输出，只是验证函数不会崩溃
			_ = logger.Setup(&cfg.Log, cfg.Env)

			// 验证环境设置
			if cfg.Env != tc.expectedEnv {
				t.Errorf("期望环境为 %s，但实际为: %s", tc.expectedEnv, cfg.Env)
			}
		})
	}
}

// 创建测试配置文件
func createTestConfigFile(filename, env string) {
	content := `# 测试配置文件

# 环境配置 (dev, qa, prod)
env: "` + env + `"

# 日志配置
log:
  level: "info" # debug, info, warn, error
  format: "json" # json, text
  output: "console" # console, file
  file: "logs/test.log" # 日志文件路径（当output为file时使用）
`
	os.WriteFile(filename, []byte(content), 0644)
}
