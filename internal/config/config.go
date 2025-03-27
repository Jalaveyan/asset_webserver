package config

import (
	"os"
)

// Этот файл читает настройки из переменных окружения (с дефолтными значениями) и используется для настройки подключения к базе данных, порта приложения и путей к TLS-сертификатам.

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	AppPort     string
	TLSCertPath string
	TLSKeyPath  string
}

func NewConfig() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "testdb"),
		AppPort:     getEnv("APP_PORT", "8443"),
		TLSCertPath: getEnv("TLS_CERT_PATH", "certs/cert.pem"), // например, cert.pem
		TLSKeyPath:  getEnv("TLS_KEY_PATH", "certs/key.pem"),   // например, key.pem
	}
}

func getEnv(key, defVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defVal
	}
	return val
}
