package main

import (
	"example.com/my-ablum/database"
	"example.com/my-ablum/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	var server *gin.Engine = gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8081")
}
