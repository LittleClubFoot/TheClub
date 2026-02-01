// Package notes provides file-based storage for notes with image/document attachments.
//
// Each note is stored as a directory under the configured data path:
//
//	data/notes/{uuid}/
//	  note.md       — Markdown content with YAML frontmatter
//	  photo.jpg     — Attached files referenced from the markdown
//
// The YAML frontmatter schema:
//
//	---
//	title: "My Note"
//	created: "2024-01-15T10:30:00Z"
//	updated: "2024-01-15T12:00:00Z"
//	tags: ["label", "inventory"]
//	public: false
//	---
package notes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Note represents a single note with its metadata and rendered content.
type Note struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Tags    []string  `json:"tags"`
	Public  bool      `json:"public"`
	Body    string    `json:"-"` // Raw markdown body (without frontmatter)
	HTML    string    `json:"-"` // Rendered HTML
	Files   []string  `json:"-"` // Attached filenames
}

// Store manages notes on the filesystem.
type Store struct {
	dataDir string
	md      goldmark.Markdown
}

// NewStore creates a Store rooted at dataDir.
// It creates the directory if it doesn't exist.
func NewStore(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
	)
	return &Store{dataDir: dataDir, md: md}, nil
}

// Create persists a new note and returns its ID.
func (s *Store) Create(title, body string, tags []string, public bool) (string, error) {
	id := uuid.New().String()[:8] // Short IDs for cleaner URLs
	dir := filepath.Join(s.dataDir, id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create note dir: %w", err)
	}

	now := time.Now().UTC()
	content := buildMarkdown(title, now, now, tags, public, body)

	if err := os.WriteFile(filepath.Join(dir, "note.md"), []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write note: %w", err)
	}
	return id, nil
}

// Get loads a single note by ID.
func (s *Store) Get(id string) (*Note, error) {
	// Sanitize ID to prevent path traversal
	id = filepath.Base(id)
	dir := filepath.Join(s.dataDir, id)

	raw, err := os.ReadFile(filepath.Join(dir, "note.md"))
	if err != nil {
		return nil, fmt.Errorf("read note: %w", err)
	}

	note, err := s.parseNote(id, raw)
	if err != nil {
		return nil, err
	}

	// List attached files
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() && e.Name() != "note.md" {
			note.Files = append(note.Files, e.Name())
		}
	}

	return note, nil
}

// Update overwrites an existing note's content and metadata.
func (s *Store) Update(id, title, body string, tags []string, public bool) error {
	id = filepath.Base(id)
	dir := filepath.Join(s.dataDir, id)

	// Read existing to preserve created time
	existing, err := s.Get(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	content := buildMarkdown(title, existing.Created, now, tags, public, body)

	if err := os.WriteFile(filepath.Join(dir, "note.md"), []byte(content), 0644); err != nil {
		return fmt.Errorf("write note: %w", err)
	}
	return nil
}

// Delete removes a note and all its attachments.
func (s *Store) Delete(id string) error {
	id = filepath.Base(id)
	return os.RemoveAll(filepath.Join(s.dataDir, id))
}

// List returns all notes sorted by creation time (newest first).
func (s *Store) List() ([]*Note, error) {
	entries, err := os.ReadDir(s.dataDir)
	if err != nil {
		return nil, fmt.Errorf("read data dir: %w", err)
	}

	var notes []*Note
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		note, err := s.Get(e.Name())
		if err != nil {
			continue // Skip malformed notes
		}
		notes = append(notes, note)
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Created.After(notes[j].Created)
	})

	return notes, nil
}

// FilePath returns the absolute filesystem path for an attached file.
// Returns an error if the file doesn't exist or the path escapes the note directory.
func (s *Store) FilePath(noteID, filename string) (string, error) {
	noteID = filepath.Base(noteID)
	filename = filepath.Base(filename)
	path := filepath.Join(s.dataDir, noteID, filename)

	// Verify the path is within the note directory
	noteDir := filepath.Join(s.dataDir, noteID)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absDir, err := filepath.Abs(noteDir)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absPath, absDir) {
		return "", fmt.Errorf("path traversal blocked")
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}
	return path, nil
}

// SaveFile writes an uploaded file into the note's directory.
func (s *Store) SaveFile(noteID, filename string, data []byte) error {
	noteID = filepath.Base(noteID)
	filename = filepath.Base(filename)
	dir := filepath.Join(s.dataDir, noteID)

	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	return os.WriteFile(filepath.Join(dir, filename), data, 0644)
}

// parseNote extracts metadata and renders HTML from raw markdown bytes.
func (s *Store) parseNote(id string, raw []byte) (*Note, error) {
	var buf bytes.Buffer
	ctx := parser.NewContext()

	if err := s.md.Convert(raw, &buf, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("parse markdown: %w", err)
	}

	metadata := meta.Get(ctx)

	note := &Note{
		ID:   id,
		HTML: buf.String(),
	}

	// Extract title
	if v, ok := metadata["title"]; ok {
		note.Title = fmt.Sprintf("%v", v)
	}
	if note.Title == "" {
		note.Title = "Untitled"
	}

	// Extract timestamps
	if v, ok := metadata["created"]; ok {
		if t, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", v)); err == nil {
			note.Created = t
		}
	}
	if v, ok := metadata["updated"]; ok {
		if t, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", v)); err == nil {
			note.Updated = t
		}
	}

	// Extract public flag
	if v, ok := metadata["public"]; ok {
		note.Public = v == true
	}

	// Extract tags
	if v, ok := metadata["tags"]; ok {
		if tags, ok := v.([]interface{}); ok {
			for _, t := range tags {
				note.Tags = append(note.Tags, fmt.Sprintf("%v", t))
			}
		}
	}

	// Extract raw body (everything after frontmatter)
	parts := strings.SplitN(string(raw), "---", 3)
	if len(parts) >= 3 {
		note.Body = strings.TrimSpace(parts[2])
	} else {
		note.Body = string(raw)
	}

	return note, nil
}

// buildMarkdown constructs a complete markdown file with YAML frontmatter.
func buildMarkdown(title string, created, updated time.Time, tags []string, public bool, body string) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", title))
	sb.WriteString(fmt.Sprintf("created: %q\n", created.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated: %q\n", updated.Format(time.RFC3339)))

	if len(tags) > 0 {
		sb.WriteString("tags:\n")
		for _, t := range tags {
			t = strings.TrimSpace(t)
			if t != "" {
				sb.WriteString(fmt.Sprintf("  - %q\n", t))
			}
		}
	} else {
		sb.WriteString("tags: []\n")
	}

	sb.WriteString(fmt.Sprintf("public: %t\n", public))
	sb.WriteString("---\n\n")
	sb.WriteString(body)
	sb.WriteString("\n")
	return sb.String()
}
