package Controllers

import "github.com/gin-gonic/gin"

type Controller interface {
	Handle() gin.HandlerFunc
	Route() string
	Method() string
}
