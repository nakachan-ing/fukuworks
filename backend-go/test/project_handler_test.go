package test

import (
	"bytes"
	"encoding/json"
	httpstd "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	userhttp "github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
)

func setupProjectRouterWithAuth() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserRepo := NewMockUserRepo()
	userHandler := userhttp.NewUserHandler(mockUserRepo)
	mockProjectRepo := NewMockProjectRepo()
	projectHandler := userhttp.NewProjectHandler(mockProjectRepo)

	r.POST("/signup", userHandler.PostUser)

	authorized := r.Group(":user")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/projects", projectHandler.PostProject)
		authorized.GET("/projects", projectHandler.GetAllProjectsByUser)
		authorized.GET("/projects/:pid", projectHandler.GetProject)
		authorized.PATCH("/projects/:pid", projectHandler.UpdateProject)
		authorized.DELETE("/projects/:pid", projectHandler.SoftDeleteProject)
	}

	return r
}

func TestGetProject_ForbiddenForOtherUser(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	signupRes := httptest.NewRecorder()
	r.ServeHTTP(signupRes, signupReq)

	project := map[string]interface{}{
		"title":         "秘密のプロジェクト",
		"platform":      "個人",
		"client":        "テストクライアント",
		"estimated_fee": 10000,
		"status":        "In progress",
		"deadline":      "2025-04-01",
	}
	projectJson, _ := json.Marshal(project)

	createReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	r.ServeHTTP(createRes, createReq)

	assert.Equal(t, httpstd.StatusCreated, createRes.Code)

	// 他人が取得しようとする
	req, _ := httpstd.NewRequest("GET", "/otheruser/projects/1", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestUpdateProject_ForbiddenForOtherUser(t *testing.T) {
	r := setupProjectRouterWithAuth()

	// ユーザー作成とプロジェクト登録
	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	_ = httptest.NewRecorder()
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "更新対象プロジェクト",
		"platform":      "個人",
		"client":        "クライアントA",
		"estimated_fee": 5000,
		"status":        "Planning",
		"deadline":      "2025-05-01",
	}
	projectJson, _ := json.Marshal(project)

	createReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createReq.Header.Set("Content-Type", "application/json")
	_ = httptest.NewRecorder()
	r.ServeHTTP(httptest.NewRecorder(), createReq)

	// 他人がPATCHで更新しようとする
	updated := map[string]interface{}{
		"title":         "不正な更新",
		"platform":      "チーム",
		"client":        "悪意あるクライアント",
		"estimated_fee": 20000,
		"status":        "Closed",
		"deadline":      "2025-06-01",
	}
	updatedJson, _ := json.Marshal(updated)

	req, _ := httpstd.NewRequest("PATCH", "/otheruser/projects/1", bytes.NewBuffer(updatedJson))
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestSoftDeleteProject_ForbiddenForOtherUser(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	_ = httptest.NewRecorder()
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "削除対象プロジェクト",
		"platform":      "個人",
		"client":        "クライアントZ",
		"estimated_fee": 8000,
		"status":        "Done",
		"deadline":      "2025-06-15",
	}
	projectJson, _ := json.Marshal(project)

	createReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createReq)

	// 他人がDELETEで削除しようとする
	req, _ := httpstd.NewRequest("DELETE", "/otheruser/projects/1", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestUpdateProject_Success(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "プロジェクトA",
		"platform":      "個人",
		"client":        "テスト",
		"estimated_fee": 3000,
		"status":        "Planning",
		"deadline":      "2025-05-01",
	}
	projectJson, _ := json.Marshal(project)
	createReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createReq)

	updated := map[string]interface{}{
		"title":         "プロジェクトA更新",
		"platform":      "チーム",
		"client":        "更新クライアント",
		"estimated_fee": 3500,
		"status":        "In progress",
		"deadline":      "2025-06-01",
	}
	updatedJson, _ := json.Marshal(updated)
	req, _ := httpstd.NewRequest("PATCH", "/kyota/projects/1", bytes.NewBuffer(updatedJson))
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "プロジェクトA更新")
}

func TestSoftDeleteProject_Success(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "削除するプロジェクト",
		"platform":      "個人",
		"client":        "削除クライアント",
		"estimated_fee": 7000,
		"status":        "Done",
		"deadline":      "2025-07-01",
	}
	projectJson, _ := json.Marshal(project)
	createReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createReq)

	req, _ := httpstd.NewRequest("DELETE", "/kyota/projects/1", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusNoContent, w.Code)
}
