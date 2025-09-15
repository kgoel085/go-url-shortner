package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"kgoel085.com/url-shortner/config"
	"kgoel085.com/url-shortner/utils"
)

var DB *sql.DB

func InitDB() {
	config := config.Config.DB
	dbUrl := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%v sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
		config.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		errStr := err.Error()
		println("Error opening database: " + errStr)
		panic(errStr)
	}
	fmt.Println("DB PINGED: ", DB.Ping())

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUserTable()
	createOtpTable()
	createUrlTable()
	createAnalyticsTable()
}

func createAnalyticsTable() {
	createAnalyticsTable := `
	CREATE TABLE IF NOT EXISTS analytics (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		url_id BIGINT NOT NULL,
		ip_address TEXT NOT NULL,
		user_agent TEXT NOT NULL,
		referrer TEXT,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (url_id) REFERENCES url(id)
	);`

	_, err := DB.Exec(createAnalyticsTable)
	if err != nil {
		errStr := fmt.Sprintf("Error creating analytics table: %v", err)
		utils.Log.Error(errStr)
		panic(errStr)
	} else {
		utils.Log.Info("Table `analytics` created or already exists")
	}
}

func createOtpTable() {
	createOtpTable := `
	DO $$
	BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_type') THEN
					CREATE TYPE otp_type AS ENUM ('email', 'phone');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_action_type') THEN
					CREATE TYPE otp_action_type AS ENUM ('login', 'signup', 'reset_password');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_status') THEN
					CREATE TYPE otp_status AS ENUM ('pending', 'success', 'expire');
			END IF;
	END$$;

	CREATE TABLE IF NOT EXISTS otp (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		key TEXT NOT NULL,
		type otp_type NOT NULL,
		action otp_action_type NOT NULL,
		otp TEXT NOT NULL,
		status otp_status NOT NULL DEFAULT 'pending',
		token UUID NOT NULL DEFAULT gen_random_uuid(),
		created_at TIMESTAMP NOT NULL
	);`

	_, err := DB.Exec(createOtpTable)
	if err != nil {
		errStr := fmt.Sprintf("Error creating otp table: %v", err)
		utils.Log.Error(errStr)
		panic(errStr)
	} else {
		utils.Log.Info("Table `otp` created or already exists")
	}
}

func createUrlTable() {
	createUrlTable := `
	DO $$
	BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'url_status') THEN
					CREATE TYPE url_status AS ENUM ('active', 'inactive', 'deleted', 'expired');
			END IF;
	END$$;

	CREATE TABLE IF NOT EXISTS url (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		user_id BIGINT NOT NULL,
		url TEXT NOT NULL,
		code TEXT NOT NULL UNIQUE,
		status url_status NOT NULL DEFAULT 'active',
		created_at TIMESTAMP NOT NULL,
		expiry_at TIMESTAMP,
		click_count BIGINT NOT NULL DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(createUrlTable)
	if err != nil {
		errStr := fmt.Sprintf("Error creating urls table: %v", err)
		utils.Log.Error(errStr)
		panic(errStr)
	} else {
		utils.Log.Info("Table `urls` created or already exists")
	}
}

func createUserTable() {
	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		errStr := fmt.Sprintf("Error creating users table: %v", err)
		utils.Log.Error(errStr)
		panic(errStr)
	} else {
		utils.Log.Info("Table `users` created or already exists")
	}
}
