package notes

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/LittleClubFoot/TheClub/internal/auth"
)

// Handler holds dependencies for note HTTP handlers.
type Handler struct {
	store     *Store
	auth      *auth.Middleware
	templates *template.Template
	baseURL   string // e.g. "https://pochita.synology.me" for QR codes
}

// NewHandler creates a Handler with the given dependencies.
func NewHandler(store *Store, authMW *auth.Middleware, templates *template.Template) *Handler {
	baseURL := os.Getenv("CLUB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://localhost"
	}
	return &Handler{
		store:     store,
		auth:      authMW,
		templates: templates,
		baseURL:   strings.TrimRight(baseURL, "/"),
	}
}

// RegisterRoutes adds all note routes to the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("/notes/login", h.handleLogin)
	mux.HandleFunc("/notes/login/submit", h.handleLoginSubmit)

	// Protected routes (require auth)
	mux.HandleFunc("/notes", h.auth.RequireAuth(h.handleList))
	mux.HandleFunc("/notes/new", h.auth.RequireAuth(h.handleNewForm))
	mux.HandleFunc("/notes/create", h.auth.RequireAuth(h.handleCreate))
	mux.HandleFunc("/notes/upload", h.auth.RequireAuth(h.handleUpload))

	// Note-specific routes — auth checked per-note (public notes bypass auth)
	mux.HandleFunc("/notes/view/", h.handleView)
	mux.HandleFunc("/notes/edit/", h.auth.RequireAuth(h.handleEditForm))
	mux.HandleFunc("/notes/update/", h.auth.RequireAuth(h.handleUpdate))
	mux.HandleFunc("/notes/delete/", h.auth.RequireAuth(h.handleDelete))
	mux.HandleFunc("/notes/qr/", h.handleQR)
	mux.HandleFunc("/notes/files/", h.handleFile)
}

// handleList renders the list of all notes.
func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.List()
	if err != nil {
		log.Printf("Error listing notes: %v", err)
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Notes   []*Note
		BaseURL string
	}{
		Title:   "Notes",
		Notes:   notes,
		BaseURL: h.baseURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "notes_list.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleNewForm renders the note creation form.
func (h *Handler) handleNewForm(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "New Note",
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "notes_new.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleCreate processes the note creation form.
func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (32 MB max)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	body := r.FormValue("body")
	tagsRaw := r.FormValue("tags")
	public := r.FormValue("public") == "on"

	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	if title == "" {
		title = "Untitled"
	}

	id, err := h.store.Create(title, body, tags, public)
	if err != nil {
		log.Printf("Error creating note: %v", err)
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	// Handle file uploads
	files := r.MultipartForm.File["files"]
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		if err := h.store.SaveFile(id, fh.Filename, data); err != nil {
			log.Printf("Error saving file %s: %v", fh.Filename, err)
		}
	}

	http.Redirect(w, r, "/notes/view/"+id, http.StatusSeeOther)
}

// handleView renders a single note as a clean markdown page.
// Public notes are accessible without auth. Private notes require auth.
func (h *Handler) handleView(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/notes/view/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	note, err := h.store.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// If note is not public, require authentication
	if !note.Public {
		// Check auth inline (same logic as middleware)
		if !h.isAuthenticated(r) {
			http.Redirect(w, r, "/notes/login?redirect="+r.URL.RequestURI(), http.StatusSeeOther)
			return
		}
	}

	data := struct {
		Title   string
		Note    *Note
		BaseURL string
	}{
		Title:   note.Title,
		Note:    note,
		BaseURL: h.baseURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "notes_view.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleEditForm renders the edit form for an existing note.
func (h *Handler) handleEditForm(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/notes/edit/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	note, err := h.store.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
		Note  *Note
	}{
		Title: "Edit: " + note.Title,
		Note:  note,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "notes_edit.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleUpdate processes the edit form submission.
func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/notes/update/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	body := r.FormValue("body")
	tagsRaw := r.FormValue("tags")
	public := r.FormValue("public") == "on"

	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	if err := h.store.Update(id, title, body, tags, public); err != nil {
		log.Printf("Error updating note: %v", err)
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	// Handle new file uploads
	files := r.MultipartForm.File["files"]
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		if err := h.store.SaveFile(id, fh.Filename, data); err != nil {
			log.Printf("Error saving file %s: %v", fh.Filename, err)
		}
	}

	http.Redirect(w, r, "/notes/view/"+id, http.StatusSeeOther)
}

// handleDelete removes a note.
func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/notes/delete/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	if err := h.store.Delete(id); err != nil {
		log.Printf("Error deleting note: %v", err)
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// handleQR generates a QR code PNG for the note's shareable URL.
func (h *Handler) handleQR(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/notes/qr/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	// Verify note exists
	if _, err := h.store.Get(id); err != nil {
		http.NotFound(w, r)
		return
	}

	url := fmt.Sprintf("%s/notes/view/%s", h.baseURL, id)

	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(png)
}

// handleFile serves attached files from a note's directory.
func (h *Handler) handleFile(w http.ResponseWriter, r *http.Request) {
	// Path format: /notes/files/{id}/{filename}
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/notes/files/"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.NotFound(w, r)
		return
	}

	noteID := parts[0]
	filename := parts[1]

	// Check if note is public or user is authenticated
	note, err := h.store.Get(noteID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if !note.Public && !h.isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	path, err := h.store.FilePath(noteID, filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, path)
}

// handleUpload handles HTMX file upload requests and returns a file list fragment.
func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	noteID := r.FormValue("note_id")
	if noteID == "" {
		http.Error(w, "note_id required", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	var uploaded []string
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		if err := h.store.SaveFile(noteID, fh.Filename, data); err != nil {
			log.Printf("Error saving file %s: %v", fh.Filename, err)
			continue
		}
		uploaded = append(uploaded, fh.Filename)
	}

	// Return HTMX fragment with uploaded file list
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, name := range uploaded {
		fmt.Fprintf(w, `<div class="file-item">Uploaded: %s</div>`, template.HTMLEscapeString(name))
	}
}

// handleLogin renders the login page.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/notes"
	}
	data := struct {
		Title    string
		Redirect string
		Error    string
	}{
		Title:    "Login",
		Redirect: redirect,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "notes_login.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleLoginSubmit processes the login form.
func (h *Handler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.FormValue("token")
	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/notes"
	}

	expectedToken := os.Getenv("CLUB_AUTH_TOKEN")
	if expectedToken == "" || token != expectedToken {
		data := struct {
			Title    string
			Redirect string
			Error    string
		}{
			Title:    "Login",
			Redirect: redirect,
			Error:    "Invalid token",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		h.templates.ExecuteTemplate(w, "notes_login.html", data)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "club_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30, // 30 days
	})

	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// isAuthenticated checks if the current request has valid auth.
func (h *Handler) isAuthenticated(r *http.Request) bool {
	// Authelia header
	if r.Header.Get("Remote-User") != "" {
		return true
	}

	token := os.Getenv("CLUB_AUTH_TOKEN")
	if token == "" {
		return false
	}

	// Bearer token
	if r.Header.Get("Authorization") == "Bearer "+token {
		return true
	}

	// Query param
	if r.URL.Query().Get("token") == token {
		return true
	}

	// Session cookie
	cookie, err := r.Cookie("club_session")
	if err == nil && cookie.Value == token {
		return true
	}

	return false
}
