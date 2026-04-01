package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// HashPassword 对密码进行 SHA256 哈希
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(password, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}

// GenerateToken 生成简单的 token（后续可以改为 JWT）
func GenerateToken(userID string) string {
	// 这里使用简单的方式生成 token
	// 实际生产环境应该使用 JWT
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d", userID, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
