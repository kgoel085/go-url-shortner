package utils

import (
	"fmt"
	"time"

	"example.com/url-shortner/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int64) (string, error) {
	expiryInMin := config.Config.JWT.ExpiryMinutes
	jwtSecretKey := config.Config.JWT.SecretKey

	fmt.Println("JWT GENERATION CONFIG ::", expiryInMin, jwtSecretKey)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Minute * time.Duration(expiryInMin)).Unix(),
	})

	token, err := jwtToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return token, err
	}

	return token, nil
}

func ValidateJWT(token string) (int64, error) {
	jwtSecretKey := config.Config.JWT.SecretKey
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

	// email, _ := claims["email"].(string)
	userId, _ := claims["userId"].(float64)

	return int64(userId), nil
}
