package Controllers

import (
	"database/sql"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/Dontpingforgank/AuthenticationService/Utils/DbConnectionUtils"
	"github.com/Dontpingforgank/AuthenticationService/Utils/SessionUtils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LogInController struct {
	config        *Models.Config
	loggerFactory Logger.LoggerFactory
	dbFactory     Database.DatabaseFactory
}

func NewLogInController(config *Models.Config, loggerFactory Logger.LoggerFactory, dbFactory Database.DatabaseFactory) Controller {
	return &LogInController{
		config:        config,
		loggerFactory: loggerFactory,
		dbFactory:     dbFactory,
	}
}

func (ctr LogInController) Handle() gin.HandlerFunc {
	logger, closeLog, logErr := ctr.loggerFactory.NewLogger()
	if logErr != nil {
		panic(logErr)
	}

	defer closeLog()

	return func(context *gin.Context) {
		var credentials Models.UserLogInModel

		err := context.Bind(&credentials)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"Message": "Incorrect request body",
			})
		}

		connection, connectionError := DbConnectionUtils.EstablishDbConnection(ctr.dbFactory)

		if connectionError != nil {
			logger.Error("Couldn't establish db connection", zap.String("err", connectionError.Error()))
		}

		defer func(connection *sql.DB) {
			_ = connection.Close()
		}(connection)

		userId, emailErr := DbConnectionUtils.CheckIfEmailIsTaken(credentials.Email, connection)
		if emailErr != nil {
			logger.Error(emailErr.Error())
		}

		if userId <= 0 {
			context.JSON(http.StatusBadRequest, gin.H{
				"Success": false,
				"Message": fmt.Sprintf("user with given email is not registered. Email: %s", credentials.Email),
			})
		}

		userRealPassword, passErr := DbConnectionUtils.GetUserPassword(userId, connection)
		if passErr != nil {
			logger.Error(passErr.Error())
		}

		compareErr := bcrypt.CompareHashAndPassword([]byte(userRealPassword), []byte(credentials.Password))
		if compareErr != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"Success": false,
				"Message": "password do not match",
			})
		} else {
			jwtToken, tokenError := SessionUtils.GenerateJwtToken(ctr.config.JwtSecret, Models.UserClaimsModel{Id: userId})
			if tokenError != nil {
				context.JSON(http.StatusBadRequest, gin.H{
					"Success": false,
					"Message": "Couldn't generate token",
				})
			} else {
				context.SetCookie("Authorization", jwtToken, 3600*24*14, "", "", false, true)
				context.JSON(http.StatusOK, gin.H{
					"Success": true,
					"Message": fmt.Sprintf("Welcome back %s", credentials.Email),
					"Token":   jwtToken,
				})
			}
		}
	}
}

func (l LogInController) Route() string {
	return "/LogIn"
}

func (l LogInController) Method() string {
	return http.MethodPost
}
