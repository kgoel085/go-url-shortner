package model

import (
	"fmt"
	"time"

	"example.com/url-shortner/db"
	"example.com/url-shortner/utils"
)

type Analytics struct {
	ID         int64  `json:"id"`
	UrlID      int64  `json:"url_id" binding:"required"`
	ClickCount int64  `json:"click_count" binding:"required"`
	IPAddress  string `json:"ip_address" binding:"required,ip"`
	UserAgent  string `json:"user_agent" binding:"required"`
	Referrer   string `json:"referrer"`
	CreatedAt  string `json:"created_at"`
}

func (a *Analytics) Save() error {
	query := `INSERT INTO analytics (url_id, ip_address, user_agent, referrer, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	logStr := fmt.Sprintf("Save analytics in DB : %s, URL ID: %d, Timestamp: %s", query, a.UrlID, time.Now().UTC())
	utils.Log.Info(logStr)

	rowErr := db.DB.QueryRow(query, a.UrlID, a.IPAddress, a.UserAgent, a.Referrer, time.Now().UTC()).Scan(&a.ID, &a.CreatedAt)

	if rowErr != nil {
		errStr := fmt.Sprintf("Error while trying to save analytics - %s !", rowErr.Error())
		return fmt.Errorf(errStr)
	}

	// Increment click count in url table
	updateQuery := `UPDATE url SET click_count = click_count + 1 WHERE id = $1`
	_, updateErr := db.DB.Exec(updateQuery, a.UrlID)
	if updateErr != nil {
		errStr := fmt.Sprintf("Error while trying to update click count - %s !", updateErr.Error())
		return fmt.Errorf(errStr)
	}

	return nil
}
