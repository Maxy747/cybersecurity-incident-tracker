package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/marwans200/cyberguard/internal/repository"
)

// AuditHandler serves audit log endpoints (Admin only).
type AuditHandler struct {
	audit *repository.AuditRepo
}

// NewAuditHandler constructs an AuditHandler.
func NewAuditHandler(audit *repository.AuditRepo) *AuditHandler {
	return &AuditHandler{audit: audit}
}

// List godoc
// @Summary      List all audit logs (Admin only)
// @Description  Returns paginated audit log entries for every action in the system.
// @Tags         audit
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query int false "Page size (default 50)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {array}  models.AuditLogWithUser
// @Failure      403 {object} map[string]string
// @Router       /api/audit-logs [get]
func (h *AuditHandler) List(c *gin.Context) {
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	logs, err := h.audit.FindAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch audit logs"})
		return
	}

	total, _ := h.audit.Count()

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
