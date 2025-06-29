package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/crizzy9/findarr/internal/config"
	"github.com/crizzy9/findarr/internal/api"
	"github.com/crizzy9/findarr/internal/web"
	"github.com/crizzy9/findarr/internal/storage"
)

func main() {
	// Load configuration
	appConfig := config.LoadConfig()

	// Initialize database
	storage.InitDB(appConfig.Database.Path)

	// Ensure config directory exists
	configDir := filepath.Dir(appConfig.Database.Path)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create config directory: %v", err)
		}
	}

	// Initialize templates
	templateFiles := []string{"web/templates/base.html", "web/templates/index.html"}
	templates := web.NewTemplateManager(templateFiles)
	if err := templates.Load(); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Create handlers and mux
	handlers := api.NewHandlers(appConfig, templates)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	// Serve static files for /static (if needed)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Start the server
	serverAddr := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	log.Printf("Server started at %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal(err)
	}
}

