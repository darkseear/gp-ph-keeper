package config

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// Сброс состояния флагов перед тестом
	defer func() {
		flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	}()

	cfg := New()
	fmt.Printf("Running  %s", cfg.Address)
	assert.Equal(t, "localhost:8080", cfg.Address)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "", cfg.DatabaseDSN)
	assert.Equal(t, "secretkey", cfg.SecretKey)
}

func TestConfigWithEnv(t *testing.T) {
	// Сброс состояния флагов и переменных окружения
	defer func() {
		flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("SECRET_KEY")
	}()

	// Установка переменных окружения для теста
	os.Setenv("SERVER_ADDRESS", "localhost:9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DATABASE_DSN", "test_dsn")
	os.Setenv("SECRET_KEY", "test_secret")

	cfg := New()

	assert.Equal(t, "localhost:9090", cfg.Address)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "test_dsn", cfg.DatabaseDSN)
	assert.Equal(t, "test_secret", cfg.SecretKey)
}
