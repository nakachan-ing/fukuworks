package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authGroup := r.Group("/:user")
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/protected", func(c *gin.Context) {
		user, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{"message": "Hello " + user.(string)})
	})

	return r
}

func TestAuthMiddleware_Success(t *testing.T) {
	r := setupAuthTestRouter()

	req, _ := http.NewRequest("GET", "/kyota/protected", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello kyota")
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupAuthTestRouter()

	req, _ := http.NewRequest("GET", "/kyota/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	r := setupAuthTestRouter()

	req, _ := http.NewRequest("GET", "/kyota/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token format")
}

func TestAuthMiddleware_UsernameMismatch(t *testing.T) {
	r := setupAuthTestRouter()

	req, _ := http.NewRequest("GET", "/otheruser/protected", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}
