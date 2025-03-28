package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	httpstd "net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	userhttp "github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
)

func setupTaskRouterWithAuth() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserRepo := NewMockUserRepo()
	userHandler := userhttp.NewUserHandler(mockUserRepo)
	mockProjectRepo := NewMockProjectRepo()
	projectHandler := userhttp.NewProjectHandler(mockProjectRepo)
	mockTaskRepo := NewMockTaskRepo()
	taskHandler := userhttp.NewTaskHandler(mockTaskRepo)

	r.POST("/signup", userHandler.PostUser)

	authorized := r.Group(":user")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/projects", projectHandler.PostProject)
		authorized.POST("/projects/:pid/tasks", taskHandler.PostTask)
		authorized.GET("/projects/:pid/tasks/:tid", taskHandler.GetTask)
		authorized.PATCH("/projects/:pid/tasks/:tid", taskHandler.UpdateTask)
		authorized.DELETE("/projects/:pid/tasks/:tid", taskHandler.SoftDeleteTask)
	}

	return r
}

func TestGetTask_ForbiddenForOtherUser(t *testing.T) {
	r := setupTaskRouterWithAuth()

	// ユーザー登録
	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	// プロジェクト作成
	project := map[string]interface{}{
		"title":         "テストPJ",
		"platform":      "個人",
		"client":        "client",
		"estimated_fee": 1000,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	projectJson, _ := json.Marshal(project)
	createProjectReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createProjectReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createProjectReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createProjectReq)

	// タスク作成
	task := map[string]interface{}{
		"title":       "タスクA",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-10",
	}
	taskJson, _ := json.Marshal(task)
	createTaskReq, _ := httpstd.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBuffer(taskJson))
	createTaskReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createTaskReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createTaskReq)

	// 他人がタスクを取得しようとする
	req, _ := httpstd.NewRequest("GET", "/otheruser/projects/1/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestUpdateTask_ForbiddenForOtherUser(t *testing.T) {
	r := setupTaskRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "テストPJ",
		"platform":      "個人",
		"client":        "client",
		"estimated_fee": 1000,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	projectJson, _ := json.Marshal(project)
	createProjectReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createProjectReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createProjectReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createProjectReq)

	task := map[string]interface{}{
		"title":       "タスクA",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-10",
	}
	taskJson, _ := json.Marshal(task)
	createTaskReq, _ := httpstd.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBuffer(taskJson))
	createTaskReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createTaskReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createTaskReq)

	updated := map[string]interface{}{
		"title":       "タスク更新",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Doing",
		"due_date":    "2025-06-15",
	}
	updatedJson, _ := json.Marshal(updated)
	updateReq, _ := httpstd.NewRequest("PATCH", "/otheruser/projects/1/tasks/1", bytes.NewBuffer(updatedJson))
	updateReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	updateReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, updateReq)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestSoftDeleteTask_ForbiddenForOtherUser(t *testing.T) {
	r := setupTaskRouterWithAuth()

	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	signupReq, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	signupReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), signupReq)

	project := map[string]interface{}{
		"title":         "テストPJ",
		"platform":      "個人",
		"client":        "client",
		"estimated_fee": 1000,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	projectJson, _ := json.Marshal(project)
	createProjectReq, _ := httpstd.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(projectJson))
	createProjectReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createProjectReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createProjectReq)

	task := map[string]interface{}{
		"title":       "タスクA",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-10",
	}
	taskJson, _ := json.Marshal(task)
	createTaskReq, _ := httpstd.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBuffer(taskJson))
	createTaskReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	createTaskReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), createTaskReq)

	deleteReq, _ := httpstd.NewRequest("DELETE", "/otheruser/projects/1/tasks/1", nil)
	deleteReq.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, deleteReq)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestPostTask_Success(t *testing.T) {
	r := setupTaskRouterWithAuth()

	// ユーザー & プロジェクト作成
	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	r.ServeHTTP(httptest.NewRecorder(), newJSONReq("POST", "/signup", signup))

	project := map[string]interface{}{
		"title":         "PJ",
		"platform":      "Web",
		"client":        "client",
		"estimated_fee": 0,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	pjJson, _ := json.Marshal(project)
	r.ServeHTTP(httptest.NewRecorder(), newAuthedJSONReq("POST", "/kyota/projects", pjJson))

	task := map[string]interface{}{
		"title":       "タスク1",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-15",
	}
	taskJson, _ := json.Marshal(task)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, newAuthedJSONReq("POST", "/kyota/projects/1/tasks", taskJson))

	assert.Equal(t, httpstd.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "タスク1")
}

func TestUpdateTask_Success(t *testing.T) {
	r := setupTaskRouterWithAuth()
	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	r.ServeHTTP(httptest.NewRecorder(), newJSONReq("POST", "/signup", signup))

	pj := map[string]interface{}{
		"title":         "PJ",
		"platform":      "Web",
		"client":        "client",
		"estimated_fee": 0,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	pjJson, _ := json.Marshal(pj)
	r.ServeHTTP(httptest.NewRecorder(), newAuthedJSONReq("POST", "/kyota/projects", pjJson))

	task := map[string]interface{}{
		"title":       "タスク1",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-15",
	}
	taskJson, _ := json.Marshal(task)
	r.ServeHTTP(httptest.NewRecorder(), newAuthedJSONReq("POST", "/kyota/projects/1/tasks", taskJson))

	update := map[string]interface{}{
		"title":       "更新済みタスク",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Doing",
		"due_date":    "2025-07-01",
	}
	updateJson, _ := json.Marshal(update)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, newAuthedJSONReq("PATCH", "/kyota/projects/1/tasks/1", updateJson))

	assert.Equal(t, httpstd.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "更新済みタスク")
}

func TestSoftDeleteTask_Success(t *testing.T) {
	r := setupTaskRouterWithAuth()
	signup := `{"name":"kyota","email":"kyota@example.com","password":"secret"}`
	r.ServeHTTP(httptest.NewRecorder(), newJSONReq("POST", "/signup", signup))

	pj := map[string]interface{}{
		"title":         "PJ",
		"platform":      "Web",
		"client":        "client",
		"estimated_fee": 0,
		"status":        "Open",
		"deadline":      "2025-06-01",
	}
	pjJson, _ := json.Marshal(pj)
	r.ServeHTTP(httptest.NewRecorder(), newAuthedJSONReq("POST", "/kyota/projects", pjJson))

	task := map[string]interface{}{
		"title":       "削除対象",
		"description": "タスクの説明", // ★ 追加
		"priority":    "High",   // ★ 追加
		"status":      "Todo",
		"due_date":    "2025-06-15",
	}
	taskJson, _ := json.Marshal(task)
	r.ServeHTTP(httptest.NewRecorder(), newAuthedJSONReq("POST", "/kyota/projects/1/tasks", taskJson))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, newAuthedJSONReq("DELETE", "/kyota/projects/1/tasks/1", nil))
	assert.Equal(t, httpstd.StatusNoContent, w.Code)
}

func newJSONReq(method, path, body string) *httpstd.Request {
	req, _ := httpstd.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newAuthedJSONReq(method, path string, body []byte) *httpstd.Request {
	req, _ := httpstd.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestPostTask_ValidationError(t *testing.T) {
	r := setupTaskRouterWithAuth()

	_ = performSignup(r, "kyota")
	_ = performCreateProject(r, "kyota", "Bearer mock-token-for-kyota")

	payload := `{"title":""}`
	req, _ := httpstd.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Title")
}

func TestGetTask_NotFound(t *testing.T) {
	r := setupTaskRouterWithAuth()
	_ = performSignup(r, "kyota")
	_ = performCreateProject(r, "kyota", "Bearer mock-token-for-kyota")

	req, _ := httpstd.NewRequest("GET", "/kyota/projects/1/tasks/999", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestGetTask_Forbidden(t *testing.T) {
	r := setupTaskRouterWithAuth()
	_ = performSignup(r, "kyota")
	_ = performSignup(r, "hacker")
	_ = performCreateProject(r, "kyota", "Bearer mock-token-for-kyota")
	_ = performCreateTask(r, "kyota", 1, "Bearer mock-token-for-kyota")

	req, _ := httpstd.NewRequest("GET", "/kyota/projects/1/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-hacker")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not authorized")
}

func TestUpdateTask_NotFound(t *testing.T) {
	r := setupTaskRouterWithAuth()
	_ = performSignup(r, "kyota")
	_ = performCreateProject(r, "kyota", "Bearer mock-token-for-kyota")

	payload := map[string]string{
		"title":       "Updated Task",
		"description": "New Desc",
		"priority":    "Medium",
		"status":      "Todo",
		"due_date":    "2025-12-31",
	}
	data, _ := json.Marshal(payload)
	req, _ := httpstd.NewRequest("PATCH", "/kyota/projects/1/tasks/999", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestSoftDeleteTask_NotFound(t *testing.T) {
	r := setupTaskRouterWithAuth()
	_ = performSignup(r, "kyota")
	_ = performCreateProject(r, "kyota", "Bearer mock-token-for-kyota")

	req, _ := httpstd.NewRequest("DELETE", "/kyota/projects/1/tasks/999", nil)
	req.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, httpstd.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestTaskValidation_InvalidPriority(t *testing.T) {
	r := setupTaskRouterWithAuth()
	signup := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req1)

	project := map[string]interface{}{
		"title":         "プロジェクト",
		"platform":      "web",
		"client":        "client",
		"estimated_fee": 1000,
		"status":        "Open",
		"deadline":      "2025-12-31",
	}
	body, _ := json.Marshal(project)
	req2, _ := http.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(body))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req2)

	task := map[string]interface{}{
		"title":       "タスク",
		"description": "説明",
		"status":      "Open",
		"priority":    "Extreme",
		"due_date":    "2025-10-01",
	}
	taskBody, _ := json.Marshal(task)
	req3, _ := http.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBuffer(taskBody))
	req3.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req3.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req3)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Priority")
}

func TestTaskValidation_TitleTooLong(t *testing.T) {
	r := setupTaskRouterWithAuth()
	signup := `{"name":"kyota","email":"kyota@example.com","password":"pass1234"}`
	req1, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(signup))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req1)

	project := map[string]interface{}{
		"title":         "プロジェクト",
		"platform":      "web",
		"client":        "client",
		"estimated_fee": 1000,
		"status":        "Open",
		"deadline":      "2025-12-31",
	}
	body, _ := json.Marshal(project)
	req2, _ := http.NewRequest("POST", "/kyota/projects", bytes.NewBuffer(body))
	req2.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req2)

	task := map[string]interface{}{
		"title":       strings.Repeat("あ", 101),
		"description": "説明",
		"status":      "Open",
		"priority":    "Low",
		"due_date":    "2025-10-01",
	}
	taskBody, _ := json.Marshal(task)
	req3, _ := http.NewRequest("POST", "/kyota/projects/1/tasks", bytes.NewBuffer(taskBody))
	req3.Header.Set("Authorization", "Bearer mock-token-for-kyota")
	req3.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req3)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Title")
}

// --- ヘルパー関数 ---

func performSignup(r *gin.Engine, username string) *httptest.ResponseRecorder {
	payload := map[string]string{
		"name":     username,
		"email":    username + "@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := httpstd.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performCreateProject(r *gin.Engine, username, token string) *httptest.ResponseRecorder {
	payload := map[string]string{
		"title":    "TestProject",
		"platform": "Web",
		"client":   "TestClient",
		"status":   "InProgress",
		"deadline": "2025-12-31",
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := httpstd.NewRequest("POST", "/"+username+"/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performCreateTask(r *gin.Engine, username string, pid int, token string) *httptest.ResponseRecorder {
	strPid := strconv.Itoa(pid)
	payload := map[string]interface{}{
		"title":       "Test Task",
		"description": "Test Description",
		"priority":    "High",
		"status":      "Todo",
		"due_date":    "2025-12-31",
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := httpstd.NewRequest("POST", "/"+username+"/projects/"+strPid+"/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
