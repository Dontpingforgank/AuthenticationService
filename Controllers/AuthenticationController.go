package Controllers

import (
	"github.com/Dontpingforgank/AuthenticationService/Database"
	"github.com/Dontpingforgank/AuthenticationService/Logger"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type AuthenticationController struct {
	configs       *Models.Config
	loggerFactory Logger.LoggerFactory
	dbFactory     Database.DatabaseFactory
}

func NewAuthenticationController(configs *Models.Config, loggerFactory Logger.LoggerFactory, dbFactory Database.DatabaseFactory) Controller {
	return &AuthenticationController{
		configs:       configs,
		loggerFactory: loggerFactory,
		dbFactory:     dbFactory,
	}
}

func (ctr AuthenticationController) Handle() gin.HandlerFunc {
	logger, closeLogger, err := ctr.loggerFactory.NewLogger()
	if err != nil {
		panic(err)
	}

	defer closeLogger()

	return func(context *gin.Context) {
		tokenHeader := context.GetHeader("Authorization")

		if tokenHeader != "" {
			token, err := jwt.ParseWithClaims(tokenHeader, &Models.UserClaimsModel{}, func(tkn *jwt.Token) (interface{}, error) {
				if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
					context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"Success": false,
						"Message": "Token Invalid",
					})
				}
				return []byte(ctr.configs.JwtSecret), nil
			})

			if claims, ok := token.Claims.(*Models.UserClaimsModel); ok && token.Valid {
				if claims.Exp <= time.Now().Unix() {
					context.JSON(http.StatusUnauthorized, gin.H{
						"Success": false,
						"Message": "Token expired",
					})
				} else {
					context.JSON(http.StatusOK, gin.H{
						"Success": true,
						"Message": "Token is valid",
					})
				}
			} else {
				logger.Error(err.Error())
				context.AbortWithStatus(http.StatusUnauthorized)
			}
		}
	}
}

func (AuthenticationController) Route() string {
	return "/Authenticate"
}

func (AuthenticationController) Method() string {
	return http.MethodGet
}
