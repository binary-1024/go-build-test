package service

import (
	"context"
	"fmt"
	"time"

	"microservice/internal/cache"
	"microservice/internal/logger"
	"microservice/internal/models"
	"microservice/internal/repository"

	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers(page, limit int) ([]*models.User, int64, error)
}

// userService 用户服务实现
type userService struct {
	repo   repository.UserRepository
	cache  *cache.RedisClient
	logger logger.Logger
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository, cache *cache.RedisClient, logger logger.Logger) UserService {
	return &userService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	s.logger.Info("创建用户", "username", req.Username)

	// 检查用户名是否已存在
	existingUser, err := s.repo.GetByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("检查用户名失败", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.repo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("检查邮箱失败", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		IsActive: true,
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		s.logger.Error("密码加密失败", "error", err)
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		s.logger.Error("创建用户失败", "error", err)
		return nil, err
	}

	s.logger.Info("用户创建成功", "user_id", user.ID)
	return user, nil
}

// GetUser 获取用户
func (s *userService) GetUser(id uint) (*models.User, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	var cachedUser models.User

	ctx := context.Background()
	if err := s.cache.Get(ctx, cacheKey, &cachedUser); err == nil {
		s.logger.Debug("从缓存获取用户", "user_id", id)
		return &cachedUser, nil
	}

	// 从数据库获取
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("获取用户失败", "user_id", id, "error", err)
		return nil, err
	}

	// 缓存用户信息
	if err := s.cache.Set(ctx, cacheKey, user, 5*time.Minute); err != nil {
		s.logger.Warn("缓存用户信息失败", "user_id", id, "error", err)
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	s.logger.Info("更新用户", "user_id", id)

	// 检查用户是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("用户不存在", "user_id", id, "error", err)
		return nil, err
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 更新用户
	if err := s.repo.Update(id, updates); err != nil {
		s.logger.Error("更新用户失败", "user_id", id, "error", err)
		return nil, err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	ctx := context.Background()
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("删除用户缓存失败", "user_id", id, "error", err)
	}

	// 返回更新后的用户
	return s.repo.GetByID(id)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	s.logger.Info("删除用户", "user_id", id)

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("删除用户失败", "user_id", id, "error", err)
		return err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	ctx := context.Background()
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("删除用户缓存失败", "user_id", id, "error", err)
	}

	return nil
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(page, limit int) ([]*models.User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(offset, limit)
}
