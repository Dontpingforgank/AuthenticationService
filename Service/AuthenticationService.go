package Service

import (
	"database/sql"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Config"
	"github.com/Dontpingforgank/AuthenticationService/Controllers"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Service interface {
	Run()
}

type authenticationService struct {
	configuration       *Config.Config
	loggerFactory       Logger.LoggerFactory
	dbConnectionFactory Database.DatabaseFactory
	controllers         []Controllers.Controller
}

func NewAuthenticationService(config *Config.Config, loggerFactory Logger.LoggerFactory, dbConnectionFactory Database.DatabaseFactory,
	controller ...Controllers.Controller) Service {
	return &authenticationService{
		configuration:       config,
		loggerFactory:       loggerFactory,
		dbConnectionFactory: dbConnectionFactory,
		controllers:         controller,
	}
}

func (authService authenticationService) Run() {
	logger, closeLogger, err := authService.loggerFactory.NewLogger()
	if err != nil {
		panic(err)
	}

	defer closeLogger()

	logger.Info("Service started")

	connection, err := authService.dbConnectionFactory.NewDbConnection()
	if err != nil {
		return
	}

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			return
		}
	}(connection)

	router, err := authService.buildRouter(logger)

	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf("localhost:%s", authService.configuration.Port))
	if err != nil {
		panic(err)
	}
}

func (authService authenticationService) buildRouter(appLogger *zap.Logger) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(ginzap.Ginzap(appLogger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(appLogger, true))

	for _, ctr := range authService.controllers {
		method := ctr.Method()

		switch method {
		case http.MethodGet:
			router.GET(ctr.Route(), ctr.Handle())
		case http.MethodPost:
			router.POST(ctr.Route(), ctr.Handle())
		default:
			return nil, fmt.Errorf("cannot create HTTP handler of method %s", method)
		}
	}

	return router, nil
}
