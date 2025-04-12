package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"go.uber.org/zap"
)

func main() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("无法创建日志目录: %v", err)
	}

	// 初始化日志系统
	utils.InitLogger("debug", "logs/example.log", true)
	defer utils.Close()

	// 记录不同级别的日志
	utils.Debug("这是一条调试日志")
	utils.Info("这是一条信息日志")
	utils.Warn("这是一条警告日志")
	utils.Error("这是一条错误日志")

	// 使用结构化日志字段
	utils.Info("用户登录成功",
		zap.String("username", "lvyunze"),
		zap.String("ip", "192.168.1.1"),
		zap.Int("user_id", 12345),
	)

	// 记录包含嵌套结构的日志
	userInfo := map[string]interface{}{
		"id":       12345,
		"username": "lvyunze",
		"roles":    []string{"admin", "editor"},
		"metadata": map[string]string{
			"last_login": time.Now().Format(time.RFC3339),
			"device":     "Web Browser",
		},
	}
	utils.Info("用户详细信息", zap.Any("user", userInfo))

	// 创建带有上下文的日志器
	contextLogger := utils.WithContext("request_id", "abc-123-xyz")
	contextLogger.Info("处理API请求",
		zap.String("method", "GET"),
		zap.String("path", "/api/v1/users"),
	)

	// 记录带有错误的日志
	err := fmt.Errorf("数据库连接失败")
	utils.Error("操作失败",
		zap.Error(err),
		zap.String("operation", "fetch_users"),
	)

	// 设置Fiber应用示例
	app := fiber.New()

	// 简单的处理函数，记录日志
	app.Get("/", func(c *fiber.Ctx) error {
		// 记录请求信息
		utils.Info("收到请求",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
		)

		// 模拟处理请求
		time.Sleep(100 * time.Millisecond)

		// 记录响应信息
		utils.Info("请求处理完成",
			zap.Int("status", fiber.StatusOK),
			zap.Duration("duration", 100*time.Millisecond),
		)

		return c.SendString("日志示例程序 - 请查看控制台和日志文件")
	})

	// 记录启动信息
	utils.Info("示例服务器启动", zap.String("port", ":3000"))

	// 启动服务器（非阻塞）
	go func() {
		if err := app.Listen(":3000"); err != nil {
			utils.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待一段时间后结束示例程序
	time.Sleep(10 * time.Second)
	utils.Info("示例程序结束")
}
