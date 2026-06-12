package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marwans200/cyberguard/internal/models"
	"github.com/marwans200/cyberguard/internal/repository"
)

// IncidentHandler groups all incident-related endpoints.
type IncidentHandler struct {
	incidents *repository.IncidentRepo
	audit     *repository.AuditRepo
}

// NewIncidentHandler constructs an IncidentHandler.
func NewIncidentHandler(incidents *repository.IncidentRepo, audit *repository.AuditRepo) *IncidentHandler {
	return &IncidentHandler{incidents: incidents, audit: audit}
}

// Create godoc
// @Summary      Create a new incident
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body models.CreateIncidentRequest true "Incident payload"
// @Success      201  {object} models.Incident
// @Failure      400  {object} map[string]string
// @Router       /api/incidents [post]
func (h *IncidentHandler) Create(c *gin.Context) {
	var req models.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	inc, err := h.incidents.Create(&req, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create incident"})
		return
	}

	_ = h.audit.Log(userID.(string),
		fmt.Sprintf("Created incident: %s", inc.Title),
		"incident", &inc.ID)

	c.JSON(http.StatusCreated, inc)
}

// List godoc
// @Summary      List all incidents
// @Description  Supports filtering via query params: title, category, severity, status, assigned_to
// @Tags         incidents
// @Produce      json
// @Security     BearerAuth
// @Param        title       query string false "Filter by title (partial match)"
// @Param        category    query string false "Filter by category"
// @Param        severity    query string false "Filter by severity"
// @Param        status      query string false "Filter by status"
// @Param        assigned_to query string false "Filter by assigned analyst UUID"
// @Success      200 {array}  models.IncidentWithNames
// @Router       /api/incidents [get]
func (h *IncidentHandler) List(c *gin.Context) {
	filter := models.IncidentFilter{
		Title:      c.Query("title"),
		Category:   c.Query("category"),
		Severity:   c.Query("severity"),
		Status:     c.Query("status"),
		AssignedTo: c.Query("assigned_to"),
	}

	incidents, err := h.incidents.FindAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch incidents"})
		return
	}

	if incidents == nil {
		incidents = []models.IncidentWithNames{}
	}
	c.JSON(http.StatusOK, gin.H{"data": incidents, "count": len(incidents)})
}

// GetByID godoc
// @Summary      Get a single incident by ID
// @Tags         incidents
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident UUID"
// @Success      200 {object} models.IncidentWithNames
// @Failure      404 {object} map[string]string
// @Router       /api/incidents/{id} [get]
func (h *IncidentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	inc, err := h.incidents.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if inc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}
	c.JSON(http.StatusOK, inc)
}

// Update godoc
// @Summary      Update an incident
// @Description  Admins can update any incident. Analysts can only update incidents they created.
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                      true "Incident UUID"
// @Param        body body models.UpdateIncidentRequest true "Fields to update"
// @Success      200  {object} models.Incident
// @Failure      400  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Router       /api/incidents/{id} [put]
func (h *IncidentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.incidents.FindByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	// Analysts may only update their own incidents
	if role.(string) == models.RoleAnalyst && existing.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update incidents you created"})
		return
	}

	var req models.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.incidents.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update incident"})
		return
	}

	_ = h.audit.Log(userID.(string),
		fmt.Sprintf("Updated incident: %s", existing.Title),
		"incident", &id)

	c.JSON(http.StatusOK, updated)
}

// Delete godoc
// @Summary      Delete an incident (Admin only)
// @Tags         incidents
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident UUID"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/incidents/{id} [delete]
func (h *IncidentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.incidents.FindByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	if err := h.incidents.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete incident"})
		return
	}

	userID, _ := c.Get("user_id")
	_ = h.audit.Log(userID.(string),
		fmt.Sprintf("Deleted incident: %s", existing.Title),
		"incident", &id)

	c.JSON(http.StatusOK, gin.H{"message": "incident deleted"})
}
