package routes

import (
	"fmt"
	"net/http"
	"strconv"

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
		"message":       "file uploaded successfully",
		"uploaded_file": fileModel,
	})
}

// /files?cursor=105
// /files?limit=20
func getFiles(context *gin.Context) {

	userId := context.GetInt64("userId")

	cursorStr := context.Query("cursor")
	var cursor int64 = 0 //if cursorStr is empty, default to 0 //we are using the basequery + limit

	if cursorStr != "" {
		var err error
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	limitStr := context.DefaultQuery("limit", "12") // Default to 12 if missing
	limit, err := strconv.ParseInt(limitStr, 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit format", "error": err.Error()})
		return
	}

	files, nextCursor, err := models.GetFilesByUserId(userId, cursor, limit)

	if err != nil {

		context.JSON(http.StatusInternalServerError,
			gin.H{
				"message": "Failed to fetch files",
				"error":   err.Error(),
			})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":        files,
		"next_cursor": nextCursor,
	})
}
