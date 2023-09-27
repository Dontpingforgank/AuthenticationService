package Models

import "go.uber.org/zap/zapcore"

type DatabaseConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	DBName     string `json:"dbName"`
	DBPassword string `json:"dbPassword"`
	SSLMode    string `json:"sslMode"`
	User       string `json:"user"`
}

type Config struct {
	Database  DatabaseConfig `json:"database"`
	Zap       ZapConfig      `json:"zap"`
	Port      string         `json:"port"`
	JwtSecret string         `json:"jwtSecret"`
	IpAddress string         `json:"ipAddress"`
}

type ZapConfig struct {
	Level zapcore.Level `json:"level"`
	Path  string        `json:"path"`
}
