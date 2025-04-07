package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/stretchr/testify/assert"
)

func setupIPLimiterTest(t *testing.T, whitelistMode bool, whitelist, blacklist, allowedNetworks []string) *fiber.App {
	app := fiber.New()

	// 创建IP限制中间件
	ipLimiter := middleware.NewIPLimiter(whitelistMode, whitelist, blacklist, allowedNetworks)

	// 应用中间件
	app.Use(ipLimiter.Handler())

	// 添加测试路由
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("访问成功")
	})

	return app
}

// 测试白名单模式
func TestIPLimiter_WhitelistMode(t *testing.T) {
	// 设置白名单
	whitelist := []string{"192.168.1.1", "10.0.0.1"}
	allowedNetworks := []string{"192.168.2.0/24"}

	app := setupIPLimiterTest(t, true, whitelist, nil, allowedNetworks)

	// 测试白名单中的IP - 应该允许访问
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 测试白名单网络中的IP - 应该允许访问
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.2.100")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 测试不在白名单中的IP - 应该禁止访问
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.3.1")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode) // HTTP状态码统一为200

	// 验证响应内容
	var response utils.Response
	err = utils.ParseResponse(resp, &response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeForbidden, response.Code)
	assert.Equal(t, "IP地址不在白名单中", response.Message)
}

// 测试黑名单模式
func TestIPLimiter_BlacklistMode(t *testing.T) {
	// 设置黑名单
	blacklist := []string{"192.168.1.2", "10.0.0.2"}

	app := setupIPLimiterTest(t, false, nil, blacklist, nil)

	// 测试黑名单中的IP - 应该禁止访问
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.2")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode) // HTTP状态码统一为200

	// 验证响应内容
	var response utils.Response
	err = utils.ParseResponse(resp, &response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeForbidden, response.Code)
	assert.Equal(t, "IP地址在黑名单中", response.Message)

	// 测试不在黑名单中的IP - 应该允许访问
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.3")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 验证正常响应内容
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	assert.Contains(t, string(body[:n]), "访问成功")
}

// 测试代理IP处理
func TestIPLimiter_ProxyIP(t *testing.T) {
	// 设置黑名单
	blacklist := []string{"192.168.1.4"}

	app := setupIPLimiterTest(t, false, nil, blacklist, nil)

	// 测试X-Forwarded-For包含多个IP的情况 - 第一个IP在黑名单中
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.4, 10.0.0.1, 172.16.0.1")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode) // HTTP状态码统一为200

	// 验证响应内容
	var response utils.Response
	err = utils.ParseResponse(resp, &response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeForbidden, response.Code)

	// 测试X-Forwarded-For包含多个IP的情况 - 第一个IP不在黑名单中
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.5, 10.0.0.1, 172.16.0.1")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// 测试CIDR网络匹配
func TestIPLimiter_CIDRNetworks(t *testing.T) {
	// 设置白名单和CIDR网络
	whitelist := []string{"192.168.1.1"}
	allowedNetworks := []string{
		"10.0.0.0/8",     // 10.0.0.0 - 10.255.255.255
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.5.0/24", // 192.168.5.0 - 192.168.5.255
	}

	app := setupIPLimiterTest(t, true, whitelist, nil, allowedNetworks)

	// 测试场景1: 10.x.x.x网段的IP - 应该允许访问
	testCases := []struct {
		ip       string
		allowed  bool
		testName string
	}{
		{"10.0.0.1", true, "10.0.0.1 在允许的CIDR范围内"},
		{"10.10.10.10", true, "10.10.10.10 在允许的CIDR范围内"},
		{"10.255.255.255", true, "10.255.255.255 在允许的CIDR范围内"},
		{"172.16.0.1", true, "172.16.0.1 在允许的CIDR范围内"},
		{"172.20.30.40", true, "172.20.30.40 在允许的CIDR范围内"},
		{"172.31.255.255", true, "172.31.255.255 在允许的CIDR范围内"},
		{"192.168.5.100", true, "192.168.5.100 在允许的CIDR范围内"},
		{"192.168.6.1", false, "192.168.6.1 不在允许的CIDR范围内"},
		{"192.168.1.1", true, "192.168.1.1 在白名单中"},
		{"192.168.1.2", false, "192.168.1.2 不在白名单/CIDR范围内"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Forwarded-For", tc.ip)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			if tc.allowed {
				// 如果应该允许访问
				assert.Equal(t, fiber.StatusOK, resp.StatusCode)

				// 验证正常响应内容
				body := make([]byte, 1024)
				n, _ := resp.Body.Read(body)
				assert.Contains(t, string(body[:n]), "访问成功")
			} else {
				// 如果应该禁止访问
				assert.Equal(t, fiber.StatusOK, resp.StatusCode) // HTTP状态码统一为200

				// 验证响应内容
				var response utils.Response
				err = utils.ParseResponse(resp, &response)
				assert.NoError(t, err)
				assert.Equal(t, utils.StatusError, response.Status)
				assert.Equal(t, utils.CodeForbidden, response.Code)
				assert.Equal(t, "IP地址不在白名单中", response.Message)
			}
		})
	}
}

// 测试无效IP处理
func TestIPLimiter_InvalidIP(t *testing.T) {
	// 设置白名单和CIDR网络
	whitelist := []string{"192.168.1.1"}
	allowedNetworks := []string{"10.0.0.0/8"}

	app := setupIPLimiterTest(t, true, whitelist, nil, allowedNetworks)

	// 测试无效的IP
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "invalid-ip")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode) // HTTP状态码统一为200

	// 验证响应内容
	var response utils.Response
	err = utils.ParseResponse(resp, &response)
	assert.NoError(t, err)
	assert.Equal(t, utils.StatusError, response.Status)
	assert.Equal(t, utils.CodeForbidden, response.Code)
	assert.Equal(t, "IP地址不在白名单中", response.Message)

	// 测试无效的CIDR
	ipLimiter := middleware.NewIPLimiter(true, whitelist, nil, []string{"10.0.0.0/888"})
	app = fiber.New()
	app.Use(ipLimiter.Handler())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("访问成功")
	})

	// 在允许的IP范围内的IP仍然应该允许访问
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
