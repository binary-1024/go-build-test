# 微服务项目

这是一个完整的微服务应用示例，展示了企业级Go应用的架构设计，包含认证、缓存、数据库、日志等完整的技术栈。

## 学习目标

- 掌握微服务架构设计
- 学习企业级应用开发
- 理解分层架构和依赖注入
- 掌握JWT认证和Redis缓存
- 学习结构化日志和中间件

## 项目特点

- ✅ 分层架构设计
- ✅ JWT认证系统
- ✅ Redis缓存集成
- ✅ SQLite数据库
- ✅ 结构化日志
- ✅ 中间件支持
- ✅ 完整的CRUD API
- ✅ 错误处理和恢复

## 技术栈

### 核心框架
- **Gin** (`github.com/gin-gonic/gin`) - HTTP Web框架
- **GORM** (`gorm.io/gorm`) - ORM数据库框架

### 数据存储
- **SQLite** (`gorm.io/driver/sqlite`) - 嵌入式数据库
- **Redis** (`github.com/go-redis/redis/v8`) - 内存缓存

### 认证安全
- **JWT** (`github.com/golang-jwt/jwt/v4`) - JSON Web Token
- **bcrypt** (`golang.org/x/crypto`) - 密码加密

### 日志配置
- **Logrus** (`github.com/sirupsen/logrus`) - 结构化日志
- **godotenv** (`github.com/joho/godotenv`) - 环境变量

## 项目架构

```
05-microservice/
├── main.go                           # 应用入口
├── go.mod                           # 模块定义
├── internal/                        # 内部包
│   ├── api/                        # API层
│   │   └── handler.go              # HTTP处理器
│   ├── auth/                       # 认证模块
│   │   └── jwt.go                  # JWT管理
│   ├── cache/                      # 缓存层
│   │   └── redis.go                # Redis客户端
│   ├── config/                     # 配置管理
│   │   └── config.go               # 配置加载
│   ├── database/                   # 数据库层
│   │   └── database.go             # 数据库连接
│   ├── logger/                     # 日志系统
│   │   └── logger.go               # 日志接口
│   ├── middleware/                 # 中间件
│   │   └── middleware.go           # HTTP中间件
│   ├── models/                     # 数据模型
│   │   ├── user.go                 # 用户模型
│   │   └── product.go              # 产品模型
│   ├── repository/                 # 数据访问层
│   │   ├── user_repository.go      # 用户仓库
│   │   └── product_repository.go   # 产品仓库
│   └── service/                    # 业务逻辑层
│       ├── user_service.go         # 用户服务
│       ├── product_service.go      # 产品服务
│       └── auth_service.go         # 认证服务
└── microservice.db                 # SQLite数据库文件
```

## 快速开始

### 1. 环境准备

#### 安装Redis（可选）
```bash
# macOS
brew install redis
brew services start redis

# Ubuntu
sudo apt-get install redis-server
sudo systemctl start redis-server

# Windows
# 下载Redis并启动
```

#### 安装依赖
```bash
cd 05-microservice
go mod tidy
```

### 2. 启动服务
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 3. 测试API

#### 健康检查
```bash
curl http://localhost:8080/health
```

#### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "测试用户"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## API文档

### 认证相关

#### 用户注册
```
POST /api/v1/users
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string",
  "full_name": "string"
}
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}
```

### 用户管理（需要认证）

#### 获取用户列表
```
GET /api/v1/users?page=1&limit=10
Authorization: Bearer {token}
```

#### 获取指定用户
```
GET /api/v1/users/{id}
Authorization: Bearer {token}
```

#### 更新用户
```
PUT /api/v1/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "full_name": "string",
  "is_active": true
}
```

#### 删除用户
```
DELETE /api/v1/users/{id}
Authorization: Bearer {token}
```

### 产品管理（需要认证）

#### 创建产品
```
POST /api/v1/products
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "string",
  "description": "string",
  "price": 99.99,
  "stock": 100,
  "category": "string"
}
```

#### 获取产品列表
```
GET /api/v1/products?page=1&limit=10&category=electronics&min_price=10&max_price=1000&search=phone
Authorization: Bearer {token}
```

#### 获取指定产品
```
GET /api/v1/products/{id}
Authorization: Bearer {token}
```

#### 更新产品
```
PUT /api/v1/products/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "string",
  "description": "string",
  "price": 89.99,
  "stock": 50,
  "category": "string",
  "is_active": true
}
```

#### 删除产品
```
DELETE /api/v1/products/{id}
Authorization: Bearer {token}
```

## 架构详解

### 1. 分层架构

```
┌─────────────────┐
│   API Handler   │  ← HTTP请求处理，参数验证，响应格式化
├─────────────────┤
│    Middleware   │  ← 认证，日志，CORS，错误恢复
├─────────────────┤
│    Service      │  ← 业务逻辑，数据处理，缓存管理
├─────────────────┤
│   Repository    │  ← 数据访问，查询封装，事务管理
├─────────────────┤
│    Database     │  ← 数据持久化，GORM ORM
└─────────────────┘
```

### 2. 核心组件

#### 认证系统 (internal/auth/jwt.go)
```go
type JWTManager struct {
    secretKey string
}

func (j *JWTManager) GenerateToken(userID uint, username string) (string, error)
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error)
```

**功能:**
- JWT token生成和验证
- 用户身份信息编码
- token过期时间控制

#### 缓存系统 (internal/cache/redis.go)
```go
type RedisClient struct {
    client *redis.Client
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error
```

**功能:**
- 用户信息缓存
- 产品信息缓存
- 缓存过期管理

#### 数据库层 (internal/database/database.go)
```go
func NewConnection(databaseURL string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    // 自动迁移
    err = db.AutoMigrate(&models.User{}, &models.Product{})
    return db, err
}
```

**功能:**
- 数据库连接管理
- 自动表结构迁移
- 查询日志记录

### 3. 中间件系统

#### 认证中间件
```go
func Auth(jwtManager *auth.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        claims, err := jwtManager.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "无效的token",
            })
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

#### 日志中间件
```go
func Logger(logger logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)

        logger.Info("HTTP请求",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "latency", latency,
            "ip", c.ClientIP(),
        )
    }
}
```

## 数据模型

### 用户模型
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"uniqueIndex;not null"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Password  string         `json:"-" gorm:"not null"`
    FullName  string         `json:"full_name"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### 产品模型
```go
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
```

## 环境配置

### 配置项说明
```bash
ENVIRONMENT=development     # 运行环境
PORT=8080                  # 服务端口
DATABASE_URL=./microservice.db  # 数据库URL
REDIS_URL=redis://localhost:6379  # Redis连接
JWT_SECRET=my-secret-key   # JWT密钥
LOG_LEVEL=info            # 日志级别
```

### 生产环境配置建议
```bash
ENVIRONMENT=production
PORT=8080
DATABASE_URL=postgres://user:pass@localhost/dbname
REDIS_URL=redis://redis-server:6379
JWT_SECRET=your-random-secret-key
LOG_LEVEL=warn
```

## 安全特性

### 1. 密码安全
- 使用bcrypt加密存储密码
- 密码强度验证（最少6位）
- 密码不在API响应中返回

### 2. JWT认证
- token有效期24小时
- 安全的密钥管理
- token格式验证

### 3. 数据验证
- 输入参数验证
- 邮箱格式验证
- 数据类型检查

### 4. CORS配置
- 跨域请求控制
- 允许的方法和头部限制

## 性能优化

### 1. 缓存策略
```go
// 用户信息缓存5分钟
err := s.cache.Set(ctx, cacheKey, user, 5*time.Minute)

// 产品信息缓存10分钟
err := s.cache.Set(ctx, cacheKey, product, 10*time.Minute)
```

### 2. 数据库优化
- 使用索引加速查询
- 实现软删除避免数据丢失
- 查询分页减少内存使用

### 3. 连接池管理
- 数据库连接复用
- Redis连接池配置

## 错误处理

### 1. 统一错误响应
```json
{
  "success": false,
  "message": "错误描述"
}
```

### 2. 错误类型
- 400 Bad Request - 请求参数错误
- 401 Unauthorized - 认证失败
- 404 Not Found - 资源不存在
- 500 Internal Server Error - 服务器错误

### 3. 错误日志
```go
logger.Error("操作失败", "error", err, "user_id", userID)
```

## 监控和日志

### 1. 结构化日志
```json
{
  "level": "info",
  "msg": "HTTP请求",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "5.2ms",
  "ip": "127.0.0.1",
  "time": "2024-01-15T10:30:45Z"
}
```

### 2. 健康检查
```bash
curl http://localhost:8080/health
```

### 3. 性能监控
- 请求响应时间记录
- 错误率统计
- 缓存命中率

## 扩展建议

### 1. 添加API文档
```go
// 使用Swagger生成API文档
// @Summary 创建用户
// @Description 创建新用户账户
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "用户信息"
// @Success 201 {object} models.User
// @Router /users [post]
```

### 2. 添加单元测试
```go
func TestUserService_CreateUser(t *testing.T) {
    // 测试逻辑
}
```

### 3. 添加数据库迁移
```go
// 版本化的数据库迁移脚本
func Migrate_001_CreateUsersTable(db *gorm.DB) error {
    return db.AutoMigrate(&User{})
}
```

### 4. 添加配置验证
```go
func (c *Config) Validate() error {
    if c.JWTSecret == "" {
        return errors.New("JWT密钥不能为空")
    }
    return nil
}
```

## 部署指南

### Docker部署
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### 构建和运行
```bash
# 构建镜像
docker build -t microservice .

# 运行容器
docker run -p 8080:8080 microservice
```

## 常见问题

### Q: Redis连接失败怎么办？
A: 检查Redis服务是否启动，确认连接URL正确，应用会自动降级到无缓存模式。

### Q: 数据库迁移失败？
A: 检查数据库权限，确认go.mod中GORM版本兼容性。

### Q: JWT token失效？
A: 检查系统时间，确认密钥配置正确，token过期需要重新登录。

### Q: 如何扩展更多业务功能？
A: 按照现有的分层架构，在相应层添加新的模型、服务和API。

## 下一步学习

完成这个项目后，建议继续学习：
- 微服务通信（gRPC、消息队列）
- 服务注册与发现
- 分布式链路追踪
- 容器编排（Kubernetes）
- 性能测试和监控
- CI/CD流水线
---
修改测试 1 

修改测试 2 

修改测试 3