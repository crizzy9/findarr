package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/crizzy9/findarr/internal/config"
	"github.com/crizzy9/findarr/internal/web"
)

// Handlers encapsulates dependencies for HTTP handlers
// Use NewHandlers to create, passing config and template manager
type Handlers struct {
	Config    *config.Config
	Templates *web.TemplateManager
}

// NewHandlers creates a new Handlers instance
func NewHandlers(cfg *config.Config, tm *web.TemplateManager) *Handlers {
	return &Handlers{Config: cfg, Templates: tm}
}

// RegisterRoutes registers all HTTP handlers to the given mux/router
func (h *Handlers) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/api/search", h.SearchHandler)
	r.HandleFunc("/api/config", h.ConfigHandler)
	r.HandleFunc("/health", h.HealthHandler)
	r.HandleFunc("/", h.IndexHandler)
}

// IndexHandler renders the main page
func (h *Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Config": h.Config,
	}

	// Get the index.html template which includes base.html
	// and execute it directly (no need to specify "base.html")
	err := h.Templates.Get("index.html").Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SearchHandler handles /api/search requests
func (h *Handlers) SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	if r.Header.Get("HX-Request") == "true" {
		results := []map[string]string{
			{"title": "Inception", "type": "movie", "year": "2010", "id": "tt1375666"},
			{"title": "The Shining", "type": "movie", "year": "1980", "id": "tt0081505"},
			{"title": "The Lord of the Rings", "type": "book", "year": "1954", "id": "978-0618640157"},
			{"title": "Dune", "type": "book", "year": "1965", "id": "978-0441172719"},
			{"title": "Dark Side of the Moon", "type": "music", "year": "1973", "id": "album-1973-pink-floyd"},
			{"title": "Stranger Things", "type": "show", "year": "2016", "id": "tt4574334"},
		}
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
			delay := fmt.Sprintf("style='animation-delay: %dms'", i*50)
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
		results := []map[string]string{
			{"title": "Example Movie", "type": "movie", "query": q},
			{"title": "Example Book", "type": "book", "query": q},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
	http.SetCookie(w, &http.Cookie{Name: "LastSearch", Value: q, Path: "/"})
}

// ConfigHandler returns application config as JSON
func (h *Handlers) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Config)
}

// HealthHandler returns basic health info
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"status":  "ok",
		"version": "0.1.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
