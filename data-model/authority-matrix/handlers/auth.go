package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/modeling-example/data-model/authority-matrix/middleware"
	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type AuthHandler struct {
	repo *repository.Repository
}

func NewAuthHandler(repo *repository.Repository) *AuthHandler {
	return &AuthHandler{
		repo: repo,
	}
}

// Login はユーザーをログインさせます
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginReq.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// ユーザーの存在確認（ロールを取得）
	roles, err := h.repo.GetUserRoles(loginReq.UserID)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました: " + err.Error(), http.StatusInternalServerError)
		return
	}
	if len(roles) == 0 {
		http.Error(w, "ユーザー " + loginReq.UserID + " にロールが設定されていません", http.StatusUnauthorized)
		return
	}

	// セッションに保存
	session, _ := middleware.Store.Get(r, "auth-session")
	session.Values["user_id"] = loginReq.UserID
	session.Values["authenticated"] = true
	
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// ユーザー情報を返す
	userInfo := map[string]interface{}{
		"user_id": loginReq.UserID,
		"roles":   roles,
		"status":  "logged_in",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// Logout はユーザーをログアウトさせます
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "auth-session")
	
	// セッションをクリア
	session.Values["user_id"] = nil
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_out"})
}

// GetCurrentSession は現在のセッション情報を返します
func (h *AuthHandler) GetCurrentSession(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "auth-session")
	
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// ユーザー情報を取得
	roles, _ := h.repo.GetUserRoles(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user_id":       userID,
		"roles":         roles,
	})
}