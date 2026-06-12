package repository

import (
	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/models"
)

// NoteRepo handles all database operations for the investigation_notes table.
type NoteRepo struct{}

// Create inserts a new investigation note.
func (r *NoteRepo) Create(incidentID, authorID, content string) (*models.NoteWithAuthor, error) {
	note := &models.NoteWithAuthor{}
	err := database.DB.QueryRow(
		`INSERT INTO investigation_notes (incident_id, author_id, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, incident_id, author_id, content, created_at`,
		incidentID, authorID, content,
	).Scan(&note.ID, &note.IncidentID, &note.AuthorID, &note.Content, &note.CreatedAt)
	return note, err
}

// FindByIncident returns all notes for an incident, newest first, with author info.
func (r *NoteRepo) FindByIncident(incidentID string) ([]models.NoteWithAuthor, error) {
	rows, err := database.DB.Query(
		`SELECT n.id, n.incident_id, n.author_id, n.content, n.created_at,
		        u.name AS author_name
		 FROM investigation_notes n
		 JOIN users u ON u.id = n.author_id
		 WHERE n.incident_id = $1
		 ORDER BY n.created_at DESC`,
		incidentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.NoteWithAuthor
	for rows.Next() {
		var note models.NoteWithAuthor
		if err := rows.Scan(&note.ID, &note.IncidentID, &note.AuthorID, &note.Content, &note.CreatedAt, &note.AuthorName); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}
