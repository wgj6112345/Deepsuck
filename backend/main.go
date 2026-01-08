package main

import (
	"database/sql"
	"deepsuck/backend/handler"
	"deepsuck/backend/middleware"
	"deepsuck/backend/repository"
	"deepsuck/backend/usecase"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 初始化数据库
	db, err := sql.Open("sqlite3", "./deepsuck.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 创建表结构
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// 初始化 Repository
	convRepo := repository.NewSQLiteConversationRepository(db)
	msgRepo := repository.NewSQLiteMessageRepository(db)
	configRepo := repository.NewSQLiteConfigRepository(db)

	// 初始化 UseCase
	configUse := usecase.NewConfigUseCase(configRepo)
	convUse := usecase.NewConversationUseCase(convRepo, msgRepo)

	// 获取配置并创建 Agent Provider
	config, err := configUse.GetActiveConfig()
	if err != nil {
		log.Fatalf("Failed to get active config: %v", err)
	}

	// 如果没有配置 ProviderType，默认使用 mimo
	if config.ProviderType == "" {
		config.ProviderType = "mimo"
	}

	chatUse := usecase.NewChatUseCase(convRepo, msgRepo, configUse)

	// 初始化 Handler
	convHandler := handler.NewConversationHandler(convUse)
	configHandler := handler.NewConfigHandler(configUse)
	chatHandler := handler.NewChatHandler(chatUse)

	// 配置路由
	http.HandleFunc("/api/conversations", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			convHandler.GetConversations(w, r)
		case http.MethodPost:
			convHandler.CreateConversation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/conversations/", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			convHandler.GetConversation(w, r)
		case http.MethodDelete:
			convHandler.DeleteConversation(w, r)
		case http.MethodPut:
			convHandler.UpdateConversation(w, r)
		case http.MethodPost:
			// 检查是否是置顶操作
			if strings.HasSuffix(r.URL.Path, "/pin") {
				convHandler.TogglePin(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/chat", middleware.CORS(chatHandler.SendMessage))

	// 配置路由 - 新版多 Provider 配置
	http.HandleFunc("/api/config/active", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			configHandler.GetActiveConfig(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/config/providers", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			configHandler.GetAllProviders(w, r)
		case http.MethodPut:
			configHandler.SaveProviderConfig(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/config/providers/", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			configHandler.GetProviderConfig(w, r)
		case http.MethodPut:
			configHandler.SaveProviderConfig(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/config/activate", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			configHandler.ActivateProvider(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// 启动 HTTP 服务器
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createTables(db *sql.DB) error {
	// 创建 conversations 表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			pinned BOOLEAN DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create conversations table: %w", err)
	}

	// 为已存在的表添加 pinned 字段（如果不存在）
	_, err = db.Exec(`
		ALTER TABLE conversations ADD COLUMN pinned BOOLEAN DEFAULT 0
	`)
	if err != nil {
		// 字段可能已存在，忽略错误
		log.Printf("Note: pinned column may already exist: %v", err)
	}

	// 创建 messages 表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			thinking TEXT,
			thinking_enabled BOOLEAN DEFAULT 0,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// 创建 config 表（用于存储 activeProvider 等全局配置）
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create config table: %w", err)
	}

	// 创建 provider_configs 表（存储多个 Provider 的配置）
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS provider_configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_type TEXT UNIQUE NOT NULL,
			api_key TEXT,
			model_name TEXT NOT NULL,
			base_url TEXT NOT NULL,
			enabled BOOLEAN DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create provider_configs table: %w", err)
	}

	return nil
}