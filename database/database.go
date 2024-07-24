package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var db *sql.DB

func connection() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		log.Fatalf("Database configuration variables are not set")
	}

	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/mysql")

	if err != nil {
		panic(err)
	}

	defer db.Close()
}

func GetDB() *sql.DB {
	return db
}
