package Config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type DatabaseConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	DBName     string `json:"dbName"`
	DBPassword string `json:"dbPassword"`
	SSLMode    string `json:"sslMode"`
	User       string `json:"user"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Zap      ZapConfig      `json:"zap"`
	Port     string         `json:"port"`
}

type ZapConfig struct {
	Level zapcore.Level `json:"level"`
	Path  string        `json:"path"`
}

// just loads configValues into Config struct
func LoadConfigValues(file string) (*Config, error) {
	configFile, err := os.Open(file)

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			return
		}
	}(configFile)

	if err != nil {
		return nil, err
	}

	var config Config

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
