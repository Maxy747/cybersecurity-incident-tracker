package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marwans200/cyberguard/internal/models"
	"github.com/marwans200/cyberguard/internal/repository"
)

// NoteHandler groups all investigation-note endpoints.
type NoteHandler struct {
	notes     *repository.NoteRepo
	incidents *repository.IncidentRepo
	audit     *repository.AuditRepo
}

// NewNoteHandler constructs a NoteHandler.
func NewNoteHandler(notes *repository.NoteRepo, incidents *repository.IncidentRepo, audit *repository.AuditRepo) *NoteHandler {
	return &NoteHandler{notes: notes, incidents: incidents, audit: audit}
}

// Create godoc
// @Summary      Add an investigation note to an incident
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                   true "Incident UUID"
// @Param        body body models.CreateNoteRequest true "Note content"
// @Success      201  {object} models.NoteWithAuthor
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Router       /api/incidents/{id}/notes [post]
func (h *NoteHandler) Create(c *gin.Context) {
	incidentID := c.Param("id")

	// Verify incident exists
	inc, err := h.incidents.FindByID(incidentID)
	if err != nil || inc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	note, err := h.notes.Create(incidentID, userID.(string), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save note"})
		return
	}

	_ = h.audit.Log(userID.(string),
		"Added investigation note",
		"incident", &incidentID)

	c.JSON(http.StatusCreated, note)
}

// List godoc
// @Summary      List all investigation notes for an incident
// @Tags         notes
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident UUID"
// @Success      200 {array}  models.NoteWithAuthor
// @Failure      404 {object} map[string]string
// @Router       /api/incidents/{id}/notes [get]
func (h *NoteHandler) List(c *gin.Context) {
	incidentID := c.Param("id")

	inc, err := h.incidents.FindByID(incidentID)
	if err != nil || inc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	notes, err := h.notes.FindByIncident(incidentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch notes"})
		return
	}

	if notes == nil {
		notes = []models.NoteWithAuthor{}
	}
	c.JSON(http.StatusOK, gin.H{"data": notes, "count": len(notes)})
}
