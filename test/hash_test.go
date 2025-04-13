package test

import (
	"fmt"
	"github.com/lvyunze/fiber-rbac/internal/pkg/hash"
	"testing"
)

// 生成测试用的密码哈希
func TestGeneratePasswordHash(t *testing.T) {
	password := "password123"
	hashedPassword, err := hash.GeneratePassword(password)
	if err != nil {
		t.Fatalf("生成密码哈希失败: %v", err)
	}
	
	fmt.Printf("密码 '%s' 的哈希值: %s\n", password, hashedPassword)
	
	// 验证哈希
	valid, err := hash.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("验证密码失败: %v", err)
	}
	if !valid {
		t.Fatalf("密码验证应该成功")
	}
}
