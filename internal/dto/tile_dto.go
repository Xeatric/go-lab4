package dto

import "time"

// CreateTileRequest DTO для создания плитки
type CreateTileRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=200" example:"Брусчатка Классик"`
	Shape       string  `json:"shape" binding:"required,oneof=square rectangle hexagon circle" example:"rectangle"`
	Color       string  `json:"color" binding:"required,min=2,max=50" example:"серый"`
	Size        string  `json:"size" binding:"required,min=2,max=50" example:"200x100x60"`
	Material    string  `json:"material" binding:"omitempty,max=100" example:"бетон"`
	PricePerM2  float64 `json:"price_per_m2" binding:"required,min=0" example:"850.50"`
	Stock       int     `json:"stock" binding:"min=0" example:"1500"`
	Description string  `json:"description" binding:"max=1000" example:"Классическая брусчатка для пешеходных зон"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url,max=500" example:"https://example.com/brick.jpg"`
}

// UpdateTileRequest DTO для полного обновления
type UpdateTileRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=200" example:"Брусчатка Премиум"`
	Shape       string  `json:"shape" binding:"required,oneof=square rectangle hexagon circle" example:"square"`
	Color       string  `json:"color" binding:"required,min=2,max=50" example:"темно-серый"`
	Size        string  `json:"size" binding:"required,min=2,max=50" example:"300x300x60"`
	Material    string  `json:"material" binding:"omitempty,max=100" example:"гранит"`
	PricePerM2  float64 `json:"price_per_m2" binding:"required,min=0" example:"2500.00"`
	Stock       int     `json:"stock" binding:"min=0" example:"500"`
	Description string  `json:"description" binding:"max=1000" example:"Премиальная гранитная брусчатка"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url,max=500" example:"https://example.com/granite.jpg"`
}

// PatchTileRequest DTO для частичного обновления
type PatchTileRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=200" example:"Обновленная брусчатка"`
	Shape       *string  `json:"shape" binding:"omitempty,oneof=square rectangle hexagon circle"`
	Color       *string  `json:"color" binding:"omitempty,min=2,max=50"`
	Size        *string  `json:"size" binding:"omitempty,min=2,max=50"`
	Material    *string  `json:"material" binding:"omitempty,max=100"`
	PricePerM2  *float64 `json:"price_per_m2" binding:"omitempty,min=0" example:"2750.00"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0" example:"450"`
	Description *string  `json:"description" binding:"omitempty,max=1000"`
	ImageURL    *string  `json:"image_url" binding:"omitempty,url,max=500"`
}

// TileResponse DTO для ответа
type TileResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"Брусчатка Классик"`
	Shape       string    `json:"shape" example:"rectangle"`
	Color       string    `json:"color" example:"серый"`
	Size        string    `json:"size" example:"200x100x60"`
	Material    string    `json:"material" example:"бетон"`
	PricePerM2  float64   `json:"price_per_m2" example:"850.50"`
	Stock       int       `json:"stock" example:"1500"`
	Description string    `json:"description" example:"Классическая брусчатка для пешеходных зон"`
	ImageURL    string    `json:"image_url" example:"https://example.com/brick.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-15T12:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-15T12:30:00Z"`
}

// PaginationRequest DTO для пагинации
type PaginationRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1" default:"1" example:"2"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100" default:"10" example:"20"`
}

// PaginationResponse DTO для ответа с пагинацией
type PaginationResponse struct {
	Data       []TileResponse `json:"data"`
	Pagination Meta           `json:"meta"`
}

type Meta struct {
	Total      int64 `json:"total" example:"42"`
	Page       int   `json:"page" example:"2"`
	Limit      int   `json:"limit" example:"10"`
	TotalPages int64 `json:"total_pages" example:"5"`
}
