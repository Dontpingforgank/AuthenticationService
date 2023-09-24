package Controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/CustomErrors"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/Dontpingforgank/AuthenticationService/Utils/DbConnectionUtils"
	"github.com/Dontpingforgank/AuthenticationService/Utils/PasswordUtils"
	"github.com/Dontpingforgank/AuthenticationService/Utils/UserUtils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	logger, closeLog, logErr := ctr.loggerFactory.NewLogger()
	if logErr != nil {
		panic(logErr)
	}

	defer closeLog()

	return func(ctx *gin.Context) {
		userRegisterInfo, userInfoParseError := getRegistrationInfo(ctx)

		if userInfoParseError != nil {
			returnJsonError(ctx, "error parsing user registration info")
		}

		connection, connectionError := DbConnectionUtils.EstablishDbConnection(ctr.dbFactory)

		if connectionError != nil {
			logger.Error("Couldn't establish db connection", zap.String("err", connectionError.Error()))
		}

		defer func(connection *sql.DB) {
			_ = connection.Close()
		}(connection)

		_, err := insertUserInDb(userRegisterInfo, connection)
		switch {
		case errors.Is(err, CustomErrors.UserTakenError{}):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Message": fmt.Sprintf("%s %s", err.Error(), userRegisterInfo.Email),
			})
		case errors.Is(err, CustomErrors.InsecurePassword{}):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Message": fmt.Sprintf("%s", err.Error()),
			})
		case err != nil:
			logger.Log(zap.ErrorLevel, "couldn't register user", zap.String("err", err.Error()))
		}
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
		"Success": false,
		"Error":   message,
	})
}

func insertUserInDb(userRegisterInfo *Models.UserRegisterModel, connection *sql.DB) (bool, error) {
	verified, verifiedErr := UserUtils.VerifyUserData(userRegisterInfo)
	if verifiedErr != nil {
		return false, verifiedErr
	}

	if verified {
		id, err := DbConnectionUtils.CheckIfEmailIsTaken(userRegisterInfo.Email, connection)

		if err != nil {
			return false, err
		}

		if id > 0 {
			return false, CustomErrors.UserTakenError{}
		}

		generatedPass, generatePassError := PasswordUtils.GenerateHashedPassword(userRegisterInfo.Password)
		if generatePassError != nil {
			return false, nil
		}

		query, prepareError := connection.Prepare("insert into user_table(name, city, country, password, email) values ($1, $2, $3, $4, $5)")
		if prepareError != nil {
			return false, prepareError
		}

		_, insertUserError := query.Exec(userRegisterInfo.Name, userRegisterInfo.City, userRegisterInfo.Country, generatedPass, userRegisterInfo.Email)
		if insertUserError != nil {
			return false, insertUserError
		}

		return true, nil
	}

	return false, nil
}
