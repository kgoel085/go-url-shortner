package model

import (
	"database/sql"
	"fmt"
	"time"

	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/db"
	"kgoel085.com/url-shortner/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" binding:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpwd"`
}

type UserOtp struct {
	OtpToken string `json:"otp_token" binding:"required"`
	OtpCode  string `json:"otp_code" binding:"required"`
}

type SignUpUser struct {
	UserCredentials
	UserOtp
}

type LoginUser struct {
	UserCredentials
	UserOtp
}

type LoginUserResponse struct {
	Token        string `json:"token" example:"JWT Token"`
	RefreshToken string `json:"refresh_token" example:"JWT Refresh Token"`
}

func (u *User) Save() error {
	userByEmail, userByEmailErr := GetUserByEmail(u.Email)
	if userByEmail.Email == u.Email || userByEmailErr == nil {
		return fmt.Errorf("user already exists !")
	}

	hashedPwd, hashPwdErr := utils.HashPwd(u.Password)
	if hashPwdErr != nil {
		errStr := fmt.Sprintf("Error while trying to hash - %s !", hashPwdErr.Error())
		return fmt.Errorf("%s", errStr)
	}

	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id, created_at`

	logStr := fmt.Sprintf("Save user in DB : %s, Email: %s, Timestamp: %s", query, u.Email, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, u.Email, hashedPwd, time.Now().UTC()).Scan(&u.ID, &u.CreatedAt)

	return rowErr
}

func (u *User) GenerateJWT() (string, error) {
	return utils.GenerateLoginJWT(u.ID)
}

func (u *User) GenerateRefreshJWT() (string, error) {
	token, tokenErr := utils.GenerateRefreshJWT(u.ID)
	if tokenErr != nil {
		return "", tokenErr
	}

	encryptedToken, encErr := utils.Encrypt(token)
	if encErr != nil {
		return "", encErr
	}

	// Use encrypted token for storage
	token = encryptedToken

	// Save refresh token in DB
	expiryInMin := config.Config.JWT.RefreshExpiryMinutes
	expiresAt := time.Now().Add(time.Duration(expiryInMin) * time.Minute)

	query := `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4)`

	logStr := fmt.Sprintf("Save refresh token in DB : %s, UserID: %d, ExpiresAt: %s", query, u.ID, expiresAt)
	utils.Log.Info(logStr)

	_, rowErr := db.DB.Exec(query, u.ID, token, expiresAt, time.Now().UTC())
	if rowErr != nil {
		return "", rowErr
	}

	return token, nil
}

func (u *User) ValidateCredentials() error {
	userByEmail, userByEmailErr := GetUserByEmail(u.Email)
	if userByEmailErr != nil {
		return userByEmailErr
	}

	u.ID = userByEmail.ID

	userPwdHash := userByEmail.Password
	userPwd := u.Password

	// Check pwd
	isValidPwd := utils.CheckHashPwd(userPwd, userPwdHash)
	if !isValidPwd {
		return fmt.Errorf("Invalid password !")
	}

	return nil
}

func GetRefreshTokenByToken(token string) (UserRefreshToken, error) {
	var refreshToken UserRefreshToken

	query := `SELECT id, token, expires_at, user_id, created_at FROM refresh_tokens WHERE token=$1`

	logStr := fmt.Sprintf("Get refresh token from DB : %s, Timestamp: %s", query, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, token).Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.UserID, &refreshToken.CreatedAt)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return refreshToken, fmt.Errorf("Refresh token not found")
		}
		errStr := fmt.Sprintf("Error while trying to get refresh token - %s !", rowErr.Error())
		return refreshToken, fmt.Errorf("%s", errStr)
	}

	return refreshToken, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	query := `SELECT id, email, password, created_at FROM users WHERE email ILIKE $1`
	row := db.DB.QueryRow(query, email)

	logStr := fmt.Sprintf("Check User via EMAIL: %s, %s", query, email)
	utils.Log.Info(logStr)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("User not found")
		}
		return user, err
	}

	return user, nil
}

func GetUserById(id int64) (User, error) {
	var user User

	query := `SELECT id, email, password, created_at FROM users WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	logStr := fmt.Sprintf("Check User via ID: %s, %d", query, id)
	utils.Log.Info(logStr)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("User not found")
		}
		return user, err
	}

	return user, nil
}
