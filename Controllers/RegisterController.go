package Controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
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
		// parse user registration info
		userRegisterInfo, userInfoParseError := getRegistrationInfo(ctx)

		if userInfoParseError != nil {
			returnJsonError(ctx, "error parsing user registration info")
		}

		connection := establishDbConnection(ctr, ctx)

		defer func(connection *sql.DB) {
			_ = connection.Close()
		}(connection)

		_, err := insertUserInDb(userRegisterInfo, connection)
		if err != nil {
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

	if verifyUserData(userRegisterInfo) {
		taken, err := checkIfEmailIsTaken(userRegisterInfo.Email, connection)

		if err != nil {
			return false, err
		}

		if taken {
			return false, errors.New(fmt.Sprintf("user with the email %s is already registered", userRegisterInfo.Email))
		}

		generatedPass, generatePassError := generateHashedPassword(userRegisterInfo.Password)
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

func checkIfEmailIsTaken(email string, connection *sql.DB) (bool, error) {
	query := fmt.Sprintf("select id from user_table where email = '%s'", email)

	var count int

	err := connection.QueryRow(query).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func generateHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func verifyUserData(userModel *Models.UserRegisterModel) bool {
	if verifyPassword(userModel.Password) &&
		len(userModel.Email) > 0 &&
		len(userModel.Name) > 0 &&
		len(userModel.Country) > 0 &&
		len(userModel.City) > 0 {
		return true
	}

	return false
}

func verifyPassword(password string) bool {
	if len(password) >= 6 {
		var upper bool
		var lower bool
		var number bool

		for _, char := range password {
			if unicode.IsLower(char) {
				lower = true
				continue
			}

			if unicode.IsUpper(char) {
				upper = true
				continue
			}

			if unicode.IsNumber(char) {
				number = true
			}
		}

		return lower && upper && number
	}

	return false
}
