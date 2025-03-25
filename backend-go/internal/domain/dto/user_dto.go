package dto

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"min=8,max=64"` // 任意（空なら変更なし）
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserResponseForOwner struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
