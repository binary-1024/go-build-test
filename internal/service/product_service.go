package service

import (
	"context"
	"fmt"
	"time"

	"github.com/binary-1024/go-build-test/internal/cache"
	"github.com/binary-1024/go-build-test/internal/logger"
	"github.com/binary-1024/go-build-test/internal/models"
	"github.com/binary-1024/go-build-test/internal/repository"
)

// ProductService 产品服务接口
type ProductService interface {
	CreateProduct(req *models.CreateProductRequest) (*models.Product, error)
	GetProduct(id uint) (*models.Product, error)
	UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(id uint) error
	ListProducts(query *models.ProductQuery) (*models.ProductListResponse, error)
}

// productService 产品服务实现
type productService struct {
	repo   repository.ProductRepository
	cache  *cache.RedisClient
	logger logger.Logger
}

// NewProductService 创建产品服务
func NewProductService(repo repository.ProductRepository, cache *cache.RedisClient, logger logger.Logger) ProductService {
	return &productService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// CreateProduct 创建产品
func (s *productService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	s.logger.Info("创建产品", "name", req.Name)

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		IsActive:    true,
	}

	if err := s.repo.Create(product); err != nil {
		s.logger.Error("创建产品失败", "error", err)
		return nil, err
	}

	s.logger.Info("产品创建成功", "product_id", product.ID)
	return product, nil
}

// GetProduct 获取产品
func (s *productService) GetProduct(id uint) (*models.Product, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("product:%d", id)
	var cachedProduct models.Product

	ctx := context.Background()
	if err := s.cache.Get(ctx, cacheKey, &cachedProduct); err == nil {
		s.logger.Debug("从缓存获取产品", "product_id", id)
		return &cachedProduct, nil
	}

	// 从数据库获取
	product, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("获取产品失败", "product_id", id, "error", err)
		return nil, err
	}

	// 缓存产品信息
	if err := s.cache.Set(ctx, cacheKey, product, 10*time.Minute); err != nil {
		s.logger.Warn("缓存产品信息失败", "product_id", id, "error", err)
	}

	return product, nil
}

// UpdateProduct 更新产品
func (s *productService) UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	s.logger.Info("更新产品", "product_id", id)

	// 检查产品是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("产品不存在", "product_id", id, "error", err)
		return nil, err
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 更新产品
	if err := s.repo.Update(id, updates); err != nil {
		s.logger.Error("更新产品失败", "product_id", id, "error", err)
		return nil, err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("product:%d", id)
	ctx := context.Background()
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("删除产品缓存失败", "product_id", id, "error", err)
	}

	// 返回更新后的产品
	return s.repo.GetByID(id)
}

// DeleteProduct 删除产品
func (s *productService) DeleteProduct(id uint) error {
	s.logger.Info("删除产品", "product_id", id)

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("删除产品失败", "product_id", id, "error", err)
		return err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("product:%d", id)
	ctx := context.Background()
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("删除产品缓存失败", "product_id", id, "error", err)
	}

	return nil
}

// ListProducts 获取产品列表
func (s *productService) ListProducts(query *models.ProductQuery) (*models.ProductListResponse, error) {
	products, total, err := s.repo.List(query)
	if err != nil {
		s.logger.Error("获取产品列表失败", "error", err)
		return nil, err
	}

	return &models.ProductListResponse{
		Products: make([]models.Product, len(products)),
		Total:    total,
		Page:     query.Page,
		Limit:    query.Limit,
	}, nil
}
