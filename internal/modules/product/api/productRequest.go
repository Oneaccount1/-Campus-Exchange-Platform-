package api

import "time"

type CreateProductRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required"`
	Images      []string  `json:"images"`
	Category    string    `json:"category"`
	Condition   string    `json:"condition"`
	UserID      uint      `json:"user_id" binding:"required"`
	Status      string    `json:"status"`
	SoldAt      time.Time `json:"sold_at"`
}

type UpdateProductRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Images      []string  `json:"images"`
	Category    string    `json:"category"`
	Condition   string    `json:"condition"`
	Status      string    `json:"status"`
	SoldAt      time.Time `json:"sold_at"`
}

type GetProductsRequest struct {
	Page uint `json:"page" form:"page" binding:"required,min=1"`
	Size uint `json:"size" form:"size" binding:"required,min=1,max=100"`
}
