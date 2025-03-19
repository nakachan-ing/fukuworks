package dto

type TaskCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
}

type TaskResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}
