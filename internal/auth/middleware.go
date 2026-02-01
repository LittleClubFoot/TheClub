// Package auth provides authentication middleware for The Club.
//
// It supports two modes:
//  1. API token — A shared secret passed via Authorization: Bearer <token> header
//     or via ?token=<token> query parameter. Configured via CLUB_AUTH_TOKEN env var.
//  2. Authelia forward-auth — When deployed behind Caddy with Authelia,
//     the Remote-User / Remote-Groups headers are trusted. This is the
//     recommended production setup.
//
// If CLUB_AUTH_TOKEN is empty, all write operations are denied (safe default).
// Read-only access to public notes is always allowed without auth.
package auth

import (
	"log"
	"net/http"
	"os"
)

// Middleware wraps an http.Handler with authentication checks.
type Middleware struct {
	token string
}

// New creates auth middleware. It reads the token from CLUB_AUTH_TOKEN env var.
func New() *Middleware {
	token := os.Getenv("CLUB_AUTH_TOKEN")
	if token == "" {
		log.Println("WARNING: CLUB_AUTH_TOKEN not set — protected endpoints will reject all requests")
	}
	return &Middleware{token: token}
}

// RequireAuth wraps a handler to require authentication.
// Authentication succeeds if any of these are true:
//   - Authorization: Bearer <token> matches CLUB_AUTH_TOKEN
//   - ?token=<token> query parameter matches CLUB_AUTH_TOKEN
//   - Remote-User header is present (set by Authelia forward-auth via Caddy)
func (m *Middleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check Authelia forward-auth header (trusted when behind Caddy)
		if r.Header.Get("Remote-User") != "" {
			next(w, r)
			return
		}

		// Check Bearer token
		if m.token != "" {
			auth := r.Header.Get("Authorization")
			if auth == "Bearer "+m.token {
				next(w, r)
				return
			}

			// Check query parameter token (for browser-based access)
			if r.URL.Query().Get("token") == m.token {
				next(w, r)
				return
			}
		}

		// Check session cookie (set after first successful token auth)
		cookie, err := r.Cookie("club_session")
		if err == nil && cookie.Value == m.token && m.token != "" {
			next(w, r)
			return
		}

		// No valid auth found — show login page for browsers, 401 for API
		if r.Header.Get("Accept") == "application/json" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Redirect to login
		http.Redirect(w, r, "/notes/login?redirect="+r.URL.RequestURI(), http.StatusSeeOther)
	}
}
