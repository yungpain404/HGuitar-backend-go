package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL          string
	DBUser         string
	DBPassword     string
	JWTSecret      []byte
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
	ServerPort     string
	AllowedOrigins []string
	FirebaseProjectID string
}

func Load() *Config {
	_ = godotenv.Load() // nạp file .env nếu có (không lỗi nếu không có file)
	cfg := &Config{
		DBURL:          getEnv("DB_URL", "sqlserver://sa:YourStrong!Passw0rd@localhost:1434?database=GUITAR_LEARNING_DB&encrypt=true&trustServerCertificate=true"),
		DBUser:         getEnv("DB_USERNAME", "sa"),
		DBPassword:     getEnv("DB_PASSWORD", "YourStrong!Passw0rd"),
		JWTSecret:      []byte(getEnv("JWT_SECRET", "replace-with-a-strong-256-bit-secret-key-before-deploying")),
		AccessExpiry:   time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MS", 900000)) * time.Millisecond,
		RefreshExpiry:  time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_MS", 604800000)) * time.Millisecond,
		ServerPort:     getEnv("SERVER_PORT", "5501"),
		AllowedOrigins: splitEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5500,http://127.0.0.1:5500"),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", "your-firebase-project-id"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func splitEnv(key, def string) []string {
	if v := os.Getenv(key); v != "" {
		parts := splitString(v, ",")
		if len(parts) > 0 {
			return parts
		}
	}
	return splitString(def, ",")
}

func splitString(s, sep string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+1] == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
