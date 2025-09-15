package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	customValidator "kgoel085.com/url-shortner/validator"
)

// ErrorDetail represents a single field error
type ErrorDetail struct {
	Field string `json:"field" example:"email"`
	Error string `json:"error" example:"invalid email"`
}

// ErrorResponse is a generic error response for API errors
type ErrorResponse struct {
	Message string        `json:"message" example:"Request failed"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

func HandleValidationError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]ErrorDetail, len(ve))
		for i, fe := range ve {
			out[i] = ErrorDetail{
				Field: fe.Field(),
				Error: customValidator.MsgForTag(fe),
			}
		}

		Log.WithFields(logrus.Fields{
			"url":   ctx.Request.URL.Path,
			"error": out,
		}).Error("Request failed")

		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Request failed",
			Errors:  out,
		})
		return
	}

	// Non-validation error
	Log.WithFields(logrus.Fields{
		"url":   ctx.Request.URL.Path,
		"error": err.Error(),
	}).Error("Request failed")

	ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Message: err.Error(),
	})
}
