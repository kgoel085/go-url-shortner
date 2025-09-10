package model

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"example.com/url-shortner/config"
	"example.com/url-shortner/db"
	"example.com/url-shortner/utils"
)

type OtpType string
type OtpActionType string
type OtpStatus string

const (
	OtpStatusPending OtpStatus = "pending"
	OtpStatusSuccess OtpStatus = "success"
	OtpStatusExpire  OtpStatus = "expire"
)

const (
	OtpActionTypeLogin         OtpActionType = "login"
	OtpActionTypeSignUp        OtpActionType = "signup"
	OtpActionTypeResetPassword OtpActionType = "reset_password"
)

const (
	OtpTypeEmail OtpType = "email"
	OtpTypePhone OtpType = "phone"
)

type Otp struct {
	ID        int64         `json:"id"`
	Key       string        `json:"key" binding:"required"`
	Type      OtpType       `json:"type" binding:"required"`
	Action    OtpActionType `json:"action" binding:"required"`
	OtpCode   string        `json:"otp"`
	Token     string        `json:"token"`
	CreatedAt time.Time     `json:"created_at"`
	Status    OtpStatus     `json:"status"`
}

type SendOtp struct {
	Type   OtpType       `json:"type" binding:"required"`
	Action OtpActionType `json:"action" binding:"required"`
	Key    string        `json:"key" binding:"required,email"`
}

type VerifyOtp struct {
	Token  string `json:"token" binding:"required"`
	Otp    string `json:"otp" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type SendOTPResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (otpVerify *VerifyOtp) Verify() error {
	return otpVerify.verifyInternal(false)
}

func (otpVerify *VerifyOtp) VerifyWithUpdate() error {
	return otpVerify.verifyInternal(true)
}

func (otpVerify *VerifyOtp) verifyInternal(performUpdate bool) error {
	var otp Otp
	row := db.DB.QueryRow("SELECT id, otp, status, created_at FROM otp WHERE token = $1 AND action = $2 AND status = $3", otpVerify.Token, otpVerify.Action, OtpStatusPending)

	scanErr := row.Scan(&otp.ID, &otp.OtpCode, &otp.Status, &otp.CreatedAt)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return fmt.Errorf("Invalid OTP details !")
		}
		return scanErr
	}

	if otp.ID == 0 {
		return fmt.Errorf("Invalid OTP token")
	}

	if otp.OtpCode != otpVerify.Otp {
		return fmt.Errorf("Invalid OTP code")
	}

	if time.Since(otp.CreatedAt) > time.Minute*time.Duration(config.Config.OTP.ExpiryMinutes) {
		// Expire the OTP
		_, updateErr := db.DB.Exec("UPDATE otp SET status=$1 WHERE id=$2", OtpStatusExpire, otp.ID)
		if updateErr != nil {
			utils.Log.Error("Error expiring OTP: ", updateErr)
		}
		return fmt.Errorf("OTP has expired. Please request a new one.")
	}

	if performUpdate { // Mark OTP as success only if performUpdate is true
		updateErr := otp.UpdateStatus(OtpStatusSuccess)
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (otp *Otp) UpdateStatus(status OtpStatus) error {
	// Mark OTP as success
	utils.Log.Info("Updating OTP status to ", status, "UPDATE otp SET status=$1 WHERE id=$2")
	_, updateErr := db.DB.Exec("UPDATE otp SET status=$1 WHERE id=$2", status, otp.ID)
	if updateErr != nil {
		return updateErr
	}

	otp.Status = status
	return nil
}

func (otp *Otp) Generate() error {
	// OTP Type checks
	switch {
	case otp.Action == OtpActionTypeLogin && otp.Type == OtpTypeEmail:
		{
			{
				userEmail, userEmailErr := getUserByEmail(otp.Key)
				if userEmailErr != nil {
					return userEmailErr
				}
				if userEmail.ID == 0 {
					return fmt.Errorf("No user found with email %s", otp.Key)
				}
			}
		}
	}

	otp.generateOtp()

	// Check if any other OTP exists with same action and type recently
	checkErr := otp.checkExistingOtp()
	if checkErr != nil {
		return checkErr
	}

	insertQuery := `INSERT INTO otp (key, type, action, otp, created_at, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, token`
	logStr := fmt.Sprintf("Insert OTP in DB : %s, Key: %s, Type: %s, Action: %s, Timestamp: %s", insertQuery, otp.Key, otp.Type, otp.Action, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(insertQuery, otp.Key, otp.Type, otp.Action, otp.OtpCode, otp.CreatedAt, OtpStatusPending).Scan(&otp.ID, &otp.Token)
	if rowErr != nil {
		return rowErr
	}

	return nil

}

func (otp *Otp) checkExistingOtp() error {
	row, rowErr := db.DB.Query("SELECT id, otp, created_at FROM otp WHERE key = $1 AND type=$2 AND action=$3 AND status = $4 ORDER BY created_at DESC", otp.Key, otp.Type, otp.Action, OtpStatusPending)
	if rowErr != nil {
		return rowErr
	}
	defer row.Close()

	if row.Next() {
		existingOtp := Otp{}
		scanErr := row.Scan(&existingOtp.ID, &existingOtp.OtpCode, &existingOtp.CreatedAt)
		if scanErr != nil {
			return scanErr
		}
		// If OTP was sent within last specified minute, do not send another one
		if time.Since(existingOtp.CreatedAt) < time.Minute*time.Duration(config.Config.OTP.ExpiryMinutes) {
			errStr := fmt.Sprintf("OTP already sent recently at %s. Please wait before requesting a new one.", existingOtp.CreatedAt.Format(config.TIME_FORMAT))
			return fmt.Errorf("%s", errStr)
		} else {
			// Expire the previous OTP
			_, updateErr := db.DB.Exec("UPDATE otp SET status=$1 WHERE id=$2", OtpStatusExpire, existingOtp.ID)
			if updateErr != nil {
				return updateErr
			}
		}
	}

	return nil
}

func (otp *Otp) generateOtp() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newOtp := fmt.Sprintf("%06d", r.Intn(1000000))

	otp.OtpCode = newOtp
	otp.CreatedAt = time.Now().UTC()

	return newOtp
}

func (ot OtpType) IsValid() bool {
	switch ot {
	case OtpTypeEmail, OtpTypePhone:
		return true
	}
	return false
}

func (ot OtpActionType) IsValid() bool {
	switch ot {
	case OtpActionTypeLogin, OtpActionTypeSignUp, OtpActionTypeResetPassword:
		return true
	}
	return false
}

func (ot OtpStatus) IsValid() bool {
	switch ot {
	case OtpStatusPending, OtpStatusSuccess, OtpStatusExpire:
		return true
	}
	return false
}
