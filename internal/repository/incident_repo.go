package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/models"
)

// IncidentRepo handles all database operations for the incidents table.
type IncidentRepo struct{}

// Create inserts a new incident record.
func (r *IncidentRepo) Create(req *models.CreateIncidentRequest, createdBy string) (*models.Incident, error) {
	inc := &models.Incident{}
	err := database.DB.QueryRow(
		`INSERT INTO incidents (title, description, category, severity, status, created_by, assigned_to)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, title, description, category, severity, status, created_by, assigned_to, created_at, updated_at`,
		req.Title, req.Description, req.Category, req.Severity, models.StatusOpen, createdBy, req.AssignedTo,
	).Scan(&inc.ID, &inc.Title, &inc.Description, &inc.Category, &inc.Severity, &inc.Status,
		&inc.CreatedBy, &inc.AssignedTo, &inc.CreatedAt, &inc.UpdatedAt)
	return inc, err
}

// FindAll returns all incidents optionally filtered.
func (r *IncidentRepo) FindAll(f models.IncidentFilter) ([]models.IncidentWithNames, error) {
	base := `
		SELECT i.id, i.title, i.description, i.category, i.severity, i.status,
		       i.created_by, i.assigned_to, i.created_at, i.updated_at,
		       cb.name AS created_by_name,
		       at.name AS assigned_to_name
		FROM incidents i
		JOIN users cb ON cb.id = i.created_by
		LEFT JOIN users at ON at.id = i.assigned_to
		WHERE 1=1`

	var args []interface{}
	idx := 1

	if f.Title != "" {
		base += fmt.Sprintf(" AND LOWER(i.title) LIKE LOWER($%d)", idx)
		args = append(args, "%"+f.Title+"%")
		idx++
	}
	if f.Category != "" {
		base += fmt.Sprintf(" AND i.category = $%d", idx)
		args = append(args, f.Category)
		idx++
	}
	if f.Severity != "" {
		base += fmt.Sprintf(" AND i.severity = $%d", idx)
		args = append(args, f.Severity)
		idx++
	}
	if f.Status != "" {
		base += fmt.Sprintf(" AND i.status = $%d", idx)
		args = append(args, f.Status)
		idx++
	}
	if f.AssignedTo != "" {
		base += fmt.Sprintf(" AND i.assigned_to = $%d", idx)
		args = append(args, f.AssignedTo)
		idx++
	}
	_ = strings.TrimSpace(base) // lint

	base += " ORDER BY i.created_at DESC"

	rows, err := database.DB.Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []models.IncidentWithNames
	for rows.Next() {
		var inc models.IncidentWithNames
		var assignedTo sql.NullString
		var assignedToName sql.NullString
		if err := rows.Scan(
			&inc.ID, &inc.Title, &inc.Description, &inc.Category, &inc.Severity, &inc.Status,
			&inc.CreatedBy, &assignedTo, &inc.CreatedAt, &inc.UpdatedAt,
			&inc.CreatedByName, &assignedToName,
		); err != nil {
			return nil, err
		}
		if assignedTo.Valid {
			inc.AssignedTo = &assignedTo.String
		}
		if assignedToName.Valid {
			inc.AssignedToName = &assignedToName.String
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

// FindByID returns a single incident with creator/assignee names.
func (r *IncidentRepo) FindByID(id string) (*models.IncidentWithNames, error) {
	inc := &models.IncidentWithNames{}
	var assignedTo sql.NullString
	var assignedToName sql.NullString

	err := database.DB.QueryRow(`
		SELECT i.id, i.title, i.description, i.category, i.severity, i.status,
		       i.created_by, i.assigned_to, i.created_at, i.updated_at,
		       cb.name AS created_by_name,
		       at.name AS assigned_to_name
		FROM incidents i
		JOIN users cb ON cb.id = i.created_by
		LEFT JOIN users at ON at.id = i.assigned_to
		WHERE i.id = $1`, id,
	).Scan(
		&inc.ID, &inc.Title, &inc.Description, &inc.Category, &inc.Severity, &inc.Status,
		&inc.CreatedBy, &assignedTo, &inc.CreatedAt, &inc.UpdatedAt,
		&inc.CreatedByName, &assignedToName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if assignedTo.Valid {
		inc.AssignedTo = &assignedTo.String
	}
	if assignedToName.Valid {
		inc.AssignedToName = &assignedToName.String
	}
	return inc, nil
}

// Update applies partial updates to an incident.
func (r *IncidentRepo) Update(id string, req *models.UpdateIncidentRequest) (*models.Incident, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", idx))
		args = append(args, *req.Title)
		idx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}
	if req.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", idx))
		args = append(args, *req.Category)
		idx++
	}
	if req.Severity != nil {
		setClauses = append(setClauses, fmt.Sprintf("severity = $%d", idx))
		args = append(args, *req.Severity)
		idx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, *req.Status)
		idx++
	}
	if req.AssignedTo != nil {
		setClauses = append(setClauses, fmt.Sprintf("assigned_to = $%d", idx))
		args = append(args, *req.AssignedTo)
		idx++
	}

	args = append(args, id)
	query := fmt.Sprintf(
		`UPDATE incidents SET %s WHERE id = $%d
		 RETURNING id, title, description, category, severity, status, created_by, assigned_to, created_at, updated_at`,
		strings.Join(setClauses, ", "), idx,
	)

	inc := &models.Incident{}
	err := database.DB.QueryRow(query, args...).
		Scan(&inc.ID, &inc.Title, &inc.Description, &inc.Category, &inc.Severity, &inc.Status,
			&inc.CreatedBy, &inc.AssignedTo, &inc.CreatedAt, &inc.UpdatedAt)
	return inc, err
}

// Delete removes an incident by ID.
func (r *IncidentRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM incidents WHERE id = $1`, id)
	return err
}

// --- Dashboard queries ---

// CountByStatus returns total, open, investigating, resolved, and closed counts.
func (r *IncidentRepo) CountByStatus() (map[string]int, error) {
	rows, err := database.DB.Query(`SELECT status, COUNT(*) FROM incidents GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{"total": 0, "Open": 0, "Investigating": 0, "Resolved": 0, "Closed": 0}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
		counts["total"] += count
	}
	return counts, rows.Err()
}

// CountBySeverity returns incident counts grouped by severity level.
func (r *IncidentRepo) CountBySeverity() (map[string]int, error) {
	rows, err := database.DB.Query(`SELECT severity, COUNT(*) FROM incidents GROUP BY severity`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			return nil, err
		}
		counts[severity] = count
	}
	return counts, rows.Err()
}

// TrendsByDay returns incident counts per day for the last N days.
func (r *IncidentRepo) TrendsByDay(days int) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		`SELECT DATE(created_at) AS day, COUNT(*) AS count
		 FROM incidents
		 WHERE created_at >= NOW() - ($1 || ' days')::INTERVAL
		 GROUP BY day ORDER BY day`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trend []map[string]interface{}
	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		trend = append(trend, map[string]interface{}{"day": day, "count": count})
	}
	return trend, rows.Err()
}
