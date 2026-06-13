package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/middleware"
	"github.com/marwans200/cyberguard/internal/models"
	"github.com/marwans200/cyberguard/internal/routes"
)

func setupTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	database.DB = db

	middleware.JWTSecret = []byte("test_secret_key_12345678901234567890")

	r := gin.New()
	routes.Setup(r)
	return r, mock
}

func TestCyberGuardEndpoints(t *testing.T) {
	r, mock := setupTestRouter(t)

	adminID := uuid.New().String()
	incidentID := uuid.New().String()
	noteID := uuid.New().String()

	// 1. MOCK REGISTER USER
	// Step A: FindByEmail() check
	mock.ExpectQuery("SELECT id, name, email, password_hash, role, created_at, updated_at.*FROM users WHERE email =").
		WithArgs("admin@cyberguard.local").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "created_at", "updated_at"})) // empty result (no duplicate)

	// Step B: CountAll() check (to determine if first user should be admin)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Step C: Create() insertion
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("Admin User", "admin@cyberguard.local", sqlmock.AnyArg(), "admin").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(adminID, "Admin User", "admin@cyberguard.local", "hashed_password", "admin", time.Now(), time.Now()))

	// Step D: AuditLog insertion
	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(adminID, "User registered", "user", adminID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute Register Request
	regReq := models.RegisterRequest{
		Name:     "Admin User",
		Email:    "admin@cyberguard.local",
		Password: "password123",
		Role:     "admin",
	}
	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Register Response]: %s", w.Body.String())

	// 2. MOCK LOGIN USER
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	mock.ExpectQuery("SELECT id, name, email, password_hash, role, created_at, updated_at.*FROM users WHERE email =").
		WithArgs("admin@cyberguard.local").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(adminID, "Admin User", "admin@cyberguard.local", string(hashedPassword), "admin", time.Now(), time.Now()))

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(adminID, "User logged in", "user", adminID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	loginReq := models.LoginRequest{
		Email:    "admin@cyberguard.local",
		Password: "password123",
	}
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Login Response]: %s", w.Body.String())

	var loginRes models.LoginResponse
	_ = json.Unmarshal(w.Body.Bytes(), &loginRes)
	token := loginRes.Token

	// 3. MOCK GET USER PROFILE
	mock.ExpectQuery("SELECT id, name, email, password_hash, role, created_at, updated_at.*FROM users WHERE id =").
		WithArgs(adminID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(adminID, "Admin User", "admin@cyberguard.local", string(hashedPassword), "admin", time.Now(), time.Now()))

	req = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Profile Response]: %s", w.Body.String())

	// 4. MOCK CREATE INCIDENT
	mock.ExpectQuery("INSERT INTO incidents").
		WithArgs("Brute Force Detected", "150 failed login attempts from IP 185.76.43.22", "Brute Force", "High", "Open", adminID, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "category", "severity", "status", "created_by", "assigned_to", "created_at", "updated_at"}).
			AddRow(incidentID, "Brute Force Detected", "150 failed login attempts from IP 185.76.43.22", "Brute Force", "High", "Open", adminID, nil, time.Now(), time.Now()))

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(adminID, "Created incident: Brute Force Detected", "incident", incidentID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	incReq := models.CreateIncidentRequest{
		Title:       "Brute Force Detected",
		Description: "150 failed login attempts from IP 185.76.43.22",
		Category:    "Brute Force",
		Severity:    "High",
		AssignedTo:  nil,
	}
	body, _ = json.Marshal(incReq)
	req = httptest.NewRequest(http.MethodPost, "/api/incidents", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Create Incident Response]: %s", w.Body.String())

	// 5. MOCK ADD INVESTIGATION NOTE
	// Verify incident exists first
	mock.ExpectQuery("SELECT i.id, i.title.*FROM incidents i.*WHERE i.id =").
		WithArgs(incidentID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "category", "severity", "status", "created_by", "assigned_to", "created_at", "updated_at", "created_by_name", "assigned_to_name"}).
			AddRow(incidentID, "Brute Force Detected", "150 failed login attempts from IP 185.76.43.22", "Brute Force", "High", "Open", adminID, nil, time.Now(), time.Now(), "Admin User", nil))

	// Insert note
	mock.ExpectQuery("INSERT INTO investigation_notes").
		WithArgs(incidentID, adminID, "Blocked the offending IP address in the edge firewall.").
		WillReturnRows(sqlmock.NewRows([]string{"id", "incident_id", "author_id", "content", "created_at"}).
			AddRow(noteID, incidentID, adminID, "Blocked the offending IP address in the edge firewall.", time.Now()))

	// Audit log
	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(adminID, "Added investigation note", "incident", incidentID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	noteReq := models.CreateNoteRequest{
		Content: "Blocked the offending IP address in the edge firewall.",
	}
	body, _ = json.Marshal(noteReq)
	req = httptest.NewRequest(http.MethodPost, "/api/incidents/"+incidentID+"/notes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Create Note Response]: %s", w.Body.String())

	// 6. MOCK DASHBOARD STATS
	mock.ExpectQuery("SELECT status, COUNT\\(\\*\\) FROM incidents GROUP BY status").
		WillReturnRows(sqlmock.NewRows([]string{"status", "count"}).
			AddRow("Open", 1).
			AddRow("Investigating", 0).
			AddRow("Resolved", 0).
			AddRow("Closed", 0))

	req = httptest.NewRequest(http.MethodGet, "/api/dashboard/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d. Response: %s", w.Code, w.Body.String())
	}
	t.Logf("[Dashboard Stats Response]: %s", w.Body.String())

	// Ensure all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("sqlmock expectations were not met: %v", err)
	}
}
