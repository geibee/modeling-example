package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// GetCurrentUser は現在のユーザー情報を返すAPIハンドラー
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// コンテキストからユーザーIDを取得（認証ミドルウェアで設定済み）
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ユーザーのロール情報を取得
	roles, err := h.repo.GetUserRoles(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ユーザーの組織情報を取得
	organizations, err := h.repo.GetUserOrganizations(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ユーザー情報をマップで作成
	userInfo := map[string]interface{}{
		"user_id": userID,
		"roles":   roles,
		"organizations": organizations,
	}

	// ロール名を取得（最初のロールを使用）
	if len(roles) > 0 {
		userInfo["role_name"] = getRoleName(roles[0])
		userInfo["permissions"] = getPermissionsSummary(roles[0])
	}

	// ユーザー名を設定（実際のシステムではDBから取得）
	userInfo["name"] = getUserDisplayName(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// getUserDisplayName はユーザーIDから表示名を取得（モック実装）
func getUserDisplayName(userID string) string {
	switch userID {
	case "user001":
		return "管理職ユーザー"
	case "user002":
		return "営業担当ユーザー"
	case "user003":
		return "在庫管理者ユーザー"
	case "admin":
		return "システム管理者"
	default:
		return "ゲストユーザー"
	}
}

// getRoleName はロールIDからロール名を取得（モック実装）
func getRoleName(roleID string) string {
	switch roleID {
	case "R01":
		return "管理職"
	case "R02":
		return "営業担当"
	case "R03":
		return "在庫管理者"
	case "R04":
		return "システム管理者"
	default:
		return "未定義"
	}
}

// getPermissionsSummary はロールIDから権限の概要を取得（モック実装）
func getPermissionsSummary(roleID string) string {
	switch roleID {
	case "R01":
		return "予算管理"
	case "R02":
		return "販売管理"
	case "R03":
		return "在庫管理"
	case "R04":
		return "全権限"
	default:
		return "権限なし"
	}
}