package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// User 用户模型
type User struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	IsActive bool      `json:"is_active"`
	CreateAt time.Time `json:"created_at"`
}

// Product 产品模型
type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	CreateAt    time.Time `json:"created_at"`
}

// Response 统一响应格式
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 内存数据库
type MemoryDB struct {
	users    map[uint]*User
	products map[uint]*Product
	userID   uint
	prodID   uint
	mutex    sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	db := &MemoryDB{
		users:    make(map[uint]*User),
		products: make(map[uint]*Product),
		userID:   1,
		prodID:   1,
	}

	// 初始化一些测试数据
	db.users[1] = &User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		FullName: "系统管理员",
		IsActive: true,
		CreateAt: time.Now(),
	}

	db.products[1] = &Product{
		ID:          1,
		Name:        "Go学习指南",
		Description: "从入门到精通的Go语言学习资料",
		Price:       99.99,
		Stock:       100,
		Category:    "书籍",
		IsActive:    true,
		CreateAt:    time.Now(),
	}

	return db
}

func (db *MemoryDB) GetAllUsers() []*User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	users := make([]*User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}
	return users
}

func (db *MemoryDB) GetUserByID(id uint) (*User, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, exists := db.users[id]
	return user, exists
}

func (db *MemoryDB) CreateUser(user *User) *User {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user.ID = db.userID
	user.CreateAt = time.Now()
	db.users[db.userID] = user
	db.userID++
	return user
}

func (db *MemoryDB) GetAllProducts() []*Product {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	products := make([]*Product, 0, len(db.products))
	for _, product := range db.products {
		products = append(products, product)
	}
	return products
}

func (db *MemoryDB) GetProductByID(id uint) (*Product, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	product, exists := db.products[id]
	return product, exists
}

func (db *MemoryDB) CreateProduct(product *Product) *Product {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	product.ID = db.prodID
	product.CreateAt = time.Now()
	product.IsActive = true
	db.products[db.prodID] = product
	db.prodID++
	return product
}

// 全局数据库实例
var db *MemoryDB

func main() {
	fmt.Println("=== 简化版微服务启动 ===")

	// 初始化内存数据库
	db = NewMemoryDB()

	// 设置路由
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/users", usersHandler)
	http.HandleFunc("/api/v1/users/", userHandler)
	http.HandleFunc("/api/v1/products", productsHandler)
	http.HandleFunc("/api/v1/products/", productHandler)

	// 添加CORS中间件
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("🚀 微服务启动成功!")
	fmt.Println("📍 服务地址: http://localhost:8080")
	fmt.Println("")
	fmt.Println("📚 API端点:")
	fmt.Println("  GET  /health              - 健康检查")
	fmt.Println("  GET  /api/v1/users        - 获取用户列表")
	fmt.Println("  POST /api/v1/users        - 创建用户")
	fmt.Println("  GET  /api/v1/users/{id}   - 获取指定用户")
	fmt.Println("  GET  /api/v1/products     - 获取产品列表")
	fmt.Println("  POST /api/v1/products     - 创建产品")
	fmt.Println("  GET  /api/v1/products/{id} - 获取指定产品")
	fmt.Println("")
	fmt.Println("🧪 测试命令:")
	fmt.Println("  curl http://localhost:8080/health")
	fmt.Println("  curl http://localhost:8080/api/v1/users")
	fmt.Println("  curl http://localhost:8080/api/v1/products")
	fmt.Println("")
	fmt.Println("按 Ctrl+C 停止服务")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := Response{
		Success: true,
		Message: "欢迎使用简化版微服务!",
		Data: map[string]interface{}{
			"service":     "Go 微服务示例",
			"version":     "1.0.0",
			"description": "这是一个简化版的微服务，演示基本的CRUD操作",
			"features": []string{
				"用户管理",
				"产品管理",
				"内存数据库",
				"RESTful API",
				"JSON响应格式",
			},
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	sendJSONResponse(w, http.StatusOK, response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "服务健康状态良好",
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"uptime":    "running",
			"version":   "1.0.0",
		},
	}

	sendJSONResponse(w, http.StatusOK, response)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users := db.GetAllUsers()
		response := Response{
			Success: true,
			Message: "获取用户列表成功",
			Data:    users,
		}
		sendJSONResponse(w, http.StatusOK, response)

	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			sendErrorResponse(w, "无效的JSON数据", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Email == "" {
			sendErrorResponse(w, "用户名和邮箱不能为空", http.StatusBadRequest)
			return
		}

		user.IsActive = true
		createdUser := db.CreateUser(&user)

		response := Response{
			Success: true,
			Message: "用户创建成功",
			Data:    createdUser,
		}
		sendJSONResponse(w, http.StatusCreated, response)

	default:
		sendErrorResponse(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// 提取ID from URL path
	pathParts := splitPath(r.URL.Path)
	if len(pathParts) < 4 || pathParts[3] == "" {
		sendErrorResponse(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(pathParts[3], 10, 32)
	if err != nil {
		sendErrorResponse(w, "无效的用户ID格式", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		user, exists := db.GetUserByID(uint(id))
		if !exists {
			sendErrorResponse(w, "用户不存在", http.StatusNotFound)
			return
		}

		response := Response{
			Success: true,
			Message: "获取用户成功",
			Data:    user,
		}
		sendJSONResponse(w, http.StatusOK, response)

	default:
		sendErrorResponse(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		products := db.GetAllProducts()
		response := Response{
			Success: true,
			Message: "获取产品列表成功",
			Data:    products,
		}
		sendJSONResponse(w, http.StatusOK, response)

	case "POST":
		var product Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			sendErrorResponse(w, "无效的JSON数据", http.StatusBadRequest)
			return
		}

		if product.Name == "" || product.Price <= 0 {
			sendErrorResponse(w, "产品名称和价格不能为空", http.StatusBadRequest)
			return
		}

		createdProduct := db.CreateProduct(&product)

		response := Response{
			Success: true,
			Message: "产品创建成功",
			Data:    createdProduct,
		}
		sendJSONResponse(w, http.StatusCreated, response)

	default:
		sendErrorResponse(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	// 提取ID from URL path
	pathParts := splitPath(r.URL.Path)
	if len(pathParts) < 4 || pathParts[3] == "" {
		sendErrorResponse(w, "无效的产品ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(pathParts[3], 10, 32)
	if err != nil {
		sendErrorResponse(w, "无效的产品ID格式", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		product, exists := db.GetProductByID(uint(id))
		if !exists {
			sendErrorResponse(w, "产品不存在", http.StatusNotFound)
			return
		}

		response := Response{
			Success: true,
			Message: "获取产品成功",
			Data:    product,
		}
		sendJSONResponse(w, http.StatusOK, response)

	default:
		sendErrorResponse(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Message: message,
	}
	sendJSONResponse(w, statusCode, response)
}

func splitPath(path string) []string {
	parts := []string{}
	for _, part := range []string{"", ""} {
		if part != "" {
			parts = append(parts, part)
		}
	}

	// 手动分割路径
	start := 0
	for i, c := range path {
		if c == '/' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}

	return parts
}
