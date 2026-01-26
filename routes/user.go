package routes

import (
	"net/http"

	"example.com/my-ablum/models"
	storage "example.com/my-ablum/storage/1"
	"example.com/my-ablum/utility"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email     string `form:"email" binding:"required,email"`
	Password  string `form:"password" binding:"required,min=8"`
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
}

func signup(context *gin.Context) {
	var signupRequest SignupRequest
	err := context.ShouldBind(&signupRequest) //not with JSON as it will be a form data :)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the user object",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	user.Email = signupRequest.Email
	user.Password = signupRequest.Password
	user.FirstName = signupRequest.FirstName
	user.LastName = signupRequest.LastName

	fileHeader, err := context.FormFile("profile_pic")

	if err == nil {
		storageKey := "profile-pics/" + user.Email
		err = storage.StoreFileInS3(fileHeader, storageKey)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upload profile picture",
			})
			return
		}

		user.ProfilePic = &storageKey
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is some problem saving the user",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": user})
}

func login(context *gin.Context) {

	var loginRequest LoginRequest
	err := context.ShouldBindBodyWithJSON(&loginRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the user object",
			"error":   err.Error(),
		})

		return
	}

	var user models.User
	user.Email = loginRequest.Email
	user.Password = loginRequest.Password

	err = user.ValidateCredential() //UserId is updated here

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Password",
			"error":   err.Error(),
		})

		return
	}

	token, err := utility.GenerateToken(user.Email, user.UserId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem in generating jwt token",
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully login",
		"token":   token,
		"user":    user,
	})

}
