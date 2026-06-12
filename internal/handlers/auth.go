package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/marwans200/cyberguard/internal/middleware"
	"github.com/marwans200/cyberguard/internal/models"
	"github.com/marwans200/cyberguard/internal/repository"
)

// AuthHandler groups all authentication endpoints.
type AuthHandler struct {
	users *repository.UserRepo
	audit *repository.AuditRepo
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(users *repository.UserRepo, audit *repository.AuditRepo) *AuthHandler {
	return &AuthHandler{users: users, audit: audit}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account. The first registered user is automatically assigned the admin role.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RegisterRequest true "Registration payload"
// @Success      201  {object} models.UserResponse
// @Failure      400  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate email
	existing, err := h.users.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// First user becomes admin automatically
	role := req.Role
	if role == "" {
		role = models.RoleViewer
	}
	count, err := h.users.CountAll()
	if err == nil && count == 0 {
		role = models.RoleAdmin
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	user, err := h.users.Create(req.Name, req.Email, string(hash), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	_ = h.audit.Log(user.ID, "User registered", "user", &user.ID)

	c.JSON(http.StatusCreated, models.UserResponse{
		ID: user.ID, Name: user.Name, Email: user.Email,
		Role: user.Role, CreatedAt: user.CreatedAt,
	})
}

// Login godoc
// @Summary      Authenticate and obtain a JWT
// @Description  Validates credentials and returns a signed JWT valid for 24 hours.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "Login credentials"
// @Success      200  {object} models.LoginResponse
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.users.FindByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Sign token
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}

	_ = h.audit.Log(user.ID, "User logged in", "user", &user.ID)

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenStr,
		User: models.UserResponse{
			ID: user.ID, Name: user.Name, Email: user.Email,
			Role: user.Role, CreatedAt: user.CreatedAt,
		},
	})
}

// Profile godoc
// @Summary      Get current user profile
// @Description  Returns the authenticated user's profile.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.UserResponse
// @Failure      401 {object} map[string]string
// @Router       /api/auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.users.FindByID(userID.(string))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, models.UserResponse{
		ID: user.ID, Name: user.Name, Email: user.Email,
		Role: user.Role, CreatedAt: user.CreatedAt,
	})
}
