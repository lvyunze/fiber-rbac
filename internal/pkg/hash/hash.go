package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// 定义参数常量
const (
	SaltLength  = 16
	KeyLength   = 32
	TimeCost    = 1
	Memory      = 64 * 1024
	Parallelism = 4
)

// 定义错误
var (
	ErrInvalidHash         = errors.New("提供的哈希格式无效")
	ErrIncompatibleVersion = errors.New("不兼容的版本")
)

// GeneratePassword 生成密码哈希
func GeneratePassword(password string) (string, error) {
	// 生成随机盐值
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// 使用argon2id算法生成哈希
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		TimeCost,
		Memory,
		Parallelism,
		KeyLength,
	)

	// 构建哈希字符串，格式：$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		Memory,
		TimeCost,
		Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析哈希字符串
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// 使用相同参数计算哈希
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.TimeCost,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	// 安全比较哈希值（防止时序攻击）
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// 哈希参数结构
type params struct {
	Memory      uint32
	TimeCost    uint32
	Parallelism uint8
	KeyLength   uint32
}

// 解析哈希字符串
func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.TimeCost, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
