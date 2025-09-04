package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type dbConfig struct {
	Host     string `env:"DB_HOST,required" envDefault:"localhost"`
	Port     int64  `env:"DB_PORT,required" envDefault:"5432"`
	DbName   string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required" envDefault:"postgres"`
	Password string `env:"DB_PWD,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type appConfig struct {
	Host string `env:"HOST,required" envDefault:""`
	Port string `env:"PORT,required" envDefault:"8000"`
}

type AllConfig struct {
	App appConfig
	DB  dbConfig
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
