package dto

type TaskCreateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
	Status      string `json:"status" binding:"required,oneof=Todo Doing Done"`
	Priority    string `json:"priority" binding:"required,oneof=Low Medium High"`
	DueDate     string `json:"due_date"`
}

type TaskUpdateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
	Status      string `json:"status" binding:"required,oneof=Todo Doing Done"`
	Priority    string `json:"priority" binding:"required,oneof=Low Medium High"`
	DueDate     string `json:"due_date"`
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
