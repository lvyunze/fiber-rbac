package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/utils"
)

// IPLimiter 定义IP限制中间件的配置
type IPLimiter struct {
	// 白名单模式：只允许白名单中的IP访问，优先级高于黑名单
	WhitelistMode bool
	// IP白名单
	Whitelist []string
	// IP黑名单
	Blacklist []string
	// 允许的IP范围（CIDR格式）
	AllowedNetworks []string
}

// New 创建一个新的IP限制中间件
func NewIPLimiter(whitelistMode bool, whitelist, blacklist, allowedNetworks []string) *IPLimiter {
	return &IPLimiter{
		WhitelistMode:   whitelistMode,
		Whitelist:       whitelist,
		Blacklist:       blacklist,
		AllowedNetworks: allowedNetworks,
	}
}

// Handler 返回IP限制中间件处理函数
func (il *IPLimiter) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取客户端IP
		ip := c.IP()

		// 如果使用了代理，可能需要从X-Forwarded-For或其它头部获取真实IP
		if forwardedIP := c.Get("X-Forwarded-For"); forwardedIP != "" {
			ips := strings.Split(forwardedIP, ",")
			ip = strings.TrimSpace(ips[0]) // 取第一个IP作为客户端真实IP
		}

		// 白名单模式：只允许白名单中的IP
		if il.WhitelistMode {
			if il.isIPAllowed(ip, il.Whitelist, il.AllowedNetworks) {
				return c.Next()
			}
			return utils.ForbiddenError(c, "IP地址不在白名单中")
		}

		// 黑名单模式：阻止黑名单中的IP
		if il.isIPBlocked(ip, il.Blacklist) {
			return utils.ForbiddenError(c, "IP地址在黑名单中")
		}

		// 默认允许访问
		return c.Next()
	}
}

// isIPAllowed 检查IP是否在白名单中或匹配允许的网络
func (il *IPLimiter) isIPAllowed(ip string, whitelist []string, networks []string) bool {
	// 检查精确匹配
	for _, allowedIP := range whitelist {
		if ip == allowedIP {
			return true
		}
	}

	// 检查CIDR网络匹配
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, network := range networks {
		_, ipNet, err := net.ParseCIDR(network)
		if err != nil {
			continue
		}
		if ipNet.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// isIPBlocked 检查IP是否在黑名单中
func (il *IPLimiter) isIPBlocked(ip string, blacklist []string) bool {
	for _, blockedIP := range blacklist {
		if ip == blockedIP {
			return true
		}
	}
	return false
}
