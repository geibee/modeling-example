package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type AuthMiddleware struct {
	repo *repository.Repository
}

func NewAuthMiddleware(repo *repository.Repository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo}
}

// AuthorizeAPI はユーザーがAPIにアクセスする権限があるかをチェックします
// 画面権限ベースで認可を行います（API画面マッピングを通じて）
func (m *AuthMiddleware) AuthorizeAPI(apiID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// コンテキストからユーザーIDを取得（認証ミドルウェアで設定されていることを前提）
			userIDValue := r.Context().Value("userID")
			if userIDValue == nil {
				http.Error(w, "認可エラー: ユーザーIDが設定されていません", http.StatusUnauthorized)
				return
			}
			userID, ok := userIDValue.(string)
			if !ok {
				http.Error(w, "認可エラー: ユーザーIDの形式が不正です", http.StatusUnauthorized)
				return
			}

			// ユーザーのロールを取得
			userRoles, err := m.repo.GetUserRoles(userID)
			if err != nil || len(userRoles) == 0 {
				log.Printf("認可エラー: ユーザー %s にロールが割り当てられていません", userID)
				http.Error(w, "認可エラー: ユーザーにロールが割り当てられていません", http.StatusUnauthorized)
				return
			}
			log.Printf("ユーザー %s のロール: %v", userID, userRoles)

			// APIに紐づく画面を取得
			screenIDs, err := m.repo.GetScreensForAPI(apiID)
			if err != nil || len(screenIDs) == 0 {
				log.Printf("認可エラー: API %s に紐づく画面が見つかりません", apiID)
				http.Error(w, "認可エラー: このAPIは画面に紐づいていません", http.StatusForbidden)
				return
			}

			// ユーザーがいずれかの画面にアクセス権を持っているかチェック
			hasAccess := false
			hasEditPermission := false
			for _, roleID := range userRoles {
				for _, screenID := range screenIDs {
					if m.repo.HasRoleScreenPermission(roleID, screenID) {
						hasAccess = true
						// 権限レベルをチェック
						permissionLevel, err := m.repo.GetRoleScreenPermissionLevel(roleID, screenID)
						if err == nil && permissionLevel == "EDITOR" {
							hasEditPermission = true
						}
					}
				}
			}

			if !hasAccess {
				log.Printf("認可エラー: ユーザー %s はAPI %s へのアクセス権がありません", userID, apiID)
				http.Error(w, "認可エラー: このAPIへのアクセス権がありません", http.StatusForbidden)
				return
			}

			// 書き込み操作（POST, PUT, DELETE）の場合はEDITOR権限が必要
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "PATCH" {
				if !hasEditPermission {
					log.Printf("認可エラー: ユーザー %s はAPI %s へのEDITOR権限がありません", userID, apiID)
					http.Error(w, "認可エラー: このAPIへの編集権限がありません", http.StatusForbidden)
					return
				}
			}

			// ロールとEDITOR権限をコンテキストに追加
			ctx := context.WithValue(r.Context(), "roles", userRoles)
			ctx = context.WithValue(ctx, "hasEditPermission", hasEditPermission)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthorizeScreen はユーザーが画面にアクセスする権限があるかをチェックします
func (m *AuthMiddleware) AuthorizeScreen(screenID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// コンテキストからユーザーIDを取得
			userIDValue := r.Context().Value("userID")
			if userIDValue == nil {
				http.Error(w, "認可エラー: ユーザーIDが設定されていません", http.StatusUnauthorized)
				return
			}
			userID, ok := userIDValue.(string)
			if !ok {
				http.Error(w, "認可エラー: ユーザーIDの形式が不正です", http.StatusUnauthorized)
				return
			}

			// ユーザーのロールを取得
			userRoles, err := m.repo.GetUserRoles(userID)
			if err != nil || len(userRoles) == 0 {
				http.Error(w, "認可エラー: ユーザーにロールが割り当てられていません", http.StatusUnauthorized)
				return
			}

			// ロールが画面にアクセス権を持っているかチェック
			hasAccess := false
			for _, roleID := range userRoles {
				if m.repo.HasRoleScreenPermission(roleID, screenID) {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, "認可エラー: この画面へのアクセス権がありません", http.StatusForbidden)
				return
			}

			// ロールをコンテキストに追加
			ctx := context.WithValue(r.Context(), "roles", userRoles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HasAuthority はユーザーが必要な権限を持っているかをチェックします
func HasAuthority(roles []string, requiredRole string) bool {
	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

// ExtractUserID はセッションからユーザーIDを抽出します
func ExtractUserID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := Store.Get(r, "auth-session")
			if err != nil {
				http.Error(w, "セッションエラー", http.StatusInternalServerError)
				return
			}

			// セッションから認証情報を取得
			authenticated, ok := session.Values["authenticated"].(bool)
			if !ok || !authenticated {
				// 未認証の場合の処理
				// APIパスまたはAcceptヘッダーでレスポンスタイプを決定
				if strings.HasPrefix(r.URL.Path, "/api/") || strings.Contains(r.Header.Get("Accept"), "application/json") {
					// APIリクエストの場合はJSONエラーを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error":"認証が必要です"}`))
				} else {
					// ブラウザからのアクセスの場合はログインページにリダイレクト
					http.Redirect(w, r, "/login", http.StatusFound)
				}
				return
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				http.Error(w, "セッションが無効です", http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}