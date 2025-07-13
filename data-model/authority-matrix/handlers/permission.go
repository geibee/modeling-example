package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type PermissionHandler struct {
	repo *repository.Repository
}

func NewPermissionHandler(repo *repository.Repository) *PermissionHandler {
	return &PermissionHandler{
		repo: repo,
	}
}

// GetUserPermissions はユーザーの画面権限一覧を返します
func (h *PermissionHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	// コンテキストからユーザーIDを取得
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ユーザーのロールを取得
	userRoles, err := h.repo.GetUserRoles(userID)
	if err != nil {
		http.Error(w, "Failed to get user roles", http.StatusInternalServerError)
		return
	}

	// 各ロールが持つ画面権限を取得
	screenPermissions := make(map[string]bool)
	for _, roleID := range userRoles {
		screens, err := h.repo.GetRoleScreenPermissions(roleID)
		if err != nil {
			continue
		}
		for _, screenID := range screens {
			screenPermissions[screenID] = true
		}
	}

	// レスポンスを作成
	response := map[string]interface{}{
		"user_id": userID,
		"roles":   userRoles,
		"screens": screenPermissions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}