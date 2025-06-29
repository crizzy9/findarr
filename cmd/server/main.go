package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yourusername/findarr/internal/config"
)

// Template cache
var templates map[string]*template.Template

// App configuration
var appConfig *config.Config

// Initialize templates and configuration
func init() {
	// Load configuration
	appConfig = config.LoadConfig()

	// Ensure config directory exists
	configDir := filepath.Dir(appConfig.Database.Path)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create config directory: %v", err)
		}
	}

	// Initialize template cache
	templates = make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/index.html",
	))
}

// indexHandler renders the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Config": appConfig,
	}
	
	err := templates["index"].ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// API routes
	r.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		
		// For demonstration, check if client accepts HTML (htmx request)
		if r.Header.Get("HX-Request") == "true" {
			// Return HTML for htmx
			results := []map[string]string{
				{"title": "Inception", "type": "movie", "year": "2010", "id": "tt1375666"},
				{"title": "The Shining", "type": "movie", "year": "1980", "id": "tt0081505"},
				{"title": "The Lord of the Rings", "type": "book", "year": "1954", "id": "978-0618640157"},
				{"title": "Dune", "type": "book", "year": "1965", "id": "978-0441172719"},
				{"title": "Dark Side of the Moon", "type": "music", "year": "1973", "id": "album-1973-pink-floyd"},
				{"title": "Stranger Things", "type": "show", "year": "2016", "id": "tt4574334"},
			}
			
			// Only show results matching the query (case insensitive)
			var filteredResults []map[string]string
			for _, result := range results {
				if strings.Contains(strings.ToLower(result["title"]), strings.ToLower(q)) {
					filteredResults = append(filteredResults, result)
				}
			}
			
			w.Header().Set("Content-Type", "text/html")
			if len(filteredResults) == 0 {
				fmt.Fprint(w, "<div class='p-4 text-gray-500'>No results found</div>")
				return
			}
			
			fmt.Fprint(w, "<div class='space-y-3'>")
			for i, result := range filteredResults {
				// Add animation delay based on index
				delay := fmt.Sprintf("style='animation-delay: %dms'", i*50)
				
				// Create badge class based on media type
				badgeClass := "media-badge media-badge-" + result["type"]
				
				fmt.Fprintf(w, `
				<div class='result-item p-3 border rounded card-findarr' %s>
					<div class='flex justify-between items-start'>
						<div>
							<div class='font-medium'>%s</div>
							<div class='text-sm text-gray-600'>%s</div>
						</div>
						<span class='%s'>%s</span>
					</div>
				</div>`, delay, result["title"], result["year"], badgeClass, result["type"])
			}
			fmt.Fprint(w, "</div>")
		} else {
			// Return JSON for API clients
			results := []map[string]string{
				{"title": "Example Movie", "type": "movie", "query": q},
				{"title": "Example Book", "type": "book", "query": q},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		}
		
		// Set cookie for last search
		http.SetCookie(w, &http.Cookie{Name: "LastSearch", Value: q, Path: "/"})
	})

	// Configuration info endpoint
	r.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(appConfig)
	})

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"status": "ok",
			"version": "0.1.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	// Frontend routes
	r.HandleFunc("/", indexHandler)

	// Start the server
	serverAddr := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	log.Printf("Server started at %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatal(err)
	}
}
