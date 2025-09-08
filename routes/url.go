package routes

import (
	"fmt"
	"net/http"
	"time"

	"example.com/url-shortner/config"
	"example.com/url-shortner/mail"
	"example.com/url-shortner/middleware"
	"example.com/url-shortner/model"
	"example.com/url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func UrlShorterRoutes(router *gin.RouterGroup) {
	router.GET("/:code", handleGetUrls)

	authenticated := router.Group("/url")
	authenticated.Use(middleware.Authenticate)

	authenticated.POST("/register", handleShortUrl)
	authenticated.GET("/list", handleListUrls)
}

func handleListUrls(ctx *gin.Context) {
	loggedInUser := ctx.GetInt64(config.JWT_LOGGED_IN_USER)
	utils.Log.Info("Get URLs for user:", loggedInUser)

	filters := model.GetUrlByUserFilter{}
	status := ctx.Query("status")
	if status != "" {
		urlStatus := model.UrlStatus(status)
		if !urlStatus.IsValid() {
			utils.HandleValidationError(ctx, fmt.Errorf("Invalid URL status !"))
			return
		}
		filters.Status = urlStatus
	}

	urls, urlsErr := model.GetUrlsByUser(loggedInUser, filters)
	if urlsErr != nil {
		utils.HandleValidationError(ctx, urlsErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"urls": urls,
	})
}

func handleGetUrls(ctx *gin.Context) {
	code := ctx.Param("code")
	utils.Log.Info("Get URL by code:", code)

	url, urlErr := model.GetUrlByCode(code)
	if urlErr != nil {
		utils.HandleValidationError(ctx, urlErr)
		return
	}

	if url.Status != model.UrlStatusActive {
		utils.HandleValidationError(ctx, fmt.Errorf("URL is not active"))
		return
	}

	if !url.ExpiryAt.IsZero() && url.ExpiryAt.Before(time.Now()) {
		updateErr := url.UpdateStatus(model.UrlStatusExpired)
		if updateErr != nil {
			utils.Log.Error("Failed to update URL status to expired:", updateErr)
		}
		utils.HandleValidationError(ctx, fmt.Errorf("URL has expired"))
		return
	}

	// Log analytics data
	analytics := model.Analytics{
		UrlID:     url.ID,
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		Referrer:  ctx.Request.Referer(),
	}
	go func() {
		if err := analytics.Save(); err != nil {
			utils.Log.Error("Failed to save analytics data:", err)
		}
	}()

	utils.Log.Info("Redirecting to URL:", url.Url)
	ctx.Redirect(http.StatusPermanentRedirect, url.Url)
}

func handleShortUrl(ctx *gin.Context) {
	var createUrl model.CreateShortUrl
	payloadErr := ctx.ShouldBindJSON(&createUrl)
	if payloadErr != nil {
		utils.HandleValidationError(ctx, payloadErr)
		return
	}

	loggedInUser := ctx.GetInt64(config.JWT_LOGGED_IN_USER)
	createUrl.UserID = loggedInUser

	// Validate payload
	url, urlErr := createUrl.Validate()
	if urlErr != nil {
		utils.HandleValidationError(ctx, urlErr)
		return
	}

	urlErr = url.Save()
	if urlErr != nil {
		utils.HandleValidationError(ctx, urlErr)
		return
	}

	go mail.SendShortUrlUserMail(url)
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Short URL created successfully",
		"short_url": utils.GetShortUrl(url.Code),
	})
}
