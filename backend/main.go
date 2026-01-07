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
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/chat", middleware.CORS(chatHandler.SendMessage))

	http.HandleFunc("/api/config", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			configHandler.GetConfig(w, r)
		case http.MethodPut:
			configHandler.UpdateConfig(w, r)
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
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create conversations table: %w", err)
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

	// 创建 config 表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create config table: %w", err)
	}

	return nil
}