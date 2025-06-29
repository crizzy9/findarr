package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// DB holds the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection
func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
}
