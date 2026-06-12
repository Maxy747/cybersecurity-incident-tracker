package models

import "time"

// InvestigationNote is a single update posted by an analyst on an incident.
type InvestigationNote struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// NoteWithAuthor enriches a note with the author's name.
type NoteWithAuthor struct {
	InvestigationNote
	AuthorName string `json:"author_name"`
}

// CreateNoteRequest is the body for POST /api/incidents/:id/notes.
type CreateNoteRequest struct {
	Content string `json:"content" binding:"required,min=5"`
}
