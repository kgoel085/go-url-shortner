package routes

import (
	"fmt"
	"net/http"

	"example.com/url-shortner/mail"
	"example.com/url-shortner/model"
	"example.com/url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/sign-up", handleSignUp)
	router.POST("/verify-credentials", handleVerifyCredentials)
	router.POST("/login", handleLogin)
}

func handleLogin(ctx *gin.Context) {
	var loginUser model.LoginUser
	payloadErr := ctx.ShouldBindJSON(&loginUser)

	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	fmt.Println("User to login payload: ", loginUser)
	user := model.User{
		Email:    loginUser.Email,
		Password: loginUser.Password,
	}

	userCredsErr := user.ValidateCredentials() // Validate email and password
	if userCredsErr != nil {
		utils.HandleValidationError(ctx, userCredsErr)
		return
	}

	utils.Log.Info("User credentials validated, proceeding to OTP verification...")
	otpVerify := model.VerifyOtp{
		Token:  loginUser.OtpToken,
		Otp:    loginUser.OtpCode,
		Action: string(model.OtpActionTypeLogin),
	}

	otpErr := otpVerify.VerifyWithUpdate() // Validate OTP and update its status to 'success' if valid
	if otpErr != nil {
		utils.HandleValidationError(ctx, otpErr)
		return
	}

	utils.Log.Info("OTP verified successfully, generating JWT...")
	token, tokenErr := user.GenerateJWT()
	if tokenErr != nil {
		utils.HandleValidationError(ctx, tokenErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully !",
		"token":   token,
	})
}

func handleVerifyCredentials(ctx *gin.Context) {
	var userCreds model.UserCredentials
	payloadErr := ctx.ShouldBindBodyWithJSON(&userCreds)

	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	fmt.Println("User to verify credentials: ", userCreds)
	user := model.User{
		Email:    userCreds.Email,
		Password: userCreds.Password,
	}

	userErr := user.ValidateCredentials()
	if userErr != nil {
		utils.HandleValidationError(ctx, userErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User credentials are valid.",
	})
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

	go mail.SendSignedUpUserMail(user)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User signed up successfully !",
	})
}
