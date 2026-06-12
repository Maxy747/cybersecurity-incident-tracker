package repository

import (
	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/models"
)

// AuditRepo handles all database operations for the audit_logs table.
type AuditRepo struct{}

// Log inserts a new audit log entry. resourceType and resourceID are optional.
func (r *AuditRepo) Log(userID, action, resourceType string, resourceID *string) error {
	_, err := database.DB.Exec(
		`INSERT INTO audit_logs (user_id, action, resource_type, resource_id)
		 VALUES ($1, $2, $3, $4)`,
		userID, action, resourceType, resourceID,
	)
	return err
}

// FindAll returns all audit logs enriched with user info, newest first.
func (r *AuditRepo) FindAll(limit, offset int) ([]models.AuditLogWithUser, error) {
	rows, err := database.DB.Query(
		`SELECT al.id, al.user_id, al.action, al.resource_type, al.resource_id, al.timestamp,
		        u.name, u.email
		 FROM audit_logs al
		 JOIN users u ON u.id = al.user_id
		 ORDER BY al.timestamp DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLogWithUser
	for rows.Next() {
		var l models.AuditLogWithUser
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.ResourceType, &l.ResourceID, &l.Timestamp,
			&l.UserName, &l.UserEmail,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// Count returns the total number of audit log entries.
func (r *AuditRepo) Count() (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM audit_logs`).Scan(&count)
	return count, err
}
