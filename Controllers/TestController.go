package Controllers

import (
	"github.com/Dontpingforgank/AuthenticationService/Service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type testController struct {
	service Service.Service
}

func NewTestController(service Service.Service) Controller {
	return &testController{
		service: service,
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
