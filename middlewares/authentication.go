package middlewares

import (
	"net/http"

	"example.com/my-ablum/utility"
	"github.com/gin-gonic/gin"
)

func Authentication(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized, empty token",
		})

		return
	}

	err := utility.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Could not verify token",
		})

		return
	}

	userId, err := utility.GetUserIdFromToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get the user Id of Event-creator",
		})

		return
	}

	context.Set("userId", userId)
	context.Next()

}
