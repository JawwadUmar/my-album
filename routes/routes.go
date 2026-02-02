package routes

import (
	"example.com/my-ablum/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	middlewares.EnableCors(server)
	registerRoutesForUser(server)
	registerRoutesForFiles(server)
}

func registerRoutesForUser(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.POST("/google", googleLogin)
}

func registerRoutesForFiles(server *gin.Engine) {
	authenticationRequiredGroup := server.Group("/")
	authenticationRequiredGroup.Use(middlewares.Authentication)

	authenticationRequiredGroup.POST("/files", createFile)
	authenticationRequiredGroup.GET("/files", getFiles)
	authenticationRequiredGroup.DELETE("/files/:id", deleteFiles)
	authenticationRequiredGroup.PATCH("/profile/:id", updateProfile)

}
