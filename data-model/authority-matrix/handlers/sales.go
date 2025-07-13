package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type SalesHandler struct {
	repo *repository.Repository
	db   *sql.DB
}

func NewSalesHandler(repo *repository.Repository, db *sql.DB) *SalesHandler {
	return &SalesHandler{
		repo: repo,
		db:   db,
	}
}

// HandleSalesScreen は/sales画面へのリクエストを処理します
func (h *SalesHandler) HandleSalesScreen(w http.ResponseWriter, r *http.Request) {
	// 新しいスタイルの販売画面を返す
	http.ServeFile(w, r, "static/sales.html")
}

// HandleSalesAPI は/back-salesのAPIリクエストを処理します
func (h *SalesHandler) HandleSalesAPI(w http.ResponseWriter, r *http.Request) {
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

	// 組織の販売データを取得
	query := `
		SELECT s.id, s.org_id, o.org_name, s.product, s.quantity, s.amount, s.date
		FROM sales s
		JOIN organization o ON s.org_id = o.org_id
		WHERE s.org_id = ANY($1)
		ORDER BY s.date DESC
	`
	
	rows, err := h.db.Query(query, pq.Array(orgIDs))
	if err != nil {
		http.Error(w, fmt.Sprintf("販売データの取得に失敗: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var salesData []map[string]interface{}
	for rows.Next() {
		var (
			id, orgID, orgName, product, date string
			quantity int
			amount float64
		)
		if err := rows.Scan(&id, &orgID, &orgName, &product, &quantity, &amount, &date); err != nil {
			continue
		}
		
		salesData = append(salesData, map[string]interface{}{
			"id":       id,
			"org_id":   orgID,
			"org_name": orgName,
			"product":  product,
			"quantity": quantity,
			"amount":   amount,
			"date":     date,
		})
	}

	// データがない場合は空の配列を返す
	if salesData == nil {
		salesData = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(salesData)
}

