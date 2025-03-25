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

func TestPostProject_ValidationError(t *testing.T) {
	r := setupProjectRouterWithAuth()

	// ユーザー作成
	signupPayload := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// 必須フィールドが欠けたプロジェクト作成
	payload := `{"title":""}`
	req2, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBufferString(payload))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, httpstd.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "Title")
}

func TestPostProject_InvalidDeadline(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signupPayload := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	payload := `{"title":"test","platform":"web","client":"client","status":"draft","deadline":"31-12-2025"}`
	req2, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBufferString(payload))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, httpstd.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "Date format is invalid")
}

func TestGetProject_NotFound(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signupPayload := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	req2, _ := httpstd.NewRequest("GET", "/kyota/projects/999", nil)
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, httpstd.StatusNotFound, w2.Code)
	assert.Contains(t, w2.Body.String(), "Project not found")
}

func TestGetProject_Forbidden(t *testing.T) {
	r := setupProjectRouterWithAuth()

	// 正規ユーザーで作成
	signupPayload1 := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	projectPayload := `{"title":"project1","platform":"web","client":"client","status":"draft","deadline":"2025-12-31"}`
	req2, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBufferString(projectPayload))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// 他人がアクセス
	signupPayload2 := `{"name":"hacker","email":"hacker@example.com","password":"hack1234"}`
	req3, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload2))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	req4, _ := httpstd.NewRequest("GET", "/kyota/projects/1", nil)
	req4.Header.Set("Authorization", "Bearer mock-token-for-hacker")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	assert.Equal(t, httpstd.StatusForbidden, w4.Code)
	assert.Contains(t, w4.Body.String(), "not authorized")
}

func TestUpdateProject_InvalidID(t *testing.T) {
	r := setupProjectRouterWithAuth()

	signupPayload := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signupPayload))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	payload := map[string]string{
		"title":    "修正タイトル",
		"platform": "Web",
		"client":   "Test",
		"status":   "Open",
		"deadline": "2025-12-31",
	}
	jsonData, _ := json.Marshal(payload)

	req2, _ := httpstd.NewRequest("PATCH", "/kyota/projects/abc", bytes.NewBuffer(jsonData))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, httpstd.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "Project ID is invalid")
}
