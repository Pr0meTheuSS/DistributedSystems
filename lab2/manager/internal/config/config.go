package config

import (
	"os"
	"strconv"
)

type Config struct {
	host                string
	port                int64
	mongoDbUrl          string
	workerResponseQueue string
	workersAmount       int64
	rabbitMQUrl         string
	dbName              string
}

func NewConfig() *Config {
	workersAmountStr := os.Getenv("WORKERS_AMOUNT")
	workersAmount, err := strconv.ParseInt(workersAmountStr, 10, 64)
	if err != nil {
		workersAmount = 3
	}

	RabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	portStr := os.Getenv("APP_PORT")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		port = 9091
	}

	return &Config{
		host:                getEnv("APP_HOST", "localhost"),
		port:                port,
		mongoDbUrl:          getEnv("MONGO_DB_URL", "mongodb://root:example@mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0&authSource=admin"),
		workerResponseQueue: getEnv("WORKER_RESPONSE_QUEUE", "answers_queue"),
		workersAmount:       workersAmount,
		rabbitMQUrl:         RabbitMQURL,
		dbName:              getEnv("DB_NAME", "bf-service"),
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
	return c.mongoDbUrl
}

func (c *Config) GetWorkerResponseQueue() string {
	return c.workerResponseQueue
}

func (c *Config) GetWorkersAmount() int64 {
	return c.workersAmount
}

func (c *Config) GetRabbitMQUrl() string {
	return c.rabbitMQUrl
}

func (c *Config) GetDBName() string {
	return c.dbName
}
