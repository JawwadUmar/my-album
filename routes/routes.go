package routes

import (
	"example.com/my-ablum/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	registerRoutesForUser(server)
	registerRoutesForFiles(server)
}

func registerRoutesForUser(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)
}

func registerRoutesForFiles(server *gin.Engine) {
	authenticationRequiredGroup := server.Group("/")
	authenticationRequiredGroup.Use(middlewares.Authentication)

	authenticationRequiredGroup.POST("/files", createFile)

}
