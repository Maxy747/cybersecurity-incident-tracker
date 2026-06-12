package models

import "time"

// AuditLog records every significant user action in the system.
type AuditLog struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type,omitempty"`
	ResourceID   *string   `json:"resource_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// AuditLogWithUser enriches an audit log entry with the actor's name and email.
type AuditLogWithUser struct {
	AuditLog
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
