package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/marwans200/cyberguard/internal/repository"
)

// DashboardHandler groups all dashboard/stats endpoints.
type DashboardHandler struct {
	incidents *repository.IncidentRepo
	notes     *repository.NoteRepo
	users     *repository.UserRepo
}

// NewDashboardHandler constructs a DashboardHandler.
func NewDashboardHandler(incidents *repository.IncidentRepo, notes *repository.NoteRepo, users *repository.UserRepo) *DashboardHandler {
	return &DashboardHandler{incidents: incidents, notes: notes, users: users}
}

// Stats godoc
// @Summary      Get incident count statistics
// @Description  Returns total, open, investigating, resolved, and closed incident counts.
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]int
// @Router       /api/dashboard/stats [get]
func (h *DashboardHandler) Stats(c *gin.Context) {
	counts, err := h.incidents.CountByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch stats"})
		return
	}
	c.JSON(http.StatusOK, counts)
}

// SeverityDistribution godoc
// @Summary      Get incident count by severity
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]int
// @Router       /api/dashboard/severity [get]
func (h *DashboardHandler) SeverityDistribution(c *gin.Context) {
	counts, err := h.incidents.CountBySeverity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch severity stats"})
		return
	}
	c.JSON(http.StatusOK, counts)
}

// Trends godoc
// @Summary      Get incident trends over time
// @Description  Returns incident counts per day. Use ?days=N to control the window (default 30).
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Param        days query int false "Number of days to look back (default 30)"
// @Success      200 {array} map[string]interface{}
// @Router       /api/dashboard/trends [get]
func (h *DashboardHandler) Trends(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 {
			days = n
		}
	}

	trend, err := h.incidents.TrendsByDay(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch trends"})
		return
	}
	if trend == nil {
		trend = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"data": trend, "days": days})
}
