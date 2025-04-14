package middleware

import (
	"log/slog"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
)

// 测试IP的本地上下文键名
const testIPKey = "test_ip"

// WithTestIP 创建一个设置测试IP的中间件，仅用于测试
func WithTestIP(ip string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(testIPKey, ip)
		return c.Next()
	}
}

// IPWhitelist 创建IP白名单中间件
func IPWhitelist(cfg *config.SecurityConfig) fiber.Handler {
	// 如果未启用白名单，返回一个空的中间件
	if !cfg.EnableWhitelist {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// 预处理IP白名单，解析CIDR
	var ipNets []*net.IPNet
	var singleIPs []net.IP

	for _, ipStr := range cfg.IPWhitelist {
		// 检查是否是CIDR格式
		if strings.Contains(ipStr, "/") {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				slog.Error("解析CIDR失败", "cidr", ipStr, "error", err)
				continue
			}
			ipNets = append(ipNets, ipNet)
		} else {
			ip := net.ParseIP(ipStr)
			if ip == nil {
				slog.Error("解析IP失败", "ip", ipStr)
				continue
			}
			singleIPs = append(singleIPs, ip)
		}
	}

	slog.Info("IP白名单已启用", "单个IP数量", len(singleIPs), "IP网段数量", len(ipNets))

	// 返回中间件处理函数
	return func(c *fiber.Ctx) error {
		// 获取客户端IP
		var clientIP string
		
		// 优先从本地上下文获取测试IP（用于测试）
		if testIP, ok := c.Locals(testIPKey).(string); ok && testIP != "" {
			clientIP = testIP
		} else {
			// 从各种请求头中获取真实IP
			clientIP = c.Get("X-Real-IP")
			if clientIP == "" {
				clientIP = c.Get("X-Forwarded-For")
				if clientIP != "" {
					// X-Forwarded-For可能包含多个IP，取第一个
					ips := strings.Split(clientIP, ",")
					clientIP = strings.TrimSpace(ips[0])
				}
			}
			
			// 如果请求头中没有IP，使用远程地址
			if clientIP == "" {
				clientIP = c.IP()
			}
		}
		
		// 解析IP
		ip := net.ParseIP(clientIP)
		if ip == nil {
			slog.Error("无法解析客户端IP", "ip", clientIP)
			return response.Forbidden(c, "无效的IP地址")
		}

		// 检查是否在白名单中
		allowed := false

		// 检查单个IP
		for _, allowedIP := range singleIPs {
			if allowedIP.Equal(ip) {
				allowed = true
				break
			}
		}

		// 如果不在单个IP列表中，检查CIDR
		if !allowed {
			for _, ipNet := range ipNets {
				if ipNet.Contains(ip) {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			slog.Warn("IP不在白名单中，拒绝访问", "ip", clientIP)
			return response.Forbidden(c, "IP不在白名单中")
		}

		// IP在白名单中，继续处理请求
		return c.Next()
	}
}
