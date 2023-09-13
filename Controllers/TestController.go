package Controllers

import (
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type testController struct {
	config        *Models.Config
	loggerFactory Logger.LoggerFactory
	dbFactory     Database.DatabaseFactory
}

func NewTestController(config *Models.Config, loggerFactory Logger.LoggerFactory, dbFactory Database.DatabaseFactory) Controller {
	return &testController{
		config:        config,
		loggerFactory: loggerFactory,
		dbFactory:     dbFactory,
	}
}

func (ctr testController) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Message": "Helllo there im working",
		})
		return
	}
}

func (ctr testController) Route() string {
	return "/"
}

func (ctr testController) Method() string {
	return http.MethodGet
}
