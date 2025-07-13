package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/modeling-example/data-model/authority-matrix/models"
	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type BudgetHandler struct {
	repo *repository.Repository
	db   *sql.DB
}

func NewBudgetHandler(repo *repository.Repository, db *sql.DB) *BudgetHandler {
	return &BudgetHandler{
		repo: repo,
		db:   db,
	}
}

// HandleBudgetAPI は/back-budget1へのリクエストを処理します
func (h *BudgetHandler) HandleBudgetAPI(w http.ResponseWriter, r *http.Request) {
	// 認可はミドルウェアで既にチェック済み
	// APIへのアクセス権限がある場合のみここに到達する

	switch r.Method {
	case "GET":
		// 予算データを返す
		h.getBudgetData(w, r)
	case "POST":
		// リクエストに基づいて作成か更新かを判断
		var req models.BudgetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "無効なリクエストボディ", http.StatusBadRequest)
			return
		}

		// 簡素化のため、IDが存在するかで作成か更新かを判断
		if len(req.BudgetList) > 0 && req.BudgetList[0].ID != "" {
			h.updateBudget(w, r, req)
		} else {
			h.createBudget(w, r, req)
		}
	default:
		http.Error(w, "許可されていないメソッド", http.StatusMethodNotAllowed)
	}
}

// updateBudget は組織認可付きの予算更新を処理します
func (h *BudgetHandler) updateBudget(w http.ResponseWriter, r *http.Request, req models.BudgetRequest) {
	userID := r.Context().Value("userID").(string)
	
	// トランザクションを開始
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "サーバー内部エラー", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// ユーザーの所属組織を取得
	userOrganizations, err := h.repo.GetUserOrganizations(userID)
	if err != nil {
		http.Error(w, "ユーザー組織の取得に失敗", http.StatusInternalServerError)
		return
	}

	// 高速検索用のマップを作成
	orgMap := make(map[string]bool)
	for _, org := range userOrganizations {
		orgMap[org.OrgID] = true
	}

	// リクエスト内の各予算を処理
	for _, budget := range req.BudgetList {
		// まず予算の組織IDを取得
		var orgID string
		err := h.db.QueryRow("SELECT org_id FROM budget WHERE id = $1", budget.ID).Scan(&orgID)
		if err != nil {
			http.Error(w, "予算情報の取得に失敗", http.StatusInternalServerError)
			return
		}

		// ユーザーがこの組織にアクセス権があるかチェック
		if !orgMap[orgID] {
			errMsg := fmt.Sprintf("認可エラー: 組織%sの予算を更新できません", orgID)
			http.Error(w, errMsg, http.StatusForbidden)
			return
		}

		// 予算を更新
		if err := h.repo.UpdateBudget(budget); err != nil {
			http.Error(w, "予算の更新に失敗", http.StatusInternalServerError)
			return
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		http.Error(w, "トランザクションのコミットに失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// createBudget は組織認可付きの予算作成を処理します
func (h *BudgetHandler) createBudget(w http.ResponseWriter, r *http.Request, req models.BudgetRequest) {
	userID := r.Context().Value("userID").(string)
	
	// トランザクションを開始
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "サーバー内部エラー", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// ユーザーの所属組織を取得
	userOrganizations, err := h.repo.GetUserOrganizations(userID)
	if err != nil {
		http.Error(w, "ユーザー組織の取得に失敗", http.StatusInternalServerError)
		return
	}

	// 高速検索用のマップを作成
	orgMap := make(map[string]bool)
	for _, org := range userOrganizations {
		orgMap[org.OrgID] = true
	}

	// リクエスト内の各予算を処理
	for _, budget := range req.BudgetList {
		// 予算作成時は、リクエストに組織IDを含める必要がある
		// ここでは仮に最初の所属組織を使用
		if len(userOrganizations) > 0 {
			budget.OrgID = userOrganizations[0].OrgID
		}

		// ユーザーがこの組織にアクセス権があるかチェック
		if !orgMap[budget.OrgID] {
			errMsg := fmt.Sprintf("認可エラー: 組織%sの予算を作成できません", budget.OrgID)
			http.Error(w, errMsg, http.StatusForbidden)
			return
		}

		// 予算を作成
		if err := h.repo.CreateBudget(budget); err != nil {
			http.Error(w, "予算の作成に失敗", http.StatusInternalServerError)
			return
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		http.Error(w, "トランザクションのコミットに失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// getBudgetData は予算データを返します
func (h *BudgetHandler) getBudgetData(w http.ResponseWriter, r *http.Request) {
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

	// 組織の予算データを取得
	query := `
		SELECT b.id, b.department_id, b.org_id, o.org_name, b.amount, b.description, b.created_at, b.updated_at
		FROM budget b
		JOIN organization o ON b.org_id = o.org_id
		WHERE b.org_id = ANY($1)
		ORDER BY b.org_id, b.department_id, b.id
	`
	
	rows, err := h.db.Query(query, pq.Array(orgIDs))
	if err != nil {
		http.Error(w, fmt.Sprintf("予算データの取得に失敗: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var budgets []map[string]interface{}
	for rows.Next() {
		var (
			id, departmentID, orgID, orgName, description string
			amount float64
			createdAt, updatedAt string
		)
		if err := rows.Scan(&id, &departmentID, &orgID, &orgName, &amount, &description, &createdAt, &updatedAt); err != nil {
			continue
		}
		
		// モックデータ形式に合わせる
		budgets = append(budgets, map[string]interface{}{
			"id":            id,
			"department":    departmentID,
			"org_id":        orgID,
			"org_name":      orgName,
			"period":        "2024年Q1", // 簡単のため固定値
			"budget":        int(amount),
			"used":          int(amount * 0.6), // 使用済額は60%と仮定
			"remaining":     int(amount * 0.4), // 残額は40%と仮定
			"usage":         60.0, // 使用率60%と仮定
		})
	}

	// データがない場合は空の配列を返す
	if budgets == nil {
		budgets = []map[string]interface{}{}
	}

	response := map[string]interface{}{
		"status":  "success",
		"budgets": budgets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// HandleBudgetScreen は/budget画面へのリクエストを処理します
func (h *BudgetHandler) HandleBudgetScreen(w http.ResponseWriter, r *http.Request) {
	// 新しいスタイルの予算画面を返す
	http.ServeFile(w, r, "static/budget.html")
}