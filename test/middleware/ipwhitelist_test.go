package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/middleware"
	"github.com/lvyunze/fiber-rbac/internal/pkg/response"
	"github.com/stretchr/testify/assert"
)

// 创建测试应用
func createTestApp(enableWhitelist bool, ipWhitelist []string, testIP string) *fiber.App {
	// 创建配置
	cfg := &config.SecurityConfig{
		EnableWhitelist: enableWhitelist,
		IPWhitelist:    ipWhitelist,
	}

	// 创建Fiber应用
	app := fiber.New()

	// 如果提供了测试IP，添加测试IP中间件
	if testIP != "" {
		app.Use(middleware.WithTestIP(testIP))
	}

	// 注册IP白名单中间件
	app.Use(middleware.IPWhitelist(cfg))

	// 添加测试路由
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	return app
}

// 测试白名单禁用时所有IP都可以访问
func TestIPWhitelist_Disabled(t *testing.T) {
	app := createTestApp(false, []string{"127.0.0.1"}, "")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证状态码
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 验证响应内容
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}

// 测试白名单启用时，白名单内的IP可以访问
func TestIPWhitelist_AllowedIP(t *testing.T) {
	app := createTestApp(true, []string{"127.0.0.1", "192.168.1.100"}, "192.168.1.100")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证状态码
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 验证响应内容
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}

// 测试白名单启用时，白名单外的IP被拒绝访问
func TestIPWhitelist_BlockedIP(t *testing.T) {
	app := createTestApp(true, []string{"127.0.0.1", "192.168.1.100"}, "192.168.1.200")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证响应状态码为200（HTTP状态码）
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 解析响应体，验证业务状态码为CodeForbidden
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var respData response.Response
	err = json.Unmarshal(body, &respData)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeForbidden, respData.Code)
	assert.Equal(t, "error", respData.Status)
}

// 测试CIDR格式的IP范围
func TestIPWhitelist_CIDR(t *testing.T) {
	app := createTestApp(true, []string{"127.0.0.1", "192.168.1.0/24"}, "192.168.1.150")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证状态码
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 验证响应内容
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}

// 测试IP不在CIDR范围内
func TestIPWhitelist_OutsideCIDR(t *testing.T) {
	app := createTestApp(true, []string{"127.0.0.1", "192.168.1.0/24"}, "192.168.2.150")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证响应状态码为200（HTTP状态码）
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 解析响应体，验证业务状态码为CodeForbidden
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var respData response.Response
	err = json.Unmarshal(body, &respData)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeForbidden, respData.Code)
	assert.Equal(t, "error", respData.Status)
}

// 测试无效IP地址的处理
func TestIPWhitelist_InvalidIP(t *testing.T) {
	app := createTestApp(true, []string{"127.0.0.1"}, "invalid-ip")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证响应状态码为200（HTTP状态码）
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 解析响应体，验证业务状态码为CodeForbidden
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var respData response.Response
	err = json.Unmarshal(body, &respData)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeForbidden, respData.Code)
	assert.Equal(t, "error", respData.Status)
}

// 测试IPv6地址在白名单内
func TestIPWhitelist_IPv6Allowed(t *testing.T) {
	app := createTestApp(true, []string{"::1", "fe80::1"}, "::1")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证状态码
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 验证响应内容
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}

// 测试IPv6地址不在白名单内
func TestIPWhitelist_IPv6Blocked(t *testing.T) {
	app := createTestApp(true, []string{"::1", "fe80::1"}, "2001:db8::1")

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	// 执行请求
	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// 验证响应状态码为200（HTTP状态码）
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// 解析响应体，验证业务状态码为CodeForbidden
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var respData response.Response
	err = json.Unmarshal(body, &respData)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeForbidden, respData.Code)
	assert.Equal(t, "error", respData.Status)
}
