package dto

import "time"

// RegisterRequest представляет данные для регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=100" example:"securePassword123"`
	Name     string `json:"name" binding:"required,min=2,max=255" example:"Иван Петров"`
}

// LoginRequest представляет данные для входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securePassword123"`
}

// RefreshRequest используется для обновления токенов
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AuthResponse представляет ответ после аутентификации
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
}

// UserResponse представляет безопасные данные пользователя
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"Иван Петров"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T12:30:00Z"`
}

// WhoamiResponse представляет ответ эндпоинта /whoami
type WhoamiResponse struct {
	Authenticated bool          `json:"authenticated" example:"true"`
	User          *UserResponse `json:"user,omitempty"`
}
