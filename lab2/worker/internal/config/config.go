package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL string
	QueueName   string
	WorkerName  string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // подгружаем .env

	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:   getEnv("QUEUE_NAME", "subtask_queue"),
		WorkerName:  getEnv("WORKER_NAME", "worker-1"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[⚠️] %s not set, using default: %s", key, fallback)
		return fallback
	}
	return val
}
