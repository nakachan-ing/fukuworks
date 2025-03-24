package dto

type ProjectCreateRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Platform     string  `json:"platform" binding:"required"`
	Client       string  `json:"client" binding:"required"`
	EstimatedFee float64 `json:"estimated_fee" binding:"gte=0"`
	Status       string  `json:"status" binding:"required"`
	Deadline     string  `json:"deadline"`
}

type ProjectUpdateRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Platform     string  `json:"platform" binding:"required"`
	Client       string  `json:"client" binding:"required"`
	EstimatedFee float64 `json:"estimated_fee" binding:"gte=0"`
	Status       string  `json:"status" binding:"required"`
	Deadline     string  `json:"deadline"`
}

type ProjectResponse struct {
	Number       uint    `json:"project_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Platform     string  `json:"platform"`
	Client       string  `json:"client"`
	EstimatedFee float64 `json:"estimated_fee"`
	Status       string  `json:"status"`
	Deadline     string  `json:"deadline"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type ProjectResponseForOwner struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	Number       uint    `json:"project_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Platform     string  `json:"platform"`
	Client       string  `json:"client"`
	EstimatedFee float64 `json:"estimated_fee"`
	Status       string  `json:"status"`
	Deadline     string  `json:"deadline"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    string  `json:"deleted_at"`
}
