package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"example.com/url-shortner/db"
	"example.com/url-shortner/utils"
)

type UrlStatus string

const (
	UrlStatusActive   UrlStatus = "active"
	UrlStatusInactive UrlStatus = "inactive"
	UrlStatusDeleted  UrlStatus = "deleted"
	UrlStatusExpired  UrlStatus = "expired"
)

func (us UrlStatus) IsValid() bool {
	switch us {
	case UrlStatusActive, UrlStatusInactive, UrlStatusDeleted, UrlStatusExpired:
		return true
	}
	return false
}

type Url struct {
	ID         int64     `json:"id" binding:"required"`
	UserID     int64     `json:"user_id" binding:"required"`
	Url        string    `json:"url" binding:"required,http_url"`
	Code       string    `json:"code" binding:"required,alphanum"`
	Status     UrlStatus `json:"status" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"required"`
	ClickCount int64     `json:"click_count"`
	ExpiryAt   time.Time `json:"expires_at"`
}

type CreateShortUrl struct {
	Url      string    `json:"url" binding:"required,http_url"`
	ExpiryAt time.Time `json:"expires_at" binding:"omitempty"`
	Code     string    `json:"code" binding:"omitempty,alphanum"`
	UserID   int64     `json:"user_id"`
}

type GetUrlByUserFilter struct {
	Status UrlStatus `json:"status" binding:"omitempty,oneof=active inactive deleted expired"`
}

type UrlWithShortCode struct {
	Url
	ShortUrl string `json:"short_url"`
}

func GetUrlByCode(code string) (Url, error) {
	var url Url
	query := `SELECT id, user_id, url, code, status, created_at, expiry_at FROM url WHERE code=$1`

	logStr := fmt.Sprintf("Get URL by Code from DB : %s, Code: %s, Timestamp: %s", query, code, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, code).Scan(&url.ID, &url.UserID, &url.Url, &url.Code, &url.Status, &url.CreatedAt, &url.ExpiryAt)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return url, fmt.Errorf("no URL found for the provided code")
		}
		errStr := fmt.Sprintf("Error while trying to get URL by code - %s !", rowErr.Error())
		return url, fmt.Errorf(errStr)
	}

	return url, nil
}

func (u *Url) UpdateStatus(status UrlStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid URL status")
	}
	query := `UPDATE url SET status=$1 WHERE id=$2`

	logStr := fmt.Sprintf("Update URL status in DB : %s, ID: %d, New Status: %s, Timestamp: %s", query, u.ID, status, time.Now().UTC())
	utils.Log.Info(logStr)

	_, execErr := db.DB.Exec(query, u.Status, u.ID)
	if execErr != nil {
		errStr := fmt.Sprintf("Error while trying to update URL status - %s !", execErr.Error())
		return fmt.Errorf(errStr)
	}

	u.Status = status
	return nil
}

func GetUrlsByUser(userID int64, filter GetUrlByUserFilter) ([]UrlWithShortCode, error) {
	var urls []UrlWithShortCode

	var args []interface{}
	var conditions []string

	// Base query
	query := `SELECT id, user_id, url, code, status, created_at, expiry_at FROM url WHERE user_id=$1`
	args = append(args, userID)

	// Add status filter if provided
	if filter.Status.IsValid() {
		conditions = append(conditions, fmt.Sprintf("status=$%d", len(args)+1))
		args = append(args, filter.Status)
	}

	// Append additional conditions
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get URLs by user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url Url
		var expiryAt sql.NullTime

		if err := rows.Scan(&url.ID, &url.UserID, &url.Url, &url.Code, &url.Status, &url.CreatedAt, &expiryAt); err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}

		if expiryAt.Valid {
			url.ExpiryAt = expiryAt.Time
		} else {
			url.ExpiryAt = time.Time{}
		}

		urls = append(urls, UrlWithShortCode{Url: url, ShortUrl: utils.GetShortUrl(url.Code)})
	}

	logStr := fmt.Sprintf("Get URLs by user from DB : %s, UserID: %d, Timestamp: %s", query, userID, time.Now().UTC())
	utils.Log.Info(logStr)

	return urls, nil
}

func (u *CreateShortUrl) Validate() (Url, error) {
	var url Url
	u.Code = utils.GenerateSlug(u.Code, 20)

	urlByCode, urlByCodeErr := getUrlByCode(u.Code)
	if urlByCode.Code == u.Code || urlByCodeErr == nil {
		return url, fmt.Errorf("URL code already exists !")
	}

	return Url{
		UserID:    u.UserID,
		Url:       u.Url,
		Code:      u.Code,
		Status:    UrlStatusActive,
		CreatedAt: time.Now(),
		ExpiryAt:  u.ExpiryAt,
	}, nil
}

func (u *Url) Save() error {

	query := `INSERT INTO url (user_id, url, code, status, created_at, expiry_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	logStr := fmt.Sprintf("Save URL in DB : %s, Code: %s, Timestamp: %s", query, u.Code, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, u.UserID, u.Url, u.Code, u.Status, u.CreatedAt, u.ExpiryAt).Scan(&u.ID, &u.CreatedAt)
	if rowErr != nil {
		errStr := fmt.Sprintf("Error while trying to save URL - %s !", rowErr.Error())
		return fmt.Errorf(errStr)
	}

	return nil
}

func getUrlByCode(code string) (Url, error) {
	var url Url
	query := `SELECT id, user_id, url, code, status, created_at, expiry_at FROM url WHERE code=$1 AND status=$2`

	logStr := fmt.Sprintf("Get URL by code from DB : %s, Code: %s, Timestamp: %s", query, code, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, code, UrlStatusActive).Scan(&url.ID, &url.UserID, &url.Url, &url.Code, &url.Status, &url.CreatedAt, &url.ExpiryAt)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return url, fmt.Errorf("no active URL found for the provided code")
		}
		errStr := fmt.Sprintf("Error while trying to get URL by code - %s !", rowErr.Error())
		return url, fmt.Errorf(errStr)
	}

	return url, nil
}
