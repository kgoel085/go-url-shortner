package utils

import (
	"errors"
	"net/http"

	customValidator "example.com/url-shortner/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func HandleValidationError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]map[string]string, len(ve))
		for i, fe := range ve {
			out[i] = map[string]string{
				"field": fe.Field(),
				"error": customValidator.MsgForTag(fe),
			}
		}

		Log.WithFields(logrus.Fields{
			"url":   ctx.Request.URL.Path,
			"error": out,
		}).Error("Request failed")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Request failed",
			"errors":  out,
		})
		return
	}

	// Non-validation error
	Log.WithFields(logrus.Fields{
		"url":   ctx.Request.URL.Path,
		"error": err.Error(),
	}).Error("Request failed")

	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": err.Error(),
	})
}
