package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/marwans200/cyberguard/internal/handlers"
	"github.com/marwans200/cyberguard/internal/middleware"
	"github.com/marwans200/cyberguard/internal/models"
	"github.com/marwans200/cyberguard/internal/repository"
)

// Setup registers all API routes on the provided Gin engine.
func Setup(r *gin.Engine) {
	// ── Repository layer ─────────────────────────────────────────────────
	userRepo     := &repository.UserRepo{}
	incidentRepo := &repository.IncidentRepo{}
	noteRepo     := &repository.NoteRepo{}
	auditRepo    := &repository.AuditRepo{}

	// ── Handler layer ────────────────────────────────────────────────────
	authH      := handlers.NewAuthHandler(userRepo, auditRepo)
	incidentH  := handlers.NewIncidentHandler(incidentRepo, auditRepo)
	noteH      := handlers.NewNoteHandler(noteRepo, incidentRepo, auditRepo)
	dashboardH := handlers.NewDashboardHandler(incidentRepo, noteRepo, userRepo)
	auditH     := handlers.NewAuditHandler(auditRepo)

	// ── Swagger UI ───────────────────────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Health check ─────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "CyberGuard"})
	})

	api := r.Group("/api")

	// ── Auth (public) ────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authH.Profile)
	}

	// ── Authenticated routes ──────────────────────────────────────────────
	protected := api.Group("/", middleware.AuthMiddleware())

	// Incidents
	incidents := protected.Group("/incidents")
	{
		incidents.GET("", incidentH.List)
		incidents.GET("/:id", incidentH.GetByID)
		incidents.POST("",
			middleware.RequireRoles(models.RoleAdmin, models.RoleAnalyst),
			incidentH.Create)
		incidents.PUT("/:id",
			middleware.RequireRoles(models.RoleAdmin, models.RoleAnalyst),
			incidentH.Update)
		incidents.DELETE("/:id",
			middleware.RequireRoles(models.RoleAdmin),
			incidentH.Delete)

		// Notes nested under incidents
		incidents.POST("/:id/notes",
			middleware.RequireRoles(models.RoleAdmin, models.RoleAnalyst),
			noteH.Create)
		incidents.GET("/:id/notes", noteH.List)
	}

	// Dashboard
	dashboard := protected.Group("/dashboard")
	{
		dashboard.GET("/stats", dashboardH.Stats)
		dashboard.GET("/severity", dashboardH.SeverityDistribution)
		dashboard.GET("/trends", dashboardH.Trends)
	}

	// Audit logs (admin only)
	protected.GET("/audit-logs",
		middleware.RequireRoles(models.RoleAdmin),
		auditH.List)
}
