package routes

import (
	"fmt"
	"net/http"

	"example.com/my-ablum/models"
	storage "example.com/my-ablum/storage/1"
	"github.com/gin-gonic/gin"
)

func createFile(context *gin.Context) {
	fileHeader, err := context.FormFile("file")

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse the file correctly",
		})

		return
	}

	var fileModel models.File

	fileModel.FileName = fileHeader.Filename
	fileModel.FileSize = fileHeader.Size
	fileModel.MimeType = fileHeader.Header.Get("Content-Type")
	fileModel.CreatedBy = context.GetInt64("userId")

	fileModel.StorageKey = fmt.Sprintf("%d/%s", fileModel.CreatedBy, fileModel.FileName)

	// filePath := "./storage/" + fileModel.StorageKey

	err = storage.StoreFileInS3(fileHeader, fileModel.StorageKey)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Problem saving to S3",
			"error":   err.Error(),
		})

		return
	}

	// Create folder if not exists (storage/1/, storage/2/, etc)
	// err = os.MkdirAll(filepath.Dir(filePath), 0755)
	// if err != nil {
	// 	context.JSON(500, gin.H{"error": "failed to create storage folder"})
	// 	return
	// }

	// // Save file to disk
	// err = context.SaveUploadedFile(fileHeader, filePath)
	// if err != nil {
	// 	context.JSON(500, gin.H{"error": "failed to save file"})
	// 	return
	// }

	err = fileModel.Save()

	if err != nil {
		context.JSON(500, gin.H{"error": "database save failed"})
		return
	}

	context.JSON(200, gin.H{
		"file_id":   fileModel.FileId,
		"file_name": fileModel.FileName,
	})
}
