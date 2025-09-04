package routes

import (
	"fmt"
	"net/http"

	"example.com/url-shortner/model"
	"example.com/url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/sign-up", handleSignUp)
}

func handleSignUp(ctx *gin.Context) {
	var userToSignUp model.SignUpUser
	payloadErr := ctx.ShouldBindBodyWithJSON(&userToSignUp)

	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	fmt.Println("User to sign up payload: ", userToSignUp)
	user := model.User{
		Email:    userToSignUp.Email,
		Password: userToSignUp.Password,
	}

	saveErr := user.Save()
	if saveErr != nil {
		utils.HandleValidationError(ctx, saveErr)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User signed up successfully !",
	})
}
