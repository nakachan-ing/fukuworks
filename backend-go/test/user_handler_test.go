package test

import (
	"bytes"
	httpstd "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	userhttp "github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http"
	"github.com/stretchr/testify/assert"
)

func setupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserRepo := NewMockUserRepo()
	userHandler := userhttp.NewUserHandler(mockUserRepo)

	r.POST("/", userHandler.PostUser)
	r.POST("/login", userHandler.Login)

	return r
}

func TestPostUser_Success(t *testing.T) {
	r := setupRouterForTest()

	body := `{"name":"kyota","email":"kyota@example.com","password":"secret123"}`
	req, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "kyota@example.com")
}

func TestPostUser_ValidationError(t *testing.T) {
	r := setupRouterForTest()

	body := `{"email":"invalid@example.com","password":"secret123"}`
	req, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusBadRequest, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "Name") || strings.Contains(w.Body.String(), "name"))
	assert.Contains(t, w.Body.String(), "is required")
}

func TestLogin_Success(t *testing.T) {
	r := setupRouterForTest()

	body := `{"name":"kyota","email":"kyota@example.com","password":"secret123"}`
	req, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	loginBody := `{"name":"kyota","password":"secret123"}`
	loginReq, _ := httpstd.NewRequest("POST", "/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	r.ServeHTTP(loginRes, loginReq)

	assert.Equal(t, httpstd.StatusOK, loginRes.Code)
	assert.Contains(t, loginRes.Body.String(), "mock-token-for-kyota")
}

func TestLogin_Failure(t *testing.T) {
	r := setupRouterForTest()

	body := `{"name":"kyota","email":"kyota@example.com","password":"secret123"}`
	req, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	loginBody := `{"name":"kyota","password":"wrongpass"}`
	loginReq, _ := httpstd.NewRequest("POST", "/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	r.ServeHTTP(loginRes, loginReq)

	assert.Equal(t, httpstd.StatusUnauthorized, loginRes.Code)
	assert.Contains(t, loginRes.Body.String(), "Invalid credentials")
}
