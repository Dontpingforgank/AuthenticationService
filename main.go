package main

import (
	"github.com/Dontpingforgank/AuthenticationService/Config"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Service"
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

	dbConnectionFactory := Database.NewDbConnectionFactory(appConfiguration, loggerFactory)

	service := Service.NewAuthenticationService(appConfiguration, loggerFactory, dbConnectionFactory)

	service.Run()
}
