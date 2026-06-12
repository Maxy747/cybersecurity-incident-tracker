package models

import "time"

// Incident category constants
const (
	CategoryMalware           = "Malware"
	CategoryPhishing          = "Phishing"
	CategoryDDoS              = "DDoS"
	CategoryUnauthorizedAccess = "Unauthorized Access"
	CategoryDataBreach        = "Data Breach"
	CategoryInsiderThreat     = "Insider Threat"
	CategoryBruteForce        = "Brute Force"
	CategoryNetworkAttack     = "Network Attack"
)

// Severity level constants
const (
	SeverityCritical      = "Critical"
	SeverityHigh          = "High"
	SeverityMedium        = "Medium"
	SeverityLow           = "Low"
	SeverityInformational = "Informational"
)

// Status constants
const (
	StatusOpen         = "Open"
	StatusInvestigating = "Investigating"
	StatusResolved     = "Resolved"
	StatusClosed       = "Closed"
)

// Incident is the core entity representing a cybersecurity incident.
type Incident struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	AssignedTo  *string   `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IncidentWithNames enriches Incident with human-readable user names.
type IncidentWithNames struct {
	Incident
	CreatedByName  string  `json:"created_by_name"`
	AssignedToName *string `json:"assigned_to_name,omitempty"`
}

// CreateIncidentRequest is the body for POST /api/incidents.
type CreateIncidentRequest struct {
	Title       string  `json:"title"       binding:"required,min=3,max=255"`
	Description string  `json:"description" binding:"required"`
	Category    string  `json:"category"    binding:"required,oneof='Malware' 'Phishing' 'DDoS' 'Unauthorized Access' 'Data Breach' 'Insider Threat' 'Brute Force' 'Network Attack'"`
	Severity    string  `json:"severity"    binding:"required,oneof=Critical High Medium Low Informational"`
	AssignedTo  *string `json:"assigned_to"`
}

// UpdateIncidentRequest is the body for PUT /api/incidents/:id.
type UpdateIncidentRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Severity    *string `json:"severity"    binding:"omitempty,oneof=Critical High Medium Low Informational"`
	Status      *string `json:"status"      binding:"omitempty,oneof=Open Investigating Resolved Closed"`
	AssignedTo  *string `json:"assigned_to"`
}

// IncidentFilter holds query parameters for filtering the incident list.
type IncidentFilter struct {
	Title      string
	Category   string
	Severity   string
	Status     string
	AssignedTo string
}
