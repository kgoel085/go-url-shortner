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
	router.POST("/login", handleLogin)
	router.POST("/verify-credentials", handleVerifyCredentials)
}

// @Summary      User Login
// @Description  Login with email, password, and OTP. Returns JWT token on success.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginUser  body  model.LoginUser  true  "Login payload"
// @Success      200  {object}  model.APIResponse{data=model.LoginUserResponse} "Success" "Example: {\"message\": \"User logged in successfully !\", \"data\": {\"token\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"}}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Request failed\", \"errors\": [{\"field\": \"email\", \"error\": \"invalid email\"}]}"
// @Router       /login [post]
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

	ctx.JSON(http.StatusOK, model.APIResponse{
		Message: "User logged in successfully !",
		Data:    model.LoginUserResponse{Token: token},
	})
}

// @Summary      Verify User Credentials
// @Description  Verifies user email and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        userCreds  body  model.UserCredentials  true  "User credentials payload"
// @Success      200  {object}  model.APIResponse "Success" "Example: {\"message\": \"User credentials are valid.\"}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Request failed\", \"errors\": [{\"field\": \"password\", \"error\": \"password too weak\"}]}"
// @Router       /verify-credentials [post]
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

	ctx.JSON(http.StatusOK, model.APIResponse{
		Message: "User credentials are valid.",
	})
}

// @Summary      User Sign Up
// @Description  Register a new user with email, password, and OTP verification.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        userToSignUp  body  model.SignUpUser  true  "Sign up payload"
// @Success      201  {object}  model.APIResponse "Success" "Example: {\"message\": \"User signed up successfully !\"}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Request failed\", \"errors\": [{\"field\": \"otp_code\", \"error\": \"invalid OTP\"}]}"
// @Router       /sign-up [post]
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

	otpVerify := model.VerifyOtp{
		Token:  userToSignUp.OtpToken,
		Otp:    userToSignUp.OtpCode,
		Action: string(model.OtpActionTypeSignUp),
	}
	otpErr := otpVerify.Verify()
	if otpErr != nil {
		utils.HandleValidationError(ctx, otpErr)
		return
	}

	utils.Log.Info("OTP verified successfully")

	saveErr := user.Save()
	if saveErr != nil {
		utils.HandleValidationError(ctx, saveErr)
		return
	}

	otpVerifyErr := otpVerify.VerifyWithUpdate() // Mark OTP as success, but ignore any error
	if otpVerifyErr != nil {
		utils.Log.Error("Error updating OTP status to success: \n", otpVerifyErr)
	}

	utils.Log.Info("User signed up successfully: ", user.Email)

	go mail.SendSignedUpUserMail(user)
	ctx.JSON(http.StatusCreated, model.APIResponse{
		Message: "User signed up successfully !",
	})
}
