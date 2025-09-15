package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/mail"
	"kgoel085.com/url-shortner/model"
	"kgoel085.com/url-shortner/utils"
)

func OtpRoutes(router *gin.RouterGroup) {
	router.POST("/send", handleSendOTP)
	router.POST("/verify", handleVerifyOTP)
}

// @Summary      Verify OTP
// @Description  Verifies the OTP code sent to the user.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        otpVerifyRequest  body  model.VerifyOtp  true  "OTP verify payload"
// @Success      200  {object}  model.APIResponse "Success" "Example: {\"message\": \"OTP verified successfully\"}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Invalid OTP code\"}"
// @Router       /otp/verify [post]
func handleVerifyOTP(ctx *gin.Context) {
	var otpVerifyRequest model.VerifyOtp
	payloadErr := ctx.ShouldBindBodyWithJSON(&otpVerifyRequest)

	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	otpErr := otpVerifyRequest.Verify()
	if otpErr != nil {
		utils.HandleValidationError(ctx, otpErr)
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Message: "OTP verified successfully",
	})
}

// @Summary      Send OTP
// @Description  Sends an OTP to the user for verification.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        otpRequest  body  model.SendOtp  true  "OTP request payload"
// @Success      200  {object}  model.APIResponse{data=model.SendOTPResponse} "Success" "Example: {\"message\": \"OTP sent successfully\", \"data\": {\"id\": \"123\", \"token\": \"abcde12345\"}}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Invalid OTP type or action type\"}"
// @Router       /otp/send [post]
func handleSendOTP(ctx *gin.Context) {
	var otpRequest model.SendOtp
	payloadErr := ctx.ShouldBindBodyWithJSON(&otpRequest)

	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	if (!otpRequest.Type.IsValid()) || (!otpRequest.Action.IsValid()) {
		utils.HandleValidationError(ctx, fmt.Errorf("Invalid OTP type or action type"))
		return
	}

	otp := model.Otp{
		Key:    otpRequest.Key,
		Type:   otpRequest.Type,
		Action: otpRequest.Action,
	}
	otpErr := otp.Generate()
	if otpErr != nil {
		utils.HandleValidationError(ctx, otpErr)
		return
	}

	// Send Email
	go mail.SendOtpUserMail(otp)
	ctx.JSON(http.StatusOK, model.APIResponse{
		Message: "OTP sent successfully",
		Data:    model.SendOTPResponse{ID: fmt.Sprintf("%d", otp.ID), Token: otp.Token},
	})
}
