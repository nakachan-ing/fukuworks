package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	httpstd "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	userhttp "github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
)

func setupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserRepo := NewMockUserRepo()
	userHandler := userhttp.NewUserHandler(mockUserRepo)

	r.POST("/", userHandler.PostUser)
	r.POST("/login", userHandler.Login)

	authorized := r.Group("/:user")
	authorized.Use(middleware.ReservedPathGuard()) // 必要ならこれも入れる
	authorized.Use(middleware.AuthMiddleware())    // 認可ミドルウェア
	{
		authorized.PATCH("", userHandler.UpdateUser)
		authorized.DELETE("", userHandler.SoftDeleteUser)
	}

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

	body := `{"name":"","email":"invalid@example.com","password":"secret123"}`
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

func TestPostUser_InvalidEmailFormat(t *testing.T) {
	r := setupRouterForTest()

	// email形式が不正
	payload := `{"name":"kyota","email":"invalid-email","password":"pass1234"}`
	req, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "email")
}

func TestGetUser_NotFound(t *testing.T) {
	r := setupRouterForTest()

	req, _ := httpstd.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestPostUser_Duplicate(t *testing.T) {
	r := setupRouterForTest()

	// 同じname/emailで2回POSTする
	payload := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	req1, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(payload))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, httpstd.StatusCreated, w1.Code)

	req2, _ := httpstd.NewRequest("POST", "/", bytes.NewBufferString(payload))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, httpstd.StatusConflict, w2.Code)
	assert.Contains(t, w2.Body.String(), "already exists")
}

func TestUpdateUser_NotFound(t *testing.T) {
	r := setupRouterForTest()

	update := map[string]string{
		"name":  "hoge",
		"email": "hoge@example.com",
	}
	data, _ := json.Marshal(update)
	req, _ := http.NewRequest("PATCH", "/ghostuser", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-token-for-ghostuser")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestSoftDeleteUser_NotFound(t *testing.T) {
	r := setupRouterForTest()

	req, _ := http.NewRequest("DELETE", "/ghostuser", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-ghostuser")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}
