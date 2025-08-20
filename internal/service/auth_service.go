package service

import (
	"fmt"

	"microservice/internal/auth"
	"microservice/internal/logger"
	"microservice/internal/models"
	"microservice/internal/repository"

	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(token string) (*auth.Claims, error)
}

// authService 认证服务实现
type authService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
	logger     logger.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwtManager *auth.JWTManager, logger logger.Logger) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Login 用户登录
func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	s.logger.Info("用户登录", "username", req.Username)

	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn("用户不存在", "username", req.Username)
			return nil, fmt.Errorf("用户名或密码错误")
		}
		s.logger.Error("获取用户失败", "error", err)
		return nil, err
	}

	// 检查用户是否激活
	if !user.IsActive {
		s.logger.Warn("用户已禁用", "username", req.Username)
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 验证密码
	if err := user.CheckPassword(req.Password); err != nil {
		s.logger.Warn("密码错误", "username", req.Username)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("生成token失败", "error", err)
		return nil, err
	}

	s.logger.Info("用户登录成功", "user_id", user.ID)

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ValidateToken 验证token
func (s *authService) ValidateToken(token string) (*auth.Claims, error) {
	return s.jwtManager.ValidateToken(token)
}
