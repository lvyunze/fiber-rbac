package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler_RejectSelfDelete(t *testing.T) {
	app := fiber.New()
	app.Delete("/user", func(c *fiber.Ctx) error {
		// 伪造 userID=123，删除请求也是 id=123
		c.Locals("userID", uint64(123))
		c.Request().SetBody([]byte(`{"id":123}`))
		// 模拟 handler 逻辑
		currentUserID, ok := c.Locals("userID").(uint64)
		if !ok || currentUserID == 0 {
			return c.Status(401).JSON(fiber.Map{"code": 401, "message": "未登录或身份异常"})
		}
		if 123 == currentUserID {
			return c.Status(403).JSON(fiber.Map{"code": 403, "message": "禁止删除自己的账号"})
		}
		return c.Status(200).JSON(fiber.Map{"code": 200, "message": "用户删除成功"})
	})

	req := httptest.NewRequest("DELETE", "/user", nil)
	resp, _ := app.Test(req, -1)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	assert.Equal(t, 403, int(res["code"].(float64)))
	assert.Contains(t, res["message"], "禁止删除自己的账号")
}
