package repository

import (
	"microservice/internal/models"

	"gorm.io/gorm"
)

// ProductRepository 产品仓库接口
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	List(query *models.ProductQuery) ([]*models.Product, int64, error)
}

// productRepository 产品仓库实现
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓库
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create 创建产品
func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID 根据ID获取产品
func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新产品
func (r *productRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除产品
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// List 获取产品列表
func (r *productRepository) List(query *models.ProductQuery) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	db := r.db.Model(&models.Product{})

	// 添加搜索条件
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if query.MinPrice > 0 {
		db = db.Where("price >= ?", query.MinPrice)
	}

	if query.MaxPrice > 0 {
		db = db.Where("price <= ?", query.MaxPrice)
	}

	if query.Search != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.Limit
	err = db.Offset(offset).Limit(query.Limit).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
