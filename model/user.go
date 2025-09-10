package model

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/url-shortner/db"
	"example.com/url-shortner/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" binding:"email"`
	Password  string    `json:"-"`
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
	Token string `json:"token" example:"JWT Token"`
}

func (u *User) Save() error {
	userByEmail, userByEmailErr := getUserByEmail(u.Email)
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
	return utils.GenerateJWT(u.ID)
}

func (u *User) ValidateCredentials() error {
	userByEmail, userByEmailErr := getUserByEmail(u.Email)
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

func getUserByEmail(email string) (User, error) {
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
