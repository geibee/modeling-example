package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/modeling-example/data-model/authority-matrix/repository"
)

type AdminHandler struct {
	repo *repository.Repository
	db   *sql.DB
}

func NewAdminHandler(repo *repository.Repository, db *sql.DB) *AdminHandler {
	return &AdminHandler{
		repo: repo,
		db:   db,
	}
}

// メニュー設定のCRUD
func (h *AdminHandler) GetMenuSettings(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT menu_id, menu_name, screen_id 
		FROM menu_setting 
		ORDER BY menu_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var menus []map[string]string
	for rows.Next() {
		var menuID, menuName, screenID string
		if err := rows.Scan(&menuID, &menuName, &screenID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		menus = append(menus, map[string]string{
			"menu_id":    menuID,
			"menu_name":  menuName,
			"screen_id":  screenID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menus)
}

func (h *AdminHandler) CreateMenuSetting(w http.ResponseWriter, r *http.Request) {
	var menu map[string]string
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
		INSERT INTO menu_setting (menu_id, menu_name, screen_id) 
		VALUES ($1, $2, $3)
	`, menu["menu_id"], menu["menu_name"], menu["screen_id"])
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) UpdateMenuSetting(w http.ResponseWriter, r *http.Request) {
	// URLからmenu_idを取得
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	menuID := parts[3]

	var menu map[string]string
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
		UPDATE menu_setting 
		SET menu_id = $1, menu_name = $2, screen_id = $3 
		WHERE menu_id = $4
	`, menu["menu_id"], menu["menu_name"], menu["screen_id"], menuID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) DeleteMenuSetting(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	menuID := parts[3]

	_, err := h.db.Exec("DELETE FROM menu_setting WHERE menu_id = $1", menuID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// API画面マッピングのCRUD
func (h *AdminHandler) GetApiScreenMappings(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT screen_id, api_id 
		FROM api_screen_mapping 
		ORDER BY screen_id, api_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mappings []map[string]string
	for rows.Next() {
		var screenID, apiID string
		if err := rows.Scan(&screenID, &apiID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mappings = append(mappings, map[string]string{
			"screen_id": screenID,
			"api_id":    apiID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mappings)
}

func (h *AdminHandler) CreateApiScreenMapping(w http.ResponseWriter, r *http.Request) {
	var mapping map[string]string
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
		INSERT INTO api_screen_mapping (screen_id, api_id) 
		VALUES ($1, $2)
	`, mapping["screen_id"], mapping["api_id"])
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) UpdateApiScreenMapping(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	oldScreenID := parts[3]
	oldApiID := parts[4]

	var mapping map[string]string
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
		UPDATE api_screen_mapping 
		SET screen_id = $1, api_id = $2 
		WHERE screen_id = $3 AND api_id = $4
	`, mapping["screen_id"], mapping["api_id"], oldScreenID, oldApiID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) DeleteApiScreenMapping(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	screenID := parts[3]
	apiID := parts[4]

	_, err := h.db.Exec(`
		DELETE FROM api_screen_mapping 
		WHERE screen_id = $1 AND api_id = $2
	`, screenID, apiID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ロール一覧取得
func (h *AdminHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT role_id, role_name 
		FROM role 
		ORDER BY role_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []map[string]string
	for rows.Next() {
		var roleID, roleName string
		if err := rows.Scan(&roleID, &roleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, map[string]string{
			"role_id":   roleID,
			"role_name": roleName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// 画面一覧取得
func (h *AdminHandler) GetScreens(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT screen_id, screen_name 
		FROM screen 
		ORDER BY screen_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var screens []map[string]string
	for rows.Next() {
		var screenID, screenName string
		if err := rows.Scan(&screenID, &screenName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		screens = append(screens, map[string]string{
			"screen_id":   screenID,
			"screen_name": screenName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screens)
}

// API一覧取得
func (h *AdminHandler) GetApis(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT api_id, api_name 
		FROM api 
		ORDER BY api_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var apis []map[string]string
	for rows.Next() {
		var apiID, apiName string
		if err := rows.Scan(&apiID, &apiName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		apis = append(apis, map[string]string{
			"api_id":   apiID,
			"api_name": apiName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis)
}

// ロール画面権限の取得
func (h *AdminHandler) GetRoleScreenPermissions(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	roleID := parts[3]

	rows, err := h.db.Query(`
		SELECT screen_id, permission_level 
		FROM role_screen_permission 
		WHERE role_id = $1
	`, roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var permissions []map[string]string
	for rows.Next() {
		var screenID, permissionLevel string
		if err := rows.Scan(&screenID, &permissionLevel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		permissions = append(permissions, map[string]string{
			"screen_id": screenID,
			"permission_level": permissionLevel,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}


// 権限の更新（追加/削除）
func (h *AdminHandler) UpdateRoleScreenPermission(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	roleID := parts[3]
	screenID := parts[4]

	if r.Method == "POST" {
		// リクエストボディから権限レベルを取得
		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			// 後方互換性のため、ボディがない場合はデフォルトでVIEWER
			reqBody = map[string]string{"permission_level": "VIEWER"}
		}
		
		permissionLevel := reqBody["permission_level"]
		if permissionLevel == "" {
			permissionLevel = "VIEWER"
		}
		
		_, err := h.db.Exec(`
			INSERT INTO role_screen_permission (role_id, screen_id, permission_level) 
			VALUES ($1, $2, $3)
			ON CONFLICT (role_id, screen_id) 
			DO UPDATE SET permission_level = $3
		`, roleID, screenID, permissionLevel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		_, err := h.db.Exec(`
			DELETE FROM role_screen_permission 
			WHERE role_id = $1 AND screen_id = $2
		`, roleID, screenID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}


// ユーザー一覧取得
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT DISTINCT u.user_id, ur.role_id, r.role_name 
		FROM "user" u
		LEFT JOIN user_role ur ON u.user_id = ur.user_id
		LEFT JOIN role r ON ur.role_id = r.role_id
		ORDER BY u.user_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var userID, roleID, roleName string
		if err := rows.Scan(&userID, &roleID, &roleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// ユーザーの組織情報を取得
		orgRows, err := h.db.Query(`
			SELECT o.org_id, o.org_name 
			FROM user_organization uo
			JOIN organization o ON uo.org_id = o.org_id
			WHERE uo.user_id = $1
		`, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		var organizations []map[string]string
		for orgRows.Next() {
			var orgID, orgName string
			if err := orgRows.Scan(&orgID, &orgName); err != nil {
				orgRows.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			organizations = append(organizations, map[string]string{
				"org_id": orgID,
				"name": orgName,
			})
		}
		orgRows.Close()
		
		users = append(users, map[string]interface{}{
			"user_id":       userID,
			"name":          userID, // 実際のシステムでは別途ユーザー名フィールドを持つべき
			"role_id":       roleID,
			"role_name":     roleName,
			"organizations": organizations,
			"status":        "active", // 実際のシステムではステータスフィールドを持つべき
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 個別ユーザー取得
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID := parts[3]

	var roleID, roleName string
	err := h.db.QueryRow(`
		SELECT ur.role_id, r.role_name 
		FROM user_role ur
		LEFT JOIN role r ON ur.role_id = r.role_id
		WHERE ur.user_id = $1
		LIMIT 1
	`, userID).Scan(&roleID, &roleName)
	
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// ユーザーの組織情報を取得
	orgRows, err := h.db.Query(`
		SELECT o.org_id, o.org_name 
		FROM user_organization uo
		JOIN organization o ON uo.org_id = o.org_id
		WHERE uo.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer orgRows.Close()
	
	var organizations []map[string]string
	for orgRows.Next() {
		var orgID, orgName string
		if err := orgRows.Scan(&orgID, &orgName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		organizations = append(organizations, map[string]string{
			"org_id": orgID,
			"name": orgName,
		})
	}
	
	user := map[string]interface{}{
		"user_id":       userID,
		"name":          userID, // 実際のシステムでは別途ユーザー名フィールドを持つべき
		"role_id":       roleID,
		"role_name":     roleName,
		"organizations": organizations,
		"status":        "active", // 実際のシステムではステータスフィールドを持つべき
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ユーザー作成
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user["user_id"].(string)
	roleID := user["role_id"].(string)
	
	// トランザクション開始
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// ユーザー作成
	_, err = tx.Exec(`
		INSERT INTO "user" (user_id) 
		VALUES ($1)
	`, userID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// ロール割り当て
	_, err = tx.Exec(`
		INSERT INTO user_role (user_id, role_id) 
		VALUES ($1, $2)
	`, userID, roleID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 組織への所属を設定
	if orgs, ok := user["organizations"].([]interface{}); ok {
		for _, orgID := range orgs {
			_, err = tx.Exec(`
				INSERT INTO user_organization (user_id, org_id) 
				VALUES ($1, $2)
			`, userID, orgID.(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	
	// コミット
	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ユーザー更新
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID := parts[3]

	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roleID := user["role_id"].(string)
	
	// トランザクション開始
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 既存のロール割り当てを削除
	_, err = tx.Exec(`
		DELETE FROM user_role 
		WHERE user_id = $1
	`, userID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 新しいロール割り当て
	_, err = tx.Exec(`
		INSERT INTO user_role (user_id, role_id) 
		VALUES ($1, $2)
	`, userID, roleID)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 既存の組織所属を削除
	_, err = tx.Exec(`
		DELETE FROM user_organization 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 新しい組織所属を設定
	if orgs, ok := user["organizations"].([]interface{}); ok {
		for _, orgID := range orgs {
			_, err = tx.Exec(`
				INSERT INTO user_organization (user_id, org_id) 
				VALUES ($1, $2)
			`, userID, orgID.(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	
	// コミット
	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ユーザー削除
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID := parts[3]

	// トランザクション開始
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 組織所属を削除
	_, err = tx.Exec(`
		DELETE FROM user_organization 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// ロール割り当てを削除
	_, err = tx.Exec(`
		DELETE FROM user_role 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// ユーザーを削除
	_, err = tx.Exec(`
		DELETE FROM "user" 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// コミット
	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// 組織一覧取得
func (h *AdminHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT org_id, org_name 
		FROM organization 
		ORDER BY org_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var organizations []map[string]string
	for rows.Next() {
		var orgID, orgName string
		if err := rows.Scan(&orgID, &orgName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		organizations = append(organizations, map[string]string{
			"org_id":   orgID,
			"name": orgName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organizations)
}

// GetUsersForSwitcher はユーザー切り替え用のユーザー一覧を返します（認証不要）
func (h *AdminHandler) GetUsersForSwitcher(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT DISTINCT u.user_id, ur.role_id, r.role_name 
		FROM "user" u
		LEFT JOIN user_role ur ON u.user_id = ur.user_id
		LEFT JOIN role r ON ur.role_id = r.role_id
		ORDER BY u.user_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var userID string
		var roleID, roleName sql.NullString
		if err := rows.Scan(&userID, &roleID, &roleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		user := map[string]interface{}{
			"user_id": userID,
		}
		
		if roleID.Valid {
			user["role_id"] = roleID.String
		} else {
			user["role_id"] = ""
		}
		
		if roleName.Valid {
			user["role_name"] = roleName.String
		} else {
			user["role_name"] = "Unknown"
		}
		
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}