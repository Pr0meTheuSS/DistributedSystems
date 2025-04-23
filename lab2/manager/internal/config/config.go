package config

type Config struct {
	host                string
	port                int64
	WorkerResponseQueue string
}

func NewConfig() *Config {
	// TODO: parse from .env
	return &Config{
		host:                "localhost",
		port:                9091,
		WorkerResponseQueue: "answers_queue",
	}
}

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() int64 {
	return c.port
}
