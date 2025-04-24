package config

import (
	"os"
	"strconv"
)

type Config struct {
	host                string
	port                int64
	MongoDbUrl          string
	WorkerResponseQueue string
}

func NewConfig() *Config {
	portStr := os.Getenv("APP_PORT")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		port = 9091 // fallback default
	}

	return &Config{
		host:                getEnv("APP_HOST", "localhost"),
		port:                port,
		MongoDbUrl:          getEnv("MONGO_DB_URL", "mongodb://root:example@mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0&authSource=admin"),
		WorkerResponseQueue: getEnv("WORKER_RESPONSE_QUEUE", "answers_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() int64 {
	return c.port
}

func (c *Config) GetMongoDbUrl() string {
	return c.MongoDbUrl
}

func (c *Config) GetWorkerResponseQueue() string {
	return c.WorkerResponseQueue
}
