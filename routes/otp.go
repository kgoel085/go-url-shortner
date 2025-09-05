package routes

import (
	"fmt"
	"net/http"

	"example.com/url-shortner/mail"
	"example.com/url-shortner/model"
	"example.com/url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func OtpRoutes(router *gin.RouterGroup) {
	router.POST("/send", handleSendOTP)
	router.POST("/verify", handleVerifyOTP)
}

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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
	})
}

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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"id":      otp.ID,
		"token":   otp.Token,
	})
}
