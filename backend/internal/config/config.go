package config

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	AppName         string
	AppEnv          string
	HTTPAddr        string
	MySQLDSN        string
	Redis           redis.Options
	JWTSecret       string
	WeChatAppID     string
	WeChatAppSecret string
}

func Load() Config {
	return Config{
		AppName:   getEnv("APP_NAME", "retrofit"),
		AppEnv:    getEnv("APP_ENV", "local"),
		HTTPAddr:  getEnv("HTTP_ADDR", ":8080"),
		MySQLDSN:  getEnv("MYSQL_DSN", "retrofit:retrofit@tcp(127.0.0.1:3306)/retrofit?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret: getEnv("JWT_SECRET", "change-me"),
		Redis: redis.Options{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		WeChatAppID:     getEnv("WECHAT_APP_ID", ""),
		WeChatAppSecret: getEnv("WECHAT_APP_SECRET", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
