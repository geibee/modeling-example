package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"org-hierarchy/handlers"
	"org-hierarchy/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// データベース接続
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 静的ファイルの配信
	e.Static("/", "public")

	// ハンドラーの初期化
	departmentRepo := models.NewDepartmentRepository(db)
	orgAttrRepo := models.NewOrganizationAttributeRepository(db)
	
	departmentHandler := handlers.NewDepartmentHandler(departmentRepo)
	orgAttrHandler := handlers.NewOrganizationAttributeHandler(orgAttrRepo)

	// ルーティング
	api := e.Group("/api")
	
	// 部門のエンドポイント
	api.GET("/departments", departmentHandler.GetAll)
	api.GET("/departments/:id", departmentHandler.GetByID)
	api.POST("/departments", departmentHandler.Create)
	api.PUT("/departments/:id", departmentHandler.Update)
	api.DELETE("/departments/:id", departmentHandler.Delete)

	// 期間別組織属性のエンドポイント
	api.GET("/organization-attributes", orgAttrHandler.GetAll)
	api.GET("/organization-attributes/:department_id/:effective_date", orgAttrHandler.GetByID)
	api.POST("/organization-attributes", orgAttrHandler.Create)
	api.PUT("/organization-attributes/:department_id/:effective_date", orgAttrHandler.Update)
	api.DELETE("/organization-attributes/:department_id/:effective_date", orgAttrHandler.Delete)
	
	// 特定部門の履歴を取得
	api.GET("/departments/:id/history", orgAttrHandler.GetDepartmentHistory)
	
	// 特定日付時点の組織階層を取得
	api.GET("/hierarchy", orgAttrHandler.GetHierarchyByDate)

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// サーバー起動
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func connectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// リトライロジック
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database, retrying... (%d/10)", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}