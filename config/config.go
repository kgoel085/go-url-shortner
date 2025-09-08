package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type otpConfig struct {
	ExpiryMinutes int64 `env:"OTP_EXPIRY_MINUTES" envDefault:"5"`
}

type dbConfig struct {
	Host     string `env:"DB_HOST,required" envDefault:"localhost"`
	Port     int64  `env:"DB_PORT,required" envDefault:"5432"`
	DbName   string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required" envDefault:"postgres"`
	Password string `env:"DB_PWD,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type smtpConfig struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     string `env:"SMTP_PORT,required"`
	Username string `env:"SMTP_USER,required" binding:"mail"`
	Password string `env:"SMTP_PWD,required"`
}

type appConfig struct {
	Name           string `env:"APP_NAME" envDefault:"URL Shortner - Go"`
	Host           string `env:"HOST,required" envDefault:""`
	Port           string `env:"PORT,required" envDefault:"8000"`
	TrustedProxies string `env:"TRUSTED_ORIGINS" envDefault:""`
}

type JWTConfig struct {
	SecretKey     string `env:"JWT_SECRET,required"`
	ExpiryMinutes int64  `env:"JWT_EXPIRY_MINUTES" envDefault:"90"`
}

type redisConfig struct {
	Addr     string `env:"REDIS_ADDR,required" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PWD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type AllConfig struct {
	APP   appConfig
	DB    dbConfig
	SMTP  smtpConfig
	OTP   otpConfig
	JWT   JWTConfig
	REDIS redisConfig
}

var Config AllConfig

func LoadConfig() {
	// Load .env if present
	_ = godotenv.Load()

	// Parse env into struct
	if err := env.Parse(&Config); err != nil {
		log.Fatalf("❌ Failed to load env: %v", err)
		panic(err)
	}
}
