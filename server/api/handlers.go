package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"findarr/server/config"
	"findarr/server/storage"
	"findarr/server/web"
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

	results, err := storage.SearchMedia(strings.ToLower(q))
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		if len(results) == 0 {
			fmt.Fprint(w, "<div class='p-4 text-gray-500'>No results found</div>")
			return
		}
		fmt.Fprint(w, "<div class='space-y-3'>")
		for i, result := range results {
			delay := fmt.Sprintf("style='animation-delay: %dms'", i*50)
			badgeClass := "media-badge media-badge-" + result.Type
			fmt.Fprintf(w, `
			<div class='result-item p-3 border rounded card-findarr' %s>
				<div class='flex justify-between items-start'>
					<div>
						<div class='font-medium'>%s</div>
						<div class='text-sm text-gray-600'>%s</div>
					</div>
					<span class='%s'>%s</span>
				</div>
			</div>`, delay, result.Title, result.Year, badgeClass, result.Type)
		}
		fmt.Fprint(w, "</div>")
	} else {
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
