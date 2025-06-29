package api

import (
	"encoding/json"
	"net/http"
)

// Router returns the API router
func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", searchHandler)
	return mux
}

// Dummy search handler
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results := []map[string]string{
		{"title": "Example Movie", "type": "movie", "query": q},
		{"title": "Example Book", "type": "book", "query": q},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

