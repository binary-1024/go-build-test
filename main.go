package main

import (
	"log"
	"os"

	"github.com/binary-1024/go-build-test/internal/api"
	"github.com/binary-1024/go-build-test/internal/auth"
	"github.com/binary-1024/go-build-test/internal/cache"
	"github.com/binary-1024/go-build-test/internal/config"
	"github.com/binary-1024/go-build-test/internal/database"
	"github.com/binary-1024/go-build-test/internal/logger"
	"github.com/binary-1024/go-build-test/internal/middleware"
	"github.com/binary-1024/go-build-test/internal/repository"
	"github.com/binary-1024/go-build-test/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title 微服务API
// @version 1.0
// @description 这是一个完整的微服务示例，包含认证、缓存、数据库等功能
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载 .env 文件: %v", err)
	}

	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger := logger.NewLogger(cfg.LogLevel)
	logger.Info("启动微服务应用")

	// 初始化数据库
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("数据库连接失败", "error", err)
	}
	logger.Info("数据库连接成功")

	// 初始化Redis缓存
	redisClient := cache.NewRedisClient(cfg.RedisURL)
	logger.Info("Redis连接成功")

	// 初始化JWT管理器
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// 初始化仓库层
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, redisClient, logger)
	productService := service.NewProductService(productRepo, redisClient, logger)
	authService := service.NewAuthService(userRepo, jwtManager, logger)

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由器
	router := gin.New()

	// 添加中间件
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// 初始化API路由
	apiHandler := api.NewHandler(userService, productService, authService, logger)
	apiHandler.SetupRoutes(router, jwtManager)

	// 启动服务器
	logger.Info("服务器启动", "port", cfg.Port, "environment", cfg.Environment)
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("服务器启动失败", "error", err)
	}
}
