package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/mail"
	"kgoel085.com/url-shortner/middleware"
	"kgoel085.com/url-shortner/model"
	"kgoel085.com/url-shortner/utils"
)

func UrlShorterRoutes(router *gin.RouterGroup) {
	router.GET("/", handleRoot)
	router.GET("/:code", handleGetUrls)

	authenticated := router.Group("/url")
	authenticated.Use(middleware.Authenticate)

	authenticated.POST("/register", handleShortUrl)
	authenticated.GET("/list", handleListUrls)
}

func handleRoot(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Welcome to %s URL Shortener Service", config.Config.APP.Name),
	})
}

// @Summary      List User URLs
// @Description  Get all shortened URLs for the authenticated user.
// @Security     BearerAuth
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        status  query  string  false  "Filter by URL status"
// @Success      200  {object}  model.APIResponse{data=model.GetUrlsByUserResponse} "Success" "Example: {\"message\": \"URLs fetched successfully\", \"data\": {\"urls\": [{\"code\": \"abc123\", \"url\": \"https://example.com\"}]}}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Invalid URL status !\"}"
// @Router       /url/list [get]
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

	ctx.JSON(http.StatusOK, model.APIResponse{
		Data:    model.GetUrlsByUserResponse{Urls: urls},
		Message: "URLs fetched successfully",
	})
}

// @Summary      Redirect Short URL
// @Description  Redirects to the original URL using the short code.
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        code  path  string  true  "Short URL code"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"URL is not active\"}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"URL has expired\"}"
// @Router       /{code} [get]
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

// @Summary      Register Short URL
// @Description  Create a new short URL for the authenticated user.
// @Security     BearerAuth
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        createUrl  body  model.CreateShortUrl  true  "Create short URL payload"
// @Success      200  {object}  model.APIResponse{data=model.CreateShortUrlResponse} "Success" "Example: {\"message\": \"Short URL created successfully\", \"data\": {\"short_url\": \"https://short.ly/abc123\"}}"
// @Failure      400  {object}  utils.ErrorResponse "Validation error" "Example: {\"message\": \"Request failed\", \"errors\": [{\"field\": \"url\", \"error\": \"invalid URL\"}]}"
// @Router       /url/register [post]
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
	ctx.JSON(http.StatusOK, model.APIResponse{
		Message: "Short URL created successfully",
		Data:    model.CreateShortUrlResponse{ShortUrl: utils.GetShortUrl(url.Code)},
	})
}
