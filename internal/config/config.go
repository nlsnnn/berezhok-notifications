package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-default:"local"`
	Db       `yaml:"db"`
	RabbitMQ `yaml:"rabbitmq"`
	Email    `yaml:"email"`
}

type Db struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"8080"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"notification"`
}

type RabbitMQ struct {
	URL string `yaml:"url" env:"RABBITMQ_URL"`
}

type Email struct {
	From string `yaml:"from" env:"EMAIL_FROM"`

	ResendApiKey string `yaml:"api_key" env:"RESEND_API_KEY"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = ".env"
		log.Println("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
