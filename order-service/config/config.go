package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	NATS     NATSConfig
	HTTP     HTTPConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type NATSConfig struct {
	URL       string
	ClusterID string
	ClientID  string
	Subject   string
}

type HTTPConfig struct {
	Port string
}

func GetConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "orderservice"),
			Password: getEnv("DB_PASSWORD", "1234"),
			DBName:   getEnv("DB_NAME", "ordersdb"),
		},
		NATS: NATSConfig{
			URL:       getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID: getEnv("NATS_CLUSTER_ID", "test-cluster"),
			ClientID:  getEnv("NATS_CLIENT_ID", "order-service"),
			Subject:   getEnv("NATS_SUBJECT", "orders"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
	}
}

func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
