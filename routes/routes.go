package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	registerRoutesForUser(server)
	registerRoutesForFiles(server)
}

func registerRoutesForUser(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)
}

func registerRoutesForFiles(server *gin.Engine) {
	server.POST("/files", createFile)
}
