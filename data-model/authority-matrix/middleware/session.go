package middleware

import (
	"net/http"
	
	"github.com/gorilla/sessions"
)

var (
	// セッションストア（本番環境では環境変数から秘密鍵を取得すべき）
	Store = sessions.NewCookieStore([]byte("authority-matrix-secret-key-change-this"))
)

func init() {
	// セッションの設定
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7日間
		HttpOnly: true,
		Secure:   false, // 本番環境ではtrueにする
		SameSite: http.SameSiteLaxMode,
	}
}