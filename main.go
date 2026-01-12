package main

import (
	"example.com/my-ablum/database"
	"example.com/my-ablum/routes"
	storage "example.com/my-ablum/storage/1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	storage.InitS3()

	database.Init()
	var server *gin.Engine = gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8081")
}
