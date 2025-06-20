package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HttpPort string
	DB       DBConfig
	AWS      AWSConfig
	Server   ServerConfig
	CORS     CORSConfig
	Telegram TelegramConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

func LoadConfig() *Config {
	return &Config{
		HttpPort: getEnv("HTTP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "414295mini"),
			DBName:   getEnv("DB_NAME", "silkroad"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "eu-north-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("AWS_BUCKET_NAME", "gosilkroadbucket"),
		},
		Server: ServerConfig{
			ReadTimeout:    parseDuration(getEnv("SERVER_READ_TIMEOUT", "10s")),
			WriteTimeout:   parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "10s")),
			MaxHeaderBytes: parseInt(getEnv("SERVER_MAX_HEADER_BYTES", "1048576")),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				getEnv("CORS_ALLOWED_ORIGINS", "*"),
			},
			AllowedMethods: []string{
				"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
			},
			AllowedHeaders: []string{
				"Origin", "Content-Length", "Content-Type", "Authorization", "Accept",
			},
			AllowCredentials: parseBool(getEnv("CORS_ALLOW_CREDENTIALS", "true")),
		},
		Telegram: TelegramConfig{
			BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
			ChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}

func parseInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return result
}

func parseBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return result
}
