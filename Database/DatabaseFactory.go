package Database

import (
	"database/sql"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Config"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type DatabaseFactory interface {
	NewDbConnection() (*sql.DB, error)
}

type databaseFactory struct {
	config        *Config.Config
	loggerFactory Logger.LoggerFactory
}

func NewDbConnectionFactory(conf *Config.Config, loggerFactory Logger.LoggerFactory) DatabaseFactory {
	return &databaseFactory{
		config:        conf,
		loggerFactory: loggerFactory,
	}
}

func (dbFactory databaseFactory) NewDbConnection() (*sql.DB, error) {
	logger, closeLogger, err := dbFactory.loggerFactory.NewLogger()
	if err != nil {
		return nil, err
	}

	defer closeLogger()

	logger.Info("initializing db connection")

	dbConnectionString := buildDbConnectionString(&dbFactory.config.Database)

	start := time.Now()
	dbConnection, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		logger.Error("Couldn't open dbConnection")
		return nil, err
	}

	errPing := dbConnection.Ping()
	if err != nil {
		logger.Error("Db is not responding to ping")
		return nil, errPing
	}

	elapsed := time.Since(start).Milliseconds()

	logger.Info("Connection opened, time took: % ", zap.Int64("elapsed_time_ms", elapsed))

	return dbConnection, nil
}

func buildDbConnectionString(config *Config.DatabaseConfig) string {
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s dbname =%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.DBName, config.DBPassword, config.SSLMode)

	return dbConnectionString
}
