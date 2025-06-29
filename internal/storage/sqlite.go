package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
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

	initSchemaAndSeed()
}

// MediaItem represents a row in the media table
type MediaItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Year  string `json:"year"`
}

// SearchMedia returns media items matching the query (case-insensitive in title)
func SearchMedia(query string) ([]MediaItem, error) {
	rows, err := DB.Query("SELECT id, title, type, year FROM media WHERE LOWER(title) LIKE ? ORDER BY id", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MediaItem
	for rows.Next() {
		var m MediaItem
		if err := rows.Scan(&m.ID, &m.Title, &m.Type, &m.Year); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, nil
}

// initSchemaAndSeed creates the media table and seeds initial data if empty
func initSchemaAndSeed() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS media (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			type TEXT NOT NULL,
			year TEXT
		)
	`)
	if err != nil {
		log.Fatalf("failed to create media table: %v", err)
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM media").Scan(&count)
	if err != nil {
		log.Fatalf("failed to count media rows: %v", err)
	}
	if count == 0 {
		seedData := []struct{ title, typ, year string }{
			{"Inception", "movie", "2010"},
			{"The Shining", "movie", "1980"},
			{"The Lord of the Rings", "book", "1954"},
			{"Dune", "book", "1965"},
			{"Dark Side of the Moon", "music", "1973"},
			{"Stranger Things", "show", "2016"},
		}
		for _, s := range seedData {
			_, err := DB.Exec("INSERT INTO media (title, type, year) VALUES (?, ?, ?)", s.title, s.typ, s.year)
			if err != nil {
				log.Printf("failed to seed row (%v): %v", s, err)
			}
		}
	}
}
