package config

import (
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/darkseear/gophkeeper/server/internal/logger"
	"go.uber.org/zap"
)

// Config структура конфигурации приложения.
type Config struct {
	Address       string `env:"SERVER_ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	SecretKey     string `env:"SECRET_KEY"`
	ConfigFile    string `env:"CONFIG"`
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

// ConfigFile структура для хранения конфигурации из файла.
// Используется для загрузки параметров из JSON-файла конфигурации.
type ConfigFile struct {
	Address       string `json:"address"`        // -a /SERVER_ADDRESS
	DatabaseDSN   string `json:"database_dsn"`   // -d /DATABASE_DSN
	TrustedSubnet string `json:"trusted_subnet"` // -t /TRUSTED_SUBNET
}

var (
	once              sync.Once
	flagAddress       string
	flagLogLevel      string
	flagDSN           string
	flagSecretKey     string
	flagConfigFile    string
	flagTrustedSubnet string
)

// registerFlags инициализирует флаги один раз.
func registerFlags() {
	once.Do(func() {
		flag.StringVar(&flagAddress, "a", "localhost:50051", "Server address")
		flag.StringVar(&flagLogLevel, "l", "info", "Log level")
		flag.StringVar(&flagDSN, "d", "host=localhost user=postgres password=1234567890 dbname=gophkeeper sslmode=disable", "Database DSN")
		flag.StringVar(&flagSecretKey, "sk", "secretkey", "Secret key for JWT")
		flag.StringVar(&flagConfigFile, "c", "", "Path to config file")
		flag.StringVar(&flagConfigFile, "config", "", "Path to config file")
		flag.StringVar(&flagTrustedSubnet, "t", "", "Trusted subnet for internal requests")
	})
}

// New - создаёт конфиг с учётом флагов, переменных окружения и значений по умолчанию.
func New() *Config {
	registerFlags()

	if !flag.Parsed() {
		flag.Parse()
	}

	cfg := &Config{
		Address:       flagAddress,
		LogLevel:      flagLogLevel,
		DatabaseDSN:   flagDSN,
		SecretKey:     flagSecretKey,
		TrustedSubnet: flagTrustedSubnet,
		ConfigFile:    flagConfigFile,
	}

	// Переопределение значений переменными окружения
	setFromEnv(cfg)

	return cfg
}

// setFromEnv - обновляет конфиг значениями из переменных окружения.
func setFromEnv(cfg *Config) {
	configFile := getConfigFile(cfg)

	setStringFields(cfg, configFile)
}

// getConfigFile - конфиг из файла.
func getConfigFile(cfg *Config) ConfigFile {
	configFile, err := cfg.configFormFile()
	if err != nil {
		logger.Log.Error("Error reading config file", zap.Error(err))
	}
	return configFile
}

// setStringFields - строки файла.
func setStringFields(cfg *Config, configFile ConfigFile) {
	envVars := map[string]*string{
		"SERVER_ADDRESS": &cfg.Address,
		"LOG_LEVEL":      &cfg.LogLevel,
		"DATABASE_DSN":   &cfg.DatabaseDSN,
		"SECRET_KEY":     &cfg.SecretKey,
		"CONFIG":         &cfg.ConfigFile,
		"TRUSTED_SUBNET": &cfg.TrustedSubnet,
	}

	for env, ptr := range envVars {
		if val, ok := os.LookupEnv(env); ok {
			*ptr = val
		} else if *ptr == "" {
			switch env {
			case "SERVER_ADDRESS":
				if configFile.Address != "" {
					*ptr = configFile.Address
				}
			case "DATABASE_DSN":
				if configFile.DatabaseDSN != "" {
					*ptr = configFile.DatabaseDSN
				}
			case "TRUSTED_SUBNET":
				if configFile.TrustedSubnet != "" {
					*ptr = configFile.TrustedSubnet
				}
			}
		}
	}
}

// configFormFile читает конфигурацию из файла, если указан путь к файлу.
// Если файл не указан, возвращает пустую структуру ConfigFile.
// Если файл указан, но не может быть прочитан или распарсен, возвращает ошибку.

func (c *Config) configFormFile() (ConfigFile, error) {
	var configFile ConfigFile
	var filePath string
	if c.ConfigFile != "" {
		filePath = c.ConfigFile
	} else if os.Getenv("CONFIG") != "" {
		filePath = os.Getenv("CONFIG")
	} else {
		filePath = ""
	}

	if filePath == "" {
		return configFile, nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return configFile, err
	}
	if err := json.Unmarshal(file, &configFile); err != nil {
		return configFile, err
	}

	return configFile, nil
}

// DatabaseDSN: "host=localhost user=postgres password=1234567890 dbname=shorten sslmode=disable"
