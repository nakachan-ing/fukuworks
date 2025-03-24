package dto

type TaskCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

type TaskUpdateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

type TaskResponse struct {
	Number      uint   `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TaskResponseForOwner struct {
	ID          uint   `json:"id"`
	ProjectID   uint   `json:"project_id"`
	Number      uint   `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}
