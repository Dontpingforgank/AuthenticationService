package Controllers

import (
	"database/sql"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type registerController struct {
	config        *Models.Config
	loggerFactory Logger.LoggerFactory
	dbFactory     Database.DatabaseFactory
}

func NewRegisterController(config *Models.Config, loggerFactory Logger.LoggerFactory, dbFactory Database.DatabaseFactory) Controller {
	return &registerController{
		config:        config,
		loggerFactory: loggerFactory,
		dbFactory:     dbFactory,
	}
}

func (ctr registerController) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// parse user registration info
		userRegisterInfo, userInfoParseError := getRegistrationInfo(ctx)

		if userInfoParseError != nil {
			returnJsonError(ctx, "error parsing user registration info")
		}

		connection := establishDbConnection(ctr, ctx)

		defer func(connection *sql.DB) {
			_ = connection.Close()
		}(connection)

		insertUserInDb(userRegisterInfo, connection)
	}
}

func (ctr registerController) Route() string {
	return "/Register"
}

func (ctr registerController) Method() string {
	return http.MethodPost
}

func getRegistrationInfo(ctx *gin.Context) (*Models.UserRegisterModel, error) {
	var body Models.UserRegisterModel

	if err := ctx.ShouldBindJSON(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func returnJsonError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"Error": message,
	})
}

func establishDbConnection(ctr registerController, ctx *gin.Context) *sql.DB {
	connection, err := ctr.dbFactory.NewDbConnection()
	if err != nil {
		returnJsonError(ctx, "Couldn't establish db connection")
		return nil
	}

	err = connection.Ping()
	if err != nil {
		returnJsonError(ctx, "Db is not in reach")
		return nil
	}

	return connection
}

func insertUserInDb(userRegisterInfo *Models.UserRegisterModel, connection *sql.DB) (bool, error) {
	taken, err := checkIfEmailIsTaken(userRegisterInfo.Email, connection)
	if err != nil {
		return false, err
	}

	// Need to handle the outcome of checkIfEmailIsTakenFunction

	return taken, nil
}

func checkIfEmailIsTaken(email string, connection *sql.DB) (bool, error) {
	querry := fmt.Sprintf("select id from user_table where email = '%s'", email)

	var count int

	err := connection.QueryRow(querry).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	} else {
		return true, nil
	}

	return false, nil
}
