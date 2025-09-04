package db

import (
	"database/sql"
	"fmt"

	"example.com/url-shortner/config"
	_ "github.com/lib/pq"
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
	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		panic(fmt.Sprintf("Error creating users table: %v", err))
	} else {
		fmt.Println("Table `users` created or already exists")
	}
}
