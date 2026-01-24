package routes

import (
	"net/http"

	"example.com/my-ablum/models"
	"example.com/my-ablum/utility"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=8"`
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	ProfilePic *string `json:"profile_pic"`
}

func signup(context *gin.Context) {
	var signupRequest SignupRequest
	err := context.ShouldBindJSON(&signupRequest)

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
	user.ProfilePic = signupRequest.ProfilePic

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
		"id":      user.UserId,
	})

}
