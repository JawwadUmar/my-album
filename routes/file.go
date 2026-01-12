package routes

import (
	"net/http"

	"example.com/my-ablum/models"
	"github.com/gin-gonic/gin"
)

func createFile(context *gin.Context) {
	var file models.Image
	err := context.ShouldBindBodyWithJSON(&file)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the file object",
			"error":   err.Error(),
		})
		return
	}

	file.CreatedBy = context.GetInt64("userId")
	err = file.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"File": file})

}
