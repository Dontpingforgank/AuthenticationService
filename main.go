package main

import (
	"github.com/Dontpingforgank/AuthenticationService/Config"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"go.uber.org/zap"
)

func main() {
	appConfiguration, err := Config.LoadConfigValues("./config.json")
	if err != nil {
		panic(err)
	}

	loggerFactory := Logger.NewLoggerFactory(appConfiguration)

	appLogger, _, initLogErr := loggerFactory.NewLogger()
	if initLogErr != nil {
		panic(err)
	}

	appLogger.Log(zap.InfoLevel, "Logger initialized")

	shit, dbConnectionErr := Database.NewDbConnectionFactory(appConfiguration, loggerFactory).NewDbConnection()
	if dbConnectionErr != nil {
		panic(dbConnectionErr)
	}

	defer shit.Close()
}
