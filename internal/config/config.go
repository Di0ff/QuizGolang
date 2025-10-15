package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ConfigQuiz struct {
	QuestionsLimit string
}

type ConfigMigrations struct {
	Path string
}

type Cfg struct {
	Database   ConfigDB
	Redis      ConfigRedis
	App        ConfigApp
	Logger     ConfigLog
	Telegram   ConfigTg
	OpenAI     ConfigAI
	Quiz       ConfigQuiz
	Migrations ConfigMigrations
}

type ConfigDB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type ConfigRedis struct {
	Addr       string
	Password   string
	DB         string
	SessionTTL string
}

type ConfigApp struct {
	Name        string
	Version     string
	Port        string
	Host        string
	StaticDir   string
	FrontendDir string
}

type ConfigLog struct {
	Env   string
	Level string
}

type ConfigTg struct {
	Token string
	URL   string
}

type ConfigAI struct {
	KeyAI     string
	Model     string
	MaxTokens string
}

func Load() (*Cfg, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Файл не найден: %v", err)
	}

	cfg := &Cfg{
		Database: ConfigDB{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Name:     getEnv("DB_NAME"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASS"),
		},

		Redis: ConfigRedis{
			Addr:       getEnv("REDIS_ADDR"),
			Password:   getEnv("REDIS_PASSWORD"),
			DB:         getEnv("REDIS_DB"),
			SessionTTL: getEnv("REDIS_SESSION_TTL"),
		},

		App: ConfigApp{
			Name:        getEnv("APP_NAME"),
			Version:     getEnv("APP_VERSION"),
			Port:        getEnv("APP_PORT"),
			Host:        getEnv("APP_HOST"),
			StaticDir:   getEnv("STATIC_DIR"),
			FrontendDir: getEnv("FRONTEND_DIR"),
		},

		Logger: ConfigLog{
			Env:   getEnv("ENV"),
			Level: getEnv("LOG_LEVEL"),
		},

		Telegram: ConfigTg{
			Token: getEnv("API_TOKEN"),
			URL:   getEnv("URL"),
		},

		OpenAI: ConfigAI{
			KeyAI:     getEnv("OPENAI_API_KEY"),
			Model:     getEnv("OPENAI_MODEL"),
			MaxTokens: getEnv("OPENAI_MAX_TOKENS"),
		},

		Quiz: ConfigQuiz{
			QuestionsLimit: getEnv("QUESTIONS_LIMIT"),
		},

		Migrations: ConfigMigrations{
			Path: getEnv("MIGRATIONS_PATH"),
		},
	}

	return cfg, nil
}

func getEnv(temp string) string {
	if value := os.Getenv(temp); value != "" {
		return value
	}

	return ""
}
