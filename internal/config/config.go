package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("ошибка загрузки файла окружения :%v", err)
	}

	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ключ не найден")
	}

	return &Config{APIKey: key}, nil
}
