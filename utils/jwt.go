package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"kgoel085.com/url-shortner/config"
)

type JwtType string

const (
	LoginJwtType   JwtType = "login"
	RefreshJwtType JwtType = "refresh"
)

type GenerateJwtWithClaims struct {
	Claims      jwt.MapClaims `binding:"required"`
	SecretKey   string        `binding:"required"`
	ExpiryInMin int64         `binding:"required"`
}

func GenerateLoginJWT(userID int64) (string, error) {
	expiryInMin := config.Config.JWT.ExpiryMinutes
	jwtSecretKey := config.Config.JWT.SecretKey

	Log.Info("JWT GENERATION CONFIG ::", expiryInMin, jwtSecretKey)
	payload := GenerateJwtWithClaims{
		Claims: jwt.MapClaims{
			"userId": userID,
		},
		SecretKey:   jwtSecretKey,
		ExpiryInMin: expiryInMin,
	}

	token, err := generateJwtWithClaims(payload)
	if err != nil {
		return token, err
	}

	return token, nil
}

func GenerateRefreshJWT(userID int64) (string, error) {
	expiryInMin := config.Config.JWT.RefreshExpiryMinutes
	jwtSecretKey := config.Config.JWT.RefreshSecretKey

	Log.Info("Refresh JWT GENERATION CONFIG ::", expiryInMin, jwtSecretKey)
	payload := GenerateJwtWithClaims{
		Claims: jwt.MapClaims{
			"userId": userID,
		},
		SecretKey:   jwtSecretKey,
		ExpiryInMin: expiryInMin,
	}

	token, err := generateJwtWithClaims(payload)
	if err != nil {
		return token, err
	}

	return token, nil
}

func generateJwtWithClaims(payload GenerateJwtWithClaims) (string, error) {
	claims := payload.Claims
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(payload.ExpiryInMin)).Unix()

	Log.Info("Generating JWT with claims: ", claims)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(payload.SecretKey))
	if err != nil {
		return token, err
	}

	return token, nil
}

func ValidateJWT(token string, tokenType JwtType) (int64, error) {
	if tokenType == "" {
		return 0, fmt.Errorf("token type is required")
	}

	jwtSecretKey := config.Config.JWT.SecretKey
	if tokenType == RefreshJwtType {
		jwtSecretKey = config.Config.JWT.RefreshSecretKey
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", fmt.Errorf("Unexpected signing method !")
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("Could not parse token - %s!", err.Error())
	}

	if !parsedToken.Valid {
		return 0, fmt.Errorf("Invalid Token !")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("Invalid claims !")
	}

	userId, _ := claims["userId"].(float64)

	return int64(userId), nil
}
