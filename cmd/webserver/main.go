// Package main implements The Club's Go/HTMX test server.
//
// This server provides custom applications and dynamic content for The Club dashboard.
// It demonstrates HTMX integration and serves as a foundation for additional features.
//
// The server is designed to run behind Caddy reverse proxy in production,
// but can be run standalone for development and testing.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// main initializes and starts the HTTP server with all routes configured.
func main() {
	// Load HTML templates from the templates directory
	// This will panic if templates are malformed, which is desired behavior
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Create HTTP multiplexer for routing
	mux := http.NewServeMux()
	
	// ========================================================================
	// ROUTE HANDLERS
	// ========================================================================
	
	// Root route handler
	// In production, Caddy handles routing to Homer dashboard
	// In development, this redirects to Homer for testing
	mux.HandleFunc("/", handleRoot)
	
	// Main test server page - demonstrates HTMX integration
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		handleTestPage(w, r, templates)
	})
	
	// Dynamic time endpoint for HTMX demonstration
	// Returns HTML fragment showing current server time
	mux.HandleFunc("/test/time", handleTimeEndpoint)
	
	// API documentation page
	mux.HandleFunc("/test/api", func(w http.ResponseWriter, r *http.Request) {
		handleAPIDocsPage(w, r, templates)
	})
	
	// General documentation page
	mux.HandleFunc("/test/docs", func(w http.ResponseWriter, r *http.Request) {
		handleDocsPage(w, r, templates)
	})

	// ========================================================================
	// SERVER STARTUP
	// ========================================================================
	
	// Log server startup information
	fmt.Println("🚀 Starting The Club Go/HTMX Test Server")
	fmt.Println("📍 Listening on port :8080")
	fmt.Println("🌐 Available endpoints:")
	fmt.Println("   • /test      - Main test page with HTMX demo")
	fmt.Println("   • /test/api  - API documentation")
	fmt.Println("   • /test/docs - General documentation")
	fmt.Println("   • /test/time - Dynamic time endpoint (HTMX)")
	
	// Start HTTP server
	// This will block until the server shuts down or encounters an error
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

// ============================================================================
// HANDLER FUNCTIONS
// ============================================================================

// handleRoot handles requests to the root path.
// In development, redirects to Homer dashboard for testing.
// In production, Caddy routes root requests directly to Homer.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// Redirect to Homer dashboard (development mode)
		// In production, Caddy handles this routing
		http.Redirect(w, r, "http://localhost:8080", http.StatusTemporaryRedirect)
		return
	}
	// Return 404 for any other root-level paths
	http.NotFound(w, r)
}

// handleTestPage renders the main test page with HTMX demonstration.
// This page showcases dynamic content loading and serves as a template
// for building additional HTMX-powered features.
func handleTestPage(w http.ResponseWriter, r *http.Request, templates *template.Template) {
	// Prepare template data
	data := struct {
		Title     string // Page title for HTML head
		Message   string // Welcome message
		Timestamp string // Server startup timestamp
	}{
		Title:     "Go/HTMX Test Server",
		Message:   "Welcome to The Club Test Server!",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	// Set content type for proper rendering
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Render template with data
	if err := templates.ExecuteTemplate(w, "test.html", data); err != nil {
		log.Printf("❌ Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleTimeEndpoint returns the current server time as an HTML fragment.
// This endpoint is designed for HTMX requests and demonstrates
// dynamic content updates without full page reloads.
func handleTimeEndpoint(w http.ResponseWriter, r *http.Request) {
	// Set content type for HTML fragment
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Return formatted time as HTML
	fmt.Fprintf(w, `<div class="time-display">Current time: %s</div>`, 
		time.Now().Format("15:04:05"))
}

// handleAPIDocsPage renders the API documentation page.
// This page documents available endpoints and their usage.
func handleAPIDocsPage(w http.ResponseWriter, r *http.Request, templates *template.Template) {
	data := struct {
		Title string
	}{
		Title: "API Documentation",
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	if err := templates.ExecuteTemplate(w, "api.html", data); err != nil {
		log.Printf("❌ Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleDocsPage renders the general documentation page.
// This page provides information about The Club system and usage.
func handleDocsPage(w http.ResponseWriter, r *http.Request, templates *template.Template) {
	data := struct {
		Title string
	}{
		Title: "Documentation",
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	if err := templates.ExecuteTemplate(w, "docs.html", data); err != nil {
		log.Printf("❌ Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

