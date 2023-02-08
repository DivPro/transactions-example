package config

import "fmt"

type Config struct {
	DB    DBConfig    `yaml:"db"`
	Kafka KafkaConfig `yaml:"kafka"`
}

type DBConfig struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Addr,
		c.Database,
	)
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
