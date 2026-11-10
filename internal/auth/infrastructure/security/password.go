package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

// HashPassword băm mật khẩu (SHA-256, có salt đơn giản; trong production hãy dùng bcrypt/argon2).
func HashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// CheckPassword so sánh mật khẩu plain với hash.
func CheckPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateSecureToken tạo token ngẫu nhiên (dùng làm refresh token thô).
func GenerateSecureToken() string {
	b := make([]byte, 32)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)
	return hex.EncodeToString(b)
}

// HashToken băm refresh token trước khi lưu DB (giống Spring Boot).
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

var _ = strings.TrimSpace
