package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type InventoryHandler struct {
	repo *repository.Repository
	db   *sql.DB
}

func NewInventoryHandler(repo *repository.Repository, db *sql.DB) *InventoryHandler {
	return &InventoryHandler{
		repo: repo,
		db:   db,
	}
}

// HandleInventoryScreen は/inventory画面へのリクエストを処理します
func (h *InventoryHandler) HandleInventoryScreen(w http.ResponseWriter, r *http.Request) {
	// 新しいスタイルの在庫画面を返す
	http.ServeFile(w, r, "static/inventory.html")
}

// HandleInventoryAPI は/back-inventoryのAPIリクエストを処理します
func (h *InventoryHandler) HandleInventoryAPI(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	// ユーザーの組織を取得
	organizations, err := h.repo.GetUserOrganizations(userID)
	if err != nil {
		http.Error(w, "組織情報の取得に失敗", http.StatusInternalServerError)
		return
	}

	// 組織IDのリストを作成
	var orgIDs []string
	for _, org := range organizations {
		orgIDs = append(orgIDs, org.OrgID)
	}

	// 組織の在庫データを取得
	query := `
		SELECT i.code, i.org_id, o.org_name, i.name, i.category, i.stock, i.status, i.updated
		FROM inventory i
		JOIN organization o ON i.org_id = o.org_id
		WHERE i.org_id = ANY($1)
		ORDER BY i.code
	`
	
	rows, err := h.db.Query(query, pq.Array(orgIDs))
	if err != nil {
		http.Error(w, fmt.Sprintf("在庫データの取得に失敗: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var inventoryData []map[string]interface{}
	for rows.Next() {
		var (
			code, orgID, orgName, name, category, status, updated string
			stock int
		)
		if err := rows.Scan(&code, &orgID, &orgName, &name, &category, &stock, &status, &updated); err != nil {
			continue
		}
		
		inventoryData = append(inventoryData, map[string]interface{}{
			"code":     code,
			"org_id":   orgID,
			"org_name": orgName,
			"name":     name,
			"category": category,
			"stock":    stock,
			"status":   status,
			"updated":  updated,
		})
	}

	// データがない場合は空の配列を返す
	if inventoryData == nil {
		inventoryData = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventoryData)
}

