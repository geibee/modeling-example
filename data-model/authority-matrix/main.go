package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/modeling-example/data-model/authority-matrix/handlers"
	"github.com/modeling-example/data-model/authority-matrix/middleware"
	"github.com/modeling-example/data-model/authority-matrix/repository"
)

func main() {
	// 環境変数からデータベース設定を取得
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "authority_user")
	dbPassword := getEnv("DB_PASSWORD", "authority_pass")
	dbName := getEnv("DB_NAME", "authority_db")

	// 接続文字列を構築
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// データベースを初期化
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("データベースのオープンに失敗:", err)
	}
	defer db.Close()

	// 接続をテスト
	if err := db.Ping(); err != nil {
		log.Fatal("データベースへの接続に失敗:", err)
	}
	log.Println("PostgreSQLデータベースに接続しました")

	// 注: テーブルと初期データはDocker ComposeのSQLスクリプトで作成されます
	// ここでテーブル作成やデータ挿入を行う必要はありません

	// リポジトリを初期化
	repo := repository.NewRepository(db)

	// ミドルウェアを初期化
	authMiddleware := middleware.NewAuthMiddleware(repo)

	// ハンドラーを初期化
	budgetHandler := handlers.NewBudgetHandler(repo, db)
	adminHandler := handlers.NewAdminHandler(repo, db)
	salesHandler := handlers.NewSalesHandler(repo, db)
	inventoryHandler := handlers.NewInventoryHandler(repo, db)
	userHandler := handlers.NewUserHandler(repo)
	authHandler := handlers.NewAuthHandler(repo)
	permissionHandler := handlers.NewPermissionHandler(repo)

	// 静的ファイルサーバー
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 管理画面ルート (SCREEN_ID=004)
	adminScreenHandler := authMiddleware.AuthorizeScreen("004")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/admin/index.html")
		}),
	)
	http.Handle("/admin", middleware.ExtractUserID()(adminScreenHandler))

	// ユーザーマスタ画面ルート (SCREEN_ID=005)
	userScreenHandler := authMiddleware.AuthorizeScreen("005")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/user/index.html")
		}),
	)
	http.Handle("/user", middleware.ExtractUserID()(userScreenHandler))

	// 管理API (API_ID=400) - すべての管理APIに認可チェックを適用
	adminAPIMiddleware := authMiddleware.AuthorizeAPI("400")
	
	// メニュー設定API
	menuSettingsHandler := adminAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminHandler.GetMenuSettings(w, r)
		case "POST":
			adminHandler.CreateMenuSetting(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/menu-settings", middleware.ExtractUserID()(menuSettingsHandler))
	
	menuSettingsDetailHandler := adminAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			adminHandler.UpdateMenuSetting(w, r)
		case "DELETE":
			adminHandler.DeleteMenuSetting(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/menu-settings/", middleware.ExtractUserID()(menuSettingsDetailHandler))

	// APIマッピングAPI
	apiMappingsHandler := adminAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminHandler.GetApiScreenMappings(w, r)
		case "POST":
			adminHandler.CreateApiScreenMapping(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/api-screen-mappings", middleware.ExtractUserID()(apiMappingsHandler))
	
	apiMappingsDetailHandler := adminAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			adminHandler.UpdateApiScreenMapping(w, r)
		case "DELETE":
			adminHandler.DeleteApiScreenMapping(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/api-screen-mappings/", middleware.ExtractUserID()(apiMappingsDetailHandler))

	// その他の管理API
	rolesHandler := adminAPIMiddleware(http.HandlerFunc(adminHandler.GetRoles))
	screensHandler := adminAPIMiddleware(http.HandlerFunc(adminHandler.GetScreens))
	apisHandler := adminAPIMiddleware(http.HandlerFunc(adminHandler.GetApis))
	organizationsHandler := adminAPIMiddleware(http.HandlerFunc(adminHandler.GetOrganizations))
	
	http.Handle("/api/roles", middleware.ExtractUserID()(rolesHandler))
	http.Handle("/api/screens", middleware.ExtractUserID()(screensHandler))
	http.Handle("/api/apis", middleware.ExtractUserID()(apisHandler))
	http.Handle("/api/organizations", middleware.ExtractUserID()(organizationsHandler))
	
	roleScreenPermHandler := adminAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			adminHandler.GetRoleScreenPermissions(w, r)
		} else if r.Method == "POST" || r.Method == "DELETE" {
			adminHandler.UpdateRoleScreenPermission(w, r)
		}
	}))
	http.Handle("/api/role-screen-permissions/", middleware.ExtractUserID()(roleScreenPermHandler))
	

	// ユーザー管理API (API_ID=500)
	userAPIMiddleware := authMiddleware.AuthorizeAPI("500")
	
	// ユーザー一覧取得・作成
	usersHandler := userAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminHandler.GetUsers(w, r)
		case "POST":
			adminHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/users", middleware.ExtractUserID()(usersHandler))
	
	// 個別ユーザー取得・更新・削除
	userDetailHandler := userAPIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminHandler.GetUser(w, r)
		case "PUT":
			adminHandler.UpdateUser(w, r)
		case "DELETE":
			adminHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.Handle("/api/users/", middleware.ExtractUserID()(userDetailHandler))

	// ルートを設定
	// 予算画面ルート (SCREEN_ID=001)
	budgetScreenHandler := authMiddleware.AuthorizeScreen("001")(
		http.HandlerFunc(budgetHandler.HandleBudgetScreen),
	)
	
	// 予算APIルート (API_ID=100)
	budgetAPIHandler := authMiddleware.AuthorizeAPI("100")(
		http.HandlerFunc(budgetHandler.HandleBudgetAPI),
	)

	// 販売画面ルート (SCREEN_ID=002)
	salesScreenHandler := authMiddleware.AuthorizeScreen("002")(
		http.HandlerFunc(salesHandler.HandleSalesScreen),
	)
	
	// 販売APIルート (API_ID=200)
	salesAPIHandler := authMiddleware.AuthorizeAPI("200")(
		http.HandlerFunc(salesHandler.HandleSalesAPI),
	)

	// 在庫画面ルート (SCREEN_ID=003)
	inventoryScreenHandler := authMiddleware.AuthorizeScreen("003")(
		http.HandlerFunc(inventoryHandler.HandleInventoryScreen),
	)
	
	// 在庫APIルート (API_ID=300)
	inventoryAPIHandler := authMiddleware.AuthorizeAPI("300")(
		http.HandlerFunc(inventoryHandler.HandleInventoryAPI),
	)

	// すべてのルートにユーザー抽出ミドルウェアを適用
	http.Handle("/budget", middleware.ExtractUserID()(budgetScreenHandler))
	http.Handle("/back-budget1", middleware.ExtractUserID()(budgetAPIHandler))
	http.Handle("/sales", middleware.ExtractUserID()(salesScreenHandler))
	http.Handle("/back-sales", middleware.ExtractUserID()(salesAPIHandler))
	http.Handle("/inventory", middleware.ExtractUserID()(inventoryScreenHandler))
	http.Handle("/back-inventory", middleware.ExtractUserID()(inventoryAPIHandler))

	// ホームページ（メニュー画面）- 認証必要
	homeHandler := middleware.ExtractUserID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}))
	http.Handle("/", homeHandler)

	// ログイン画面（認証不要）
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})

	// 認証API（認証不要）
	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/logout", authHandler.Logout)
	http.HandleFunc("/api/session", authHandler.GetCurrentSession)
	
	// ユーザー一覧取得用の認証不要エンドポイント
	http.HandleFunc("/api/users-list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminHandler.GetUsersForSwitcher(w, r)
	})

	// ユーザー情報API（認証必要）
	http.Handle("/api/current-user", middleware.ExtractUserID()(http.HandlerFunc(userHandler.GetCurrentUser)))
	http.Handle("/api/user-permissions", middleware.ExtractUserID()(http.HandlerFunc(permissionHandler.GetUserPermissions)))

	log.Println("サーバーを:8080で起動します")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}