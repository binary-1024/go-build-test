package models

import (
	"time"
	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	Category    string         `json:"category"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	Stock       *int    `json:"stock" binding:"omitempty,min=0"`
	Category    string  `json:"category"`
	IsActive    *bool   `json:"is_active"`
}

// ProductQuery 产品查询参数
type ProductQuery struct {
	Page     int     `form:"page,default=1" binding:"min=1"`
	Limit    int     `form:"limit,default=10" binding:"min=1,max=100"`
	Category string  `form:"category"`
	MinPrice float64 `form:"min_price" binding:"min=0"`
	MaxPrice float64 `form:"max_price" binding:"min=0"`
	Search   string  `form:"search"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}
