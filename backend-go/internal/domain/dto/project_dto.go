package dto

type ProjectCreateRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Platform     string  `json:"platform" binding:"required"`
	Client       string  `json:"client" binding:"required"`
	EstimatedFee float64 `json:"estimated_fee" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	Deadline     string  `json:"deadline"`
}

type ProjectUpdateRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Platform     string  `json:"platform" binding:"required"`
	Client       string  `json:"client" binding:"required"`
	EstimatedFee float64 `json:"estimated_fee" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	Deadline     string  `json:"deadline"`
}

type ProjectResponse struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Platform     string  `json:"platform"`
	Client       string  `json:"client"`
	EstimatedFee float64 `json:"estimated_fee"`
	Status       string  `json:"status"`
	Deadline     string  `json:"deadline"`
}
