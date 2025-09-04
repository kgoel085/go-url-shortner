package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/url-shortner/db"
	"example.com/url-shortner/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type SignUpUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,strongpwd"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	userByEmail, userByEmailErr := getUserByEmail(u.Email)
	if userByEmail.Email == u.Email || userByEmailErr == nil {
		return errors.New("User already exists !")
	}

	hashedPwd, hashPwdErr := utils.HashPwd(u.Password)
	if hashPwdErr != nil {
		errStr := fmt.Sprintf("Error while trying to hash - %s !", hashPwdErr.Error())
		return errors.New(errStr)
	}

	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id, created_at`

	logStr := fmt.Sprintf("Save user in DB : %s, Email: %s, Timestamp: %s", query, u.Email, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, u.Email, hashedPwd, time.Now().UTC()).Scan(&u.ID, &u.CreatedAt)

	return rowErr
}

func getUserByEmail(email string) (User, error) {
	var user User

	query := `SELECT * FROM users WHERE email ILIKE $1`
	row := db.DB.QueryRow(query, email)

	logStr := fmt.Sprintf("Check User via EMAIL: %s, %s", query, email)
	utils.Log.Info(logStr)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("User not found")
		}
		return user, err
	}

	return user, nil
}
